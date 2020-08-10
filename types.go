package edge_driver_go

import "errors"

const (
	messageVersion = "v0.0.1"
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
