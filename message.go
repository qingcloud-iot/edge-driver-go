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
	edgeServiceLen = 8
)
const (
	deviceStatusReport     = "/as/mqtt/status/%s/%s"
	devicePropertiesReport = "/sys/%s/%s/thing/property/platform/post"
	deviceEventsReport     = "/sys/%s/%s/thing/event/%s/post"
)

type message struct {
}

//build device status topic
func (m message) buildStatusTopic(deviceId, thingId string) string {
	return fmt.Sprintf(deviceStatusReport, thingId, deviceId)
}

// build device status struct
func (m message) buildHeartbeatMsg(deviceId, thingId, status string) []byte {
	data := &DeviceStatus{
		DeviceId: deviceId,
		ThingId:  thingId,
		Status:   status,
		Time:     time.Now().UnixNano() / 1e6,
	}
	buf, _ := json.Marshal(data)
	return buf
}

//build device property topic
func (m message) buildPropertyTopic(deviceId, thingId string) string {
	return fmt.Sprintf(devicePropertiesReport, thingId, deviceId)
}
func (m message) buildPropertyMsg(deviceId, thingId string, meta Metadata) []byte {
	id := uuid.NewV4().String()
	params := make(map[string]*property)
	for k, _ := range meta {
		property := &property{
			Value: meta[k],
			Time:  time.Now().Unix() * 1000,
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
			EpochTime: time.Now().Unix() * 1000,
		},
		Params: params,
	}
	buf, _ := json.Marshal(message)
	return buf
}

//build device event topic
func (m message) buildEventTopic(deviceId, thingId, eventName string) string {
	return fmt.Sprintf(deviceEventsReport, thingId, deviceId, eventName)
}
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
			EpochTime: time.Now().Unix() * 1000,
		},
		Params: &eventData{
			Value: meta,
			Time:  time.Now().Unix() * 1000,
		},
	}
	buf, _ := json.Marshal(message)
	return buf
}
func (m message) parseServiceName(topic string) (string, error) {
	kv := strings.Split(topic, "/")
	if len(kv) != edgeServiceLen {
		return "", topicError
	}
	return kv[6], nil
}
func (m message) parseServiceMsg(payload []byte) (*serviceRequest, error) {
	message := &serviceRequest{}
	err := json.Unmarshal(payload, message)
	if err != nil {
		return message, err
	}
	return message, nil
}
