/*
 * Copyright (C) 2020 Yunify, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this work except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package edge_driver_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"sync/atomic"
	"time"
)

type edgeDriver struct {
	ctx             context.Context
	cancel          context.CancelFunc
	validate        validate
	client          mqtt.Client //hub client
	status          uint32      //0:not connected, 1:connected
	url             string      //meta address
	deviceId        string
	thingId         string
	edgeServiceCall EdgeCallService //service call func
	endServiceCall  CallService     //service call func
	userServiceCall UserCallService //user service call func
	logger          Logger
}

// edge sdk init
func NewClient(opt ...ServerOption) Client {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	edge := &edgeDriver{
		validate:        newDataValidate(),
		edgeServiceCall: opts.edgeServiceCall,
		endServiceCall:  opts.endServiceCall,
		userServiceCall: opts.userServiceCall,
		logger:          opts.logger,
		ctx:             ctx,
		cancel:          cancel,
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(opts.broker).
		SetClientID(opts.name).
		SetUsername(opts.name).
		SetPassword(opts.name).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetConnectionLostHandler(func(client mqtt.Client, e error) {
			if edge.logger != nil {
				edge.logger.Warn("edge connect lost")
			}
			//heartbeat lost
			atomic.StoreUint32(&edge.status, hubNotConnected)
		}).
		SetOnConnectHandler(func(client mqtt.Client) {
			if edge.logger != nil {
				edge.logger.Warn("edge connect success call")
			}
			atomic.StoreUint32(&edge.status, hubConnected)
			//edge service
			for _, v := range opts.edgeServices {
				if token := edge.client.Subscribe(v, byte(0), func(client mqtt.Client, i mqtt.Message) {
					edge.edgeCall(i.Topic(), i.Payload())
				}); token.Wait() && token.Error() != nil {
					if edge.logger != nil {
						edge.logger.Warn("edge sub error")
					}
				}
			}
			//end service
		})
	client := mqtt.NewClient(options)
	go edge.connect(client) //reconnected
	edge.client = client
	return edge
}
func (e *edgeDriver) edgeCall(topic string, payload []byte) {
	var (
		msg  message
		req  *serviceRequest
		name string
		data Metadata
		resp *serviceReply
		buf  []byte
		err  error
	)
	defer func() {
		if err != nil {
			if e.logger != nil {
				e.logger.Error(topic, err.Error())
			}
		}
	}()
	name, err = msg.parseServiceName(topic)
	if err != nil {
		return
	}
	req, err = msg.parseServiceMsg(payload)
	if err != nil {
		return
	}
	if err = e.validate.validateServiceInput(context.Background(), e.deviceId, name, req.Params); err != nil {
		return
	}
	if e.logger != nil {
		e.logger.Warn(topic, payload)
	}
	resp = &serviceReply{
		Id:   req.Id,
		Code: RPC_SUCCESS,
		Data: make(Metadata),
	}
	if e.edgeServiceCall != nil {
		if data, err = e.edgeServiceCall(name, req.Params); err != nil {
			resp.Code = RPC_FAIL
		}
		if err = e.validate.validateServiceOutput(context.Background(), e.deviceId, name, data); err != nil {
			resp.Code = RPC_FAIL
		}
		resp.Data = data
	}
	buf, err = json.Marshal(resp)
	if err != nil {
		return
	}
	if token := e.client.Publish(topic+"_reply", byte(0), false, buf); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply err:%s", token.Error()))
		}
	} else {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply  topic:%s,data:%s", topic+"_reply", string(buf)))
		}
	}
	return
}
func (e *edgeDriver) getSubDevice(device string) (string, error) {
	return "iott-1ac28fzjUM", nil
}
func (e *edgeDriver) connect(client mqtt.Client) {
	for {
		select {
		case <-e.ctx.Done():
			return
		default:
		}
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			if e.logger != nil {
				e.logger.Warn("edge connect retry......")
			}
			time.Sleep(3 * time.Second)
			continue
		} else {
			if e.logger != nil {
				e.logger.Info("edge connect success")
			}
			atomic.StoreUint32(&e.status, hubConnected)
			return
		}
	}
}

func (e *edgeDriver) GetEdgeDeviceConfig(ctx context.Context) error {
	return nil
}
func (e *edgeDriver) GetEndDeviceConfig(ctx context.Context) error {
	return nil
}
func (e *edgeDriver) Online(ctx context.Context, deviceId string) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildStatusTopic(deviceId, thingId)
	data = msg.buildHeartbeatMsg(deviceId, thingId, online)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) Offline(ctx context.Context, deviceId string) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildStatusTopic(deviceId, thingId)
	data = msg.buildHeartbeatMsg(deviceId, thingId, offline)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) ReportProperties(ctx context.Context, deviceId string, params Metadata) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if err = e.validate.validateProperties(ctx, deviceId, params); err != nil {
		return err
	}
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildPropertyTopic(deviceId, thingId)
	data = msg.buildPropertyMsg(deviceId, thingId, params)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) ReportEvent(ctx context.Context, deviceId string, eventName string, params Metadata) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if err = e.validate.validateEvent(ctx, deviceId, eventName, params); err != nil {
		return err
	}
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildEventTopic(deviceId, thingId, eventName)
	data = msg.buildEventMsg(deviceId, thingId, eventName, params)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) GetDriverInfo(ctx context.Context) (interface{}, error) {
	return nil, nil
}
func (e *edgeDriver) SetProperties(ctx context.Context, params Metadata) error {
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	return nil
}
func (e *edgeDriver) GetProperties(ctx context.Context, properties []string) (Metadata, error) {
	return nil, nil
}
func (e *edgeDriver) Close() error {
	atomic.StoreUint32(&e.status, hubNotConnected)
	if e.client != nil {
		e.client.Disconnect(250)
	}
	e.cancel()
	return nil
}
