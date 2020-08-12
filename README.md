# edge-driver-go

## edge 设备接入　SDK　GO版本
边缘端sdk是为边缘端设备接入提供基础能力，所有接入到edge的设备都是需要通过驱动进行接入

## 设备驱动分为三个部分（用户实现）
- 连接管理　设备和edge设备建立连接，我们不限制建立通信连接的协议，可根据客户业务灵活选择
- 数据转换（上行）  设备接入驱动需要将设备数据转换成符合Qingcloud　IoT物模型规范的数据格式（可选）
- 数据和命令处理（下行）　驱动可以处理云端对于设备的操作请求，并且将结果返回云端（可选）

其中,数据转换和数据和命令处理部分为可选,一种是转换成Qingcloud　IoT物模型规范格式数据，另外一种是将数据直接透传不做解析，直接上传云端

## 设备驱动分SDK接口
```go
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
```