package edge_driver_go

import "errors"

const (
	messageVersion = "v0.0.1"
)

const (
	RPC_SUCCESS = 200  //success
	RPC_FAIL    = 201  //unkonw fail
	RPC_TIMEOUT = 1001 // rpc timeout
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

var (
	notConnected    = errors.New("not connected")
	pubMessageError = errors.New("pub message fail")
	topicError      = errors.New("parse topic error")
)

//device status report
type DeviceStatus struct {
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
type serviceReply struct {
	Code int         `json:"code"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}
