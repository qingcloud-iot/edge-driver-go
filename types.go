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

type TokenStatus string

const (
	Enable  TokenStatus = "enable"  //token enable
	Disable TokenStatus = "disable" //token disable
)

const (
	messageVersion    = "v0.0.1"
	hubBroker         = "tcp://127.0.0.1:1883"
	metadataBroker    = "http://127.0.0.1:9611"
	edgeInfoRequest   = metadataBroker + "/internal/data/edge_info/"   //request edge info
	edgeDriverRequest = metadataBroker + "/internal/data/edge_driver/" //request driver info
	subDeviceRequest  = metadataBroker + "/internal/data/child_device/"
	userThingId       = "iott-end-user-system"
)
const (
	EdgeDeviceChanged   = "edgeDeviceChanged"   //edge device config change
	EdgeConfigChanged   = "edgeConfigChanged"   //edge thing config change
	DriverConfigChanged = "driverConfigChanged" //driver config change
	SubDeviceChanged    = "subDeviceChanged"    //sub device config change
)
const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
)

const (
	RpcSuccess = 200 //success
	RpcFail    = 201 //rpc timeout
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
	devicePropertyType   = "thing.property.post"
	deviceDeviceInfoType = "thing.deviceinfo.post"
	deviceEventType      = "thing.event.%s.post"
)

var (
	rpcTimeout      = errors.New("rpc timeout")
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
	Id       string                 `json:"id"`
	Version  string                 `json:"version"`
	Type     string                 `json:"type"`
	Metadata *messageMeta           `json:"metadata"`
	Params   map[string]interface{} `json:"params"`
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

//dev info
type edgeDevInfo struct {
	Id      string `json:"deviceId"`
	ThingId string `json:"thingId"`
}
type driverResult struct {
	DriverId  string    `json:"driverId"`
	DriverCfg string    `json:"driverCfg"`
	Channels  []channel `json:"channels"`
}
type channel struct {
	SubDeviceId  string `json:"subDeviceId"`
	SubDeviceCfg string `json:"subDeviceCfg"`
	ChannelCfg   string `json:"channelCfg"`
}

//sub device info
type SubDeviceInfo struct {
	Token       string                 `json:"token"`        //device token
	TokenStatus TokenStatus            `json:"token_status"` //device token status, enable or disable
	DeviceId    string                 `json:"device_id"`    //device id
	Ext         map[string]interface{} `json:"ext"`          //device custom config
	ChannelCfg  map[string]interface{} `json:"channel_cfg"`  //sub device config, example
}

//sub device info
type device struct {
	DeviceId     string `json:"deviceId"`
	TokenContent string `json:"tokenContent"`
	TokenStatus  string `json:"tokenStatus"`
	ThingId      string `json:"thingId"`
}

//driver info
type driverInfo struct {
	Id       string                 `json:"id"`
	Protocol string                 `json:"protocol"`
	Version  string                 `json:"version"`
	Custom   map[string]interface{} `json:"custom"`
}

type reply struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
