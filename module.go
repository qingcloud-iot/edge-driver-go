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

//get edge sub device list
func GetConfig() (config []byte, err error) {
	return getSessionIns().getConfig()
}

//get edge sub device model
func GetModel(id string) (config []byte, err error) {
	return getSessionIns().getModel(id)
}

//get edge sub driver info
func GetDriverInfo() (info []byte, err error) {
	return getSessionIns().getDriver()
}

//register edge device service
func RegisterEdgeService(serviceId string, call OnEdgeServiceCall) (err error) {
	var (
		msg    message
		req    *serviceRequest
		reply  *Reply
		resp   *serviceReply
		buf    []byte
		logger Logger
		//methodName string
	)
	logger = newLogger()
	err = getSessionIns().subscribes(msg.buildServiceTopic(deviceId, thingId, []string{serviceId}), func(topic string, payload []byte) {
		defer func() {
			if err != nil {
				logger.Error(topic, err.Error())
			}
		}()
		deviceId, _, err = msg.parseServiceMethod(topic)
		if err != nil {
			return
		}
		req, err = msg.parseResponseMsg(payload)
		if err != nil {
			return
		}
		resp = &serviceReply{
			Id:   req.Id,
			Code: RpcSuccess,
			Data: make(Metadata),
		}
		if call != nil {
			if reply, err = call(req.Params); err != nil {
				resp.Code = RpcFail
			} else {
				resp.Code = reply.Code
				resp.Data = reply.Data
			}
			buf, err = json.Marshal(resp)
			if err != nil {
				return
			}
			if err = getSessionIns().publish(topic+"_reply", buf); err != nil {
				logger.Error(fmt.Sprintf("edge requestServiceReply err:%s", err.Error()))
			} else {
				logger.Error(fmt.Sprintf("edge requestServiceReply  topic:%s,data:%s", topic+"_reply", string(buf)))
			}
		} else {
			logger.Warn("edge callback not set")
		}
	})
	if err != nil {
		return err
	}
	return nil
}

//report edge device property
func ReportEdgeProperties(ctx context.Context, params Metadata) (err error) {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
		)
		topic = msg.buildPropertyTopic(deviceId, thingId)
		data = msg.buildPropertyMsg(deviceId, thingId, params)
		return getSessionIns().publish(topic, data)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}

//report edge device event
func ReportEdgeEvent(ctx context.Context, eventId string, params Metadata) (err error) {
	done := wait(func() error {
		var (
			topic string
			msg   message
			data  []byte
		)
		topic = msg.buildEventTopic(deviceId, thingId, eventId)
		data = msg.buildEventMsg(deviceId, thingId, eventId, params)
		return getSessionIns().publish(topic, data)
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return rpcTimeout
	}
}

//set lost call
func SetConnectLost(call ConnectLost) {
	getSessionIns().setConnectLost(call)
}

//set config change call
func SetConfigChange(call ConfigChangeFunc) {
	getSessionIns().setConfigChange(call)
}
