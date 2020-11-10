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

import "context"

type ValueData struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}
type Metadata map[string]interface{}
type MetadataMsg map[string]ValueData

//edge service call
type OnEdgeServiceCall func(args Metadata) (*Reply, error)
type OnEndServiceCall func(name string, args Metadata) (*Reply, error)

//edge set and get
type OnSetServiceCall func(args Metadata) error
type OnGetServiceCall func(args []string) (Metadata, error)

//user service call
type OnUserServiceCall func(data []byte) ([]byte, error)

//config change call
type ConfigChangeFunc func(t string, config []byte)

//sub device interface
type Client interface {
	//report device online to cloud
	Online(ctx context.Context) error
	//report device offline to cloud
	Offline(ctx context.Context) error
	//report device property to cloud
	ReportProperties(ctx context.Context, params Metadata) error
	//report device property to cloud with tags
	ReportPropertiesWithTags(ctx context.Context, params Metadata, tags Metadata) error
	//report device property to cloud with tags and time
	ReportPropertiesWithTagsEx(ctx context.Context, params MetadataMsg, tags Metadata) error
	//report device event to cloud
	ReportEvent(ctx context.Context, eventId string, params Metadata) error
	//report user device message to cloud
	ReportUserMessage(ctx context.Context, data []byte) error
	//report device info to cloud
	ReportDeviceInfo(ctx context.Context, params *DeviceMsg) error
}
type ConnectLost func(err error)
type messageArrived func(topic string, payload []byte)

//describe device info
type config interface {
	DeviceId() string                 //device id
	ThingId() string                  //thing id
	Token() string                    //token
	Services() []string               //services
	Metadata() map[string]interface{} //device metadata
}
