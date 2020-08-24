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
)

type edgeClient struct {
	ctx      context.Context
	cancel   context.CancelFunc
	validate validate
	config   Config
	//edgeServiceCall OnEdgeServiceCall //service call func
	endServiceCall  OnEndServiceCall  //service call func
	userServiceCall OnUserServiceCall //user service call func
	setServiceCall  OnSetServiceCall  //set service call func
	getServiceCall  OnGetServiceCall  //get service call func
	logger          Logger
}

// edge sdk init
//
func NewEdgeClient(token string, opt ...ServerOption) Client {
	var (
		config Config
		err    error
		opts   options
	)
	if token == "" {
		panic("token is illegal")
	}
	config, err = newDeviceConfig(token)
	if err != nil {
		panic(err)
	}
	opts = defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	edge := &edgeClient{
		validate: newDataValidate(),
		//edgeServiceCall: opts.edgeServiceCall,
		endServiceCall:  opts.endServiceCall,
		userServiceCall: opts.userServiceCall,
		setServiceCall:  opts.setServiceCall,
		getServiceCall:  opts.getServiceCall,
		logger:          opts.logger,
		config:          config,
		ctx:             ctx,
		cancel:          cancel,
	}
	return edge
}
func (e *edgeClient) init() error {
	var (
		err error
		msg message
	)
	//end service
	err = getSessionIns().subscribe(msg.buildSetTopic(e.config.DeviceId(), e.config.ThingId()), e.endCall)
	if err != nil {
		return err
	}
	err = getSessionIns().subscribe(msg.buildGetTopic(e.config.DeviceId(), e.config.ThingId()), e.getCall)
	if err != nil {
		return err
	}
	err = getSessionIns().subscribes(msg.buildServiceTopic(e.config.DeviceId(), e.config.ThingId(), e.config.Services()), e.endCall)
	if err != nil {
		return err
	}
	err = getSessionIns().subscribe(msg.buildUserServiceTopic(e.config.DeviceId(), e.config.ThingId()), e.userCall)
	if err != nil {
		return err
	}
	return nil
}
func (e *edgeClient) setCall(topic string, payload []byte) {
	var (
		msg  message
		req  *serviceRequest
		resp *serviceReply
		buf  []byte
		err  error
	)
	req, err = msg.parseResponseMsg(payload)
	if err != nil {
		return
	}
	resp = &serviceReply{
		Id:   req.Id,
		Code: RPC_SUCCESS,
		Data: make(Metadata),
	}
	if e.setServiceCall != nil {
		if err = e.setServiceCall(req.Params); err != nil {
			resp.Code = RPC_FAIL
		}
	}
	buf, err = json.Marshal(resp)
	if err != nil {
		return
	}
	if err = getSessionIns().publish(topic+"_reply", buf); err != nil {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply err:%s", err.Error()))
		}
	} else {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply  topic:%s,data:%s", topic+"_reply", string(buf)))
		}
	}
}
func (e *edgeClient) getCall(topic string, payload []byte) {
	var (
		msg  message
		req  *serviceGetRequest
		resp *serviceReply
		data Metadata
		buf  []byte
		err  error
	)
	req, err = msg.parseGetServiceMsg(payload)
	if err != nil {
		return
	}
	resp = &serviceReply{
		Id:   req.Id,
		Code: RPC_SUCCESS,
		Data: make(Metadata),
	}
	if e.setServiceCall != nil {
		if data, err = e.getServiceCall(req.Params); err != nil {
			resp.Code = RPC_FAIL
		} else {
			resp.Data = data
		}
	}
	buf, err = json.Marshal(resp)
	if err != nil {
		return
	}
	if err = getSessionIns().publish(topic+"_reply", buf); err != nil {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply err:%s", err.Error()))
		}
	} else {
		if e.logger != nil {
			e.logger.Error(fmt.Sprintf("requestServiceReply  topic:%s,data:%s", topic+"_reply", string(buf)))
		}
	}
}
func (e *edgeClient) endCall(topic string, payload []byte) {
	var (
		msg        message
		req        *serviceRequest
		methodName string
		deviceId   string
		data       Metadata
		reply      *Reply
		resp       *serviceReply
		buf        []byte
		err        error
	)
	defer func() {
		if err != nil {
			if e.logger != nil {
				e.logger.Error(topic, err.Error())
			}
		}
	}()
	deviceId, methodName, err = msg.parseServiceMethod(topic)
	if err != nil {
		return
	}
	req, err = msg.parseResponseMsg(payload)
	if err != nil {
		return
	}
	if err = e.validate.validateServiceInput(context.Background(), deviceId, methodName, req.Params); err != nil {
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
	if e.endServiceCall != nil {
		if reply, err = e.endServiceCall(methodName, req.Params); err != nil {
			resp.Code = RPC_FAIL
		} else {
			resp.Code = reply.Code
			resp.Data = reply.Data
		}
		if err = e.validate.validateServiceOutput(context.Background(), deviceId, methodName, data); err != nil {
			resp.Code = RPC_FAIL
		}
		buf, err = json.Marshal(resp)
		if err != nil {
			return
		}
		if err = getSessionIns().publish(topic+"_reply", buf); err != nil {
			if e.logger != nil {
				e.logger.Error(fmt.Sprintf("requestServiceReply err:%s", err.Error()))
			}
		} else {
			if e.logger != nil {
				e.logger.Error(fmt.Sprintf("requestServiceReply  topic:%s,data:%s", topic+"_reply", string(buf)))
			}
		}
	} else {
		//if e.logger != nil {
		//	e.logger.Warn("callback not set")
		//}
	}
}
func (e *edgeClient) userCall(topic string, payload []byte) {
	var (
		data []byte
		err  error
	)
	defer func() {
		if err != nil {
			if e.logger != nil {
				e.logger.Error(topic, err.Error())
			}
		}
	}()
	if e.logger != nil {
		e.logger.Info(topic, payload)
	}
	if e.userServiceCall != nil {
		if data, err = e.userServiceCall(payload); err != nil {
			return
		} else {
			if err = getSessionIns().publish(topic+"_reply", data); err != nil {
				if e.logger != nil {
					e.logger.Error(fmt.Sprintf("userCall err:%s", err.Error()))
				}
			} else {
				if e.logger != nil {
					e.logger.Error(fmt.Sprintf("userCall  topic:%s,data:%s", topic+"_reply", string(data)))
				}
			}
		}
	} else {
		//if e.logger != nil {
		//	e.logger.Warn("user callback not set")
		//}
	}
}
func (e *edgeClient) ReportUserMessage(ctx context.Context, payload []byte) error {
	done := wait(func() error {
		var (
			topic string
			msg   message
		)
		topic = msg.buildUserTopic(e.config.DeviceId(), e.config.ThingId())
		return getSessionIns().publish(topic, payload)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}
func (e *edgeClient) Online(ctx context.Context) error {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
			err   error
		)
		topic = msg.buildStatusTopic(e.config.DeviceId(), e.config.ThingId())
		data = msg.buildHeartbeatMsg(e.config.DeviceId(), e.config.ThingId(), online)
		err = getSessionIns().publish(topic, data)
		if err != nil {
			return err
		}
		return e.init()
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}
func (e *edgeClient) Offline(ctx context.Context) error {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
		)
		topic = msg.buildStatusTopic(e.config.DeviceId(), e.config.ThingId())
		data = msg.buildHeartbeatMsg(e.config.DeviceId(), e.config.ThingId(), offline)
		return getSessionIns().publish(topic, data)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}

func (e *edgeClient) ReportProperties(ctx context.Context, params Metadata) error {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
			//thingId string
			err error
		)
		if err = e.validate.validateProperties(ctx, e.config.DeviceId(), params); err != nil {
			return err
		}
		topic = msg.buildPropertyTopic(e.config.DeviceId(), e.config.ThingId())
		data = msg.buildPropertyMsg(e.config.DeviceId(), e.config.ThingId(), params)
		return getSessionIns().publish(topic, data)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}
func (e *edgeClient) ReportEvent(ctx context.Context, eventId string, params Metadata) error {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
			//thingId string
			err error
		)
		if err = e.validate.validateProperties(ctx, e.config.DeviceId(), params); err != nil {
			return err
		}
		topic = msg.buildEventTopic(e.config.DeviceId(), e.config.ThingId(), eventId)
		data = msg.buildEventMsg(e.config.DeviceId(), e.config.ThingId(), eventId, params)
		return getSessionIns().publish(topic, data)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}
