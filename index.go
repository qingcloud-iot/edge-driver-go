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
