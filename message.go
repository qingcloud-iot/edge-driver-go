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
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	edgeServiceLen = 8 //service topic len
	configLen      = 5 //config topic len
)
const (
	deviceSetProperty          = "/sys/%s/%s/thing/property/base/set"
	deviceGetProperty          = "/sys/%s/%s/thing/property/base/get"
	deviceService              = "/sys/%s/%s/thing/service/%s/call"
	deviceStatusReport         = "/as/mqtt/status/%s/%s"
	devicePropertiesReport     = "/sys/%s/%s/thing/property/platform/post"
	userDevicePropertiesReport = "/sys/%s/%s/user/msg"
	userDeviceService          = "/sys/%s/%s/user/down/+/call"
	deviceEventsReport         = "/sys/%s/%s/thing/event/%s/post"
	deviceInfoReport           = "/sys/%s/%s/thing/deviceinfo/post"
	configChange               = "/iot/internal/%s/notify"
	deviceDiscoveryReport      = "/sys/%s/device/discovery/post"
)

type message struct {
}

//build device status topic
func (m message) buildStatusTopic(deviceId, thingId string) string {
	return fmt.Sprintf(deviceStatusReport, thingId, deviceId)
}

//build device status topic
func (m message) buildUserTopic(deviceId, thingId string) string {
	return fmt.Sprintf(userDevicePropertiesReport, thingId, deviceId)
}

//build device set topic
func (m message) buildSetTopic(deviceId, thingId string) string {
	return fmt.Sprintf(deviceSetProperty, thingId, deviceId)
}

//build device get topic
func (m message) buildGetTopic(deviceId, thingId string) string {
	return fmt.Sprintf(deviceGetProperty, thingId, deviceId)
}

//build device get topic
func (m message) buildUserServiceTopic(deviceId, thingId string) string {
	return fmt.Sprintf(userDeviceService, thingId, deviceId)
}

//build discovery topic
func (m message) buildDiscoveryTopic(t string) string {
	return fmt.Sprintf(deviceDiscoveryReport, t)
}

//build device get topic
func (m message) buildServiceTopic(deviceId, thingId string, services []string) []string {
	result := make([]string, len(services))
	for k, v := range services {
		result[k] = fmt.Sprintf(deviceService, thingId, deviceId, v)
	}
	return result
}

//build device property topic
func (m message) buildPropertyTopic(deviceId, thingId string) string {
	return fmt.Sprintf(devicePropertiesReport, thingId, deviceId)
}

//build device property topic
func (m message) buildDeviceInfoTopic(deviceId, thingId string) string {
	return fmt.Sprintf(deviceInfoReport, thingId, deviceId)
}

//build device property topic
func (m message) buildUserPropertyTopic(deviceId, thingId string) string {
	return fmt.Sprintf(devicePropertiesReport, thingId, deviceId)
}

// build device status struct
func (m message) buildHeartbeatMsg(deviceId, thingId, status string) []byte {
	data := &deviceStatus{
		DeviceId: deviceId,
		ThingId:  thingId,
		Status:   status,
		Time:     time.Now().UnixNano() / 1e6,
	}
	buf, _ := json.Marshal(data)
	return buf
}

//build device property data
func (m message) buildPropertyMsg(deviceId, thingId string, meta Metadata) []byte {
	id := uuid.NewV4().String()
	params := make(map[string]interface{})
	for k, _ := range meta {
		property := &property{
			Value: meta[k],
			Time:  time.Now().UnixNano() / 1e6,
		}
		params[k] = property
	}
	message := &thingPropertyMsg{
		Id:      id,
		Version: messageVersion,
		Type:    devicePropertyType,
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
		},
		Params: params,
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device property data
func (m message) buildPropertyMsgWithTags(deviceId, thingId string, meta Metadata, tags Metadata) []byte {
	id := uuid.NewV4().String()
	params := make(map[string]interface{})
	for k, _ := range meta {
		property := &property{
			Value: meta[k],
			Time:  time.Now().UnixNano() / 1e6,
		}
		params[k] = property
	}
	message := &thingPropertyMsg{
		Id:      id,
		Version: messageVersion,
		Type:    devicePropertyType,
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
			Tags:      tags,
		},
		Params: params,
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device property data
func (m message) buildPropertyMsgWithTagsEx(deviceId, thingId string, meta MetadataMsg, tags Metadata) []byte {
	id := uuid.NewV4().String()
	//params := make(map[string]interface{})
	//for k, _ := range meta {
	//	property := &property{
	//		Value: meta[k],
	//		Time:  time.Now().UnixNano() / 1e6,
	//	}
	//	params[k] = property
	//}
	message := &thingPropertyMsg{
		Id:      id,
		Version: messageVersion,
		Type:    devicePropertyType,
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
			Tags:      tags,
		},
		Params: meta.Data(),
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device info data
func (m message) buildDiscoveryMsg(deviceId, thingId string, meta Metadata) []byte {
	id := uuid.NewV4().String()
	message := &thingPropertyMsg{
		Id:      id,
		Version: messageVersion,
		Type:    deviceDeviceInfoType,
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
		},
		Params: meta,
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device info data
func (m message) buildDeviceInfoMsg(deviceId, thingId string, meta Metadata) []byte {
	id := uuid.NewV4().String()
	meta["deviceId"] = deviceId
	meta["thingId"] = thingId
	message := &thingPropertyMsg{
		Id:      id,
		Version: messageVersion,
		Type:    deviceDeviceInfoType,
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
		},
		Params: meta,
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device event topic
func (m message) buildEventTopic(deviceId, thingId, eventName string) string {
	return fmt.Sprintf(deviceEventsReport, thingId, deviceId, eventName)
}

//build device event data
func (m message) buildEventMsg(deviceId, thingId string, eventName string, meta Metadata) []byte {
	id := uuid.NewV4().String()
	message := &thingEventMsg{
		Id:      id,
		Version: messageVersion,
		Type:    fmt.Sprintf(deviceEventType, eventName),
		Metadata: &messageMeta{
			DeviceId:  deviceId,
			ThingId:   thingId,
			SourceId:  []string{deviceId},
			EpochTime: time.Now().UnixNano() / 1e6,
		},
		Params: &eventData{
			Value: meta,
			Time:  time.Now().UnixNano() / 1e6,
		},
	}
	buf, _ := json.Marshal(message)
	return buf
}

//parse device service method
func (m message) parseServiceMethod(topic string) (string, string, error) {
	kv := strings.Split(topic, "/")
	if len(kv) != edgeServiceLen {
		return "", "", topicError
	}
	return kv[2], kv[6], nil
}

//parse device config type
func (m message) parseConfigType(topic string) (string, error) {
	kv := strings.Split(topic, "/")
	if len(kv) != configLen {
		return "", topicError
	}
	return kv[4], nil
}

//parse device config type
func (m message) parseResponseMsg(payload []byte) (*serviceRequest, error) {
	message := &serviceRequest{}
	err := json.Unmarshal(payload, message)
	if err != nil {
		return message, err
	}
	return message, nil
}

//parse device get method type
func (m message) parseGetServiceMsg(payload []byte) (*serviceGetRequest, error) {
	message := &serviceGetRequest{}
	err := json.Unmarshal(payload, message)
	if err != nil {
		return message, err
	}
	return message, nil
}
