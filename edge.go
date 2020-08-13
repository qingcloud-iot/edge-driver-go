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
	name            string
	broker          string
	validate        validate
	client          mqtt.Client //hub client
	status          uint32      //0:not connected, 1:connected
	url             string      //meta address
	edgeServices    []string    `json:"services"` //edge service define
	deviceId        string
	thingId         string
	cache           configCache
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
		cache:           newCache(opts.metaBroker),
		name:            opts.name,
		broker:          opts.broker,
		logger:          opts.logger,
		ctx:             ctx,
		cancel:          cancel,
	}
	edge.init()
	return edge
}
func (e *edgeDriver) init() {
	options := mqtt.NewClientOptions()
	options.AddBroker(e.broker).
		SetClientID(e.name).
		SetUsername(e.name).
		SetPassword(e.name).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			if e.logger != nil {
				e.logger.Warn("edge connect lost")
			}
			//heartbeat lost
			atomic.StoreUint32(&e.status, hubNotConnected)
		}).
		SetOnConnectHandler(func(client mqtt.Client) {
			if e.logger != nil {
				e.logger.Warn("edge connect success call")
			}
			atomic.StoreUint32(&e.status, hubConnected)
			//edge service
			if err := e.edgeServiceInit(e.edgeServices); err != nil {
				if e.logger != nil {
					e.logger.Error(err.Error())
				}
			}
			//end service
			if info, err := e.cache.GetEndDevicesConfig(e.ctx); err != nil {
				for _, v := range info {
					if err := e.endServiceInit(v.GetServices()); err != nil {
						if e.logger != nil {
							e.logger.Error(err.Error())
						}
					}
				}
			}
			e.client.Subscribe(message_notify, byte(0), func(client mqtt.Client, i mqtt.Message) {
				if e.logger != nil {
					e.logger.Warn("client restart ", i.Topic(), i.Payload())
				}
				go e.restart()
			})
		})
	client := mqtt.NewClient(options)
	go e.connect(client) //reconnected
	e.client = client
}
func (e *edgeDriver) restart() {
	e.client.Disconnect(250)
	e.init()
}

//register end service
func (e *edgeDriver) endServiceInit(topics []string) error {
	filters := make(map[string]byte)
	for _, v := range topics {
		filters[v] = byte(0)
	}
	token := e.client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
		e.endCall(message.Topic(), message.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (e *edgeDriver) endCall(topic string, payload []byte) {
	var (
		msg      message
		req      *serviceRequest
		name     string
		deviceId string
		data     Metadata
		resp     *serviceReply
		buf      []byte
		err      error
	)
	defer func() {
		if err != nil {
			if e.logger != nil {
				e.logger.Error(topic, err.Error())
			}
		}
	}()
	deviceId, name, err = msg.parseServiceName(topic)
	if err != nil {
		return
	}
	req, err = msg.parseServiceMsg(payload)
	if err != nil {
		return
	}
	if err = e.validate.validateServiceInput(context.Background(), deviceId, name, req.Params); err != nil {
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
		if data, err = e.endServiceCall(deviceId, name, req.Params); err != nil {
			resp.Code = RPC_FAIL
		}
		if err = e.validate.validateServiceOutput(context.Background(), deviceId, name, data); err != nil {
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
}

//register end service
func (e *edgeDriver) edgeServiceInit(topics []string) error {
	filters := make(map[string]byte)
	for _, v := range topics {
		filters[v] = byte(0)
	}
	token := e.client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
		e.edgeCall(message.Topic(), message.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
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
	_, name, err = msg.parseServiceName(topic)
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
func (e *edgeDriver) getSubDevice(deviceId string) (string, error) {
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
