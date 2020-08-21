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

/*
 * 边端获取本配置信息(包括子设备属性　token等)
 * 阻塞接口, 成功返回nil,  失败返回错误信息.
 *
 */
func GetConfig() ([]byte, error) {
	return getSessionIns().getConfig()
}

/*
 * 边端获取本驱动信息
 * 阻塞接口, 成功返回nil,  失败返回错误信息.
 *
 */
func GetDriverInfo() ([]byte, error) {
	return getSessionIns().getDriver()
}

/*
 * 边端注册服务, 设备注册的服务在设备能力描述在设备物模型规定.
 *
 * 上报属性, 可以上报一个, 也可以多个一起上报.
 *
 * ctx:          接口超时控制上下文
 * serviceId:    @serviceId, 服务标识符.
 * call:    	 @call, 服务回调接口.
 *
 * 阻塞接口, 成功返回nil,  失败返回错误信息.
 *
 */
func RegisterEdgeService(serviceId string, call OnEdgeServiceCall) error {
	var (
		msg    message
		err    error
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
			Code: RPC_SUCCESS,
			Data: make(Metadata),
		}
		if call != nil {
			if reply, err = call(req.Params); err != nil {
				resp.Code = RPC_FAIL
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

/*
 * 边端上报属性, 设备具有的属性在设备能力描述在设备物模型规定.
 *
 * 上报属性, 可以上报一个, 也可以多个一起上报.
 *
 * ctx:          接口超时控制上下文
 * params:       @Metadata, 属性数组.
 *
 * 阻塞接口, 成功返回nil,  失败返回错误信息.
 *
 */
func ReportEdgeProperties(ctx context.Context, params Metadata) error {
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

/*
 * 边端上报事件, 设备具有的事件在设备能力描述在设备物模型规定.
 *
 * 上报事件, 单个事件上报.
 *
 * ctx:          接口超时控制上下文
 * eventId:      @eventId, 事件标识符.
 * params:       @Metadata, 属性数组.
 *
 * 阻塞接口, 成功返回nil,  失败返回错误信息.
 *
 */
func ReportEdgeEvent(ctx context.Context, eventId string, params Metadata) error {
	var (
		topic string
		msg   message
		data  []byte
	)
	topic = msg.buildEventTopic(deviceId, thingId, eventId)
	data = msg.buildEventMsg(deviceId, thingId, eventId, params)
	return getSessionIns().publish(topic, data)
}

//set connect lost handle
func SetConnectLost(call ConnectLost) {
	getSessionIns().setConnectLost(call)
}

//set config change handle
func SetConfigChange(call ConfigChangeFunc) {
	getSessionIns().setConfigChange(call)
}
