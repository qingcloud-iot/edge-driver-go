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

type Metadata map[string]interface{}

//edge service call
//type OnEdgeServiceCall func(name string, args Metadata) (Metadata, error)
type OnEndServiceCall func(deviceId, name string, args Metadata) (*Reply, error)

//edge set and get
type OnSetServiceCall func(args Metadata) error
type OnGetServiceCall func(args []string) (Metadata, error)

//user service call
type OnUserServiceCall func(data []byte) ([]byte, error)

//config change call
type ConfigChangeFunc func(t string, config []byte)

//边端设备sdk接口
type Client interface {
	//Init() error                                       //init
	Online(context.Context) error                        //report device online to cloud
	Offline(context.Context) error                       //report device offline to cloud
	ReportProperties(context.Context, Metadata) error    //report device property to cloud
	ReportEvent(context.Context, string, Metadata) error //report device event to cloud
	ReportUserMessage(context.Context, []byte) error     //report user device message to cloud
}
type ConnectLost func(err error)
type messageArrived func(topic string, payload []byte)

//describe device info
type DeviceConfig interface {
	DeviceId() string //device id
	ThingId() string  //thing id
	Services() []string
	Metadata() map[string]interface{} //device metadata
}
