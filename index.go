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
type EdgeCallService func(name string, params Metadata) (Metadata, error)
type CallService func(deviceId, name string, params Metadata) (Metadata, error)

//user service call
type UserCallService func(data []byte) ([]byte, error)

//config change call
type ConfigChangeFunc func(config interface{})

//边端设备sdk接口
type Client interface {
	GetEdgeDeviceConfig(context.Context) error                   //获取边设备配置
	GetEndDeviceConfig(context.Context) error                    //获取子设备配置
	Online(context.Context, string) error                        //设备上线通知
	Offline(context.Context, string) error                       //设备下线通知
	ReportProperties(context.Context, string, Metadata) error    //上报属性
	ReportEvent(context.Context, string, string, Metadata) error //上报事件
	GetDriverInfo(context.Context) (interface{}, error)          //获取驱动配置
	SetProperties(context.Context, Metadata) error               //设置设备属性
	GetProperties(context.Context, []string) (Metadata, error)   //获取设备属性
	Close() error                                                //销毁驱动
}
