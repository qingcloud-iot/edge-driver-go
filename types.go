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
	"errors"
)

const (
	messageVersion = "v0.0.1"
	hubBroker      = "tcp://127.0.0.1:1883"
	metadataBroker = "127.0.0.1:1889"
	fileToken      = "/etc/token"
)
const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
)

const (
	RPC_SUCCESS = 200 //success
	RPC_FAIL    = 201 //rpc timeout
)

const (
	hubNotConnected = 0 //not connected
	hubConnected    = 1 //connected
)
const (
	online  = "online"  //online
	offline = "offline" //offline
)
const (
	devicePropertyType = "thing.property.post"
	deviceEventType    = "thing.event.%s.post"
)
const (
	message_notify = "/iot/config/change"
)

var (
	notConnected    = errors.New("hub not connected")
	pubMessageError = errors.New("pub message fail")
	topicError      = errors.New("parse topic error")
)

//device status report
type deviceStatus struct {
	DeviceId string `json:"device_id"`
	ThingId  string `json:"thing_id"`
	Status   string `json:"status"`
	Time     int64  `json:"time"`
}
type property struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}
type messageMeta struct {
	DeviceId  string   `json:"entityId"`
	ThingId   string   `json:"modelId"`
	SourceId  []string `json:"sourceId,omitempty"`
	EpochTime int64    `json:"epochTime,omitempty"`
}

//device property
type thingPropertyMsg struct {
	Id       string               `json:"id"`
	Version  string               `json:"version"`
	Type     string               `json:"type"`
	Metadata *messageMeta         `json:"metadata"`
	Params   map[string]*property `json:"params"`
}

//device event
type thingEventMsg struct {
	Id       string       `json:"id"`
	Version  string       `json:"version"`
	Type     string       `json:"type"`
	Metadata *messageMeta `json:"metadata"`
	Params   *eventData   `json:"params"`
}
type eventData struct {
	Value Metadata `json:"value"`
	Time  int64    `json:"time"`
}
type serviceRequest struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  Metadata `json:"params"`
}
type serviceGetRequest struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  []string `json:"params"`
}
type Reply struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
type serviceReply struct {
	Code int         `json:"code"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}
