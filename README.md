# edge-driver-go

## edge 设备接入　SDK　GO版本
边缘端sdk是为边缘端设备接入提供基础能力，所有接入到edge的设备都是需要通过驱动进行接入

## 设备驱动分为三个部分（用户实现）
- 连接管理　设备和edge设备建立连接，我们不限制建立通信连接的协议，可根据客户业务灵活选择
- 数据转换（上行）  设备接入驱动需要将设备数据转换成符合Qingcloud　IoT物模型规范的数据格式（可选）
- 数据和命令处理（下行）　驱动可以处理云端对于设备的操作请求，并且将结果返回云端（可选）

其中,数据转换和数据和命令处理部分为可选,一种是转换成Qingcloud　IoT物模型规范格式数据，另外一种是将数据直接透传不做解析，直接上传云端

## 设备驱动分SDK接口
### 驱动配置管理接口
```go
//获取设备相关信息
func GetConfig() (Metadata, error)
//获取驱动相关信息
func GetDriverInfo() (Metadata, error)
//设置连接丢失回调
func SetConnectLost(call ConnectLost) 
//设置配置变更回调
func SetConfigChange(call ConfigChangeFunc) 

```
### 边设备模块接口
```go
//注册边端服务调用回调
func RegisterEdgeService(string,OnEdgeServiceCall)
//上报边端属性
func ReportEdgeProperties(context.Context,Metadata) error 
//上报边端事件
func ReportEdgeEvent(context.Context,string,Metadata) error
```
### 子设备模块管理接口
```go
//设备下行数据相关接口
//设置属性设置接口
func SetSetServiceCall(call OnSetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.setServiceCall = call
	})
}
//设置属性获取接口
func SetGetServiceCall(call OnGetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.getServiceCall = call
	})
}
//设置终端设备服务调用接口
func SetEndServiceCall(call OnEndServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.endServiceCall = call
	})
}
//设置自定义格式设备接口
func SetUserServiceCall(call OnUserServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.userServiceCall = call
	})
}
//子设备sdk接口
type Client interface {
	Init() error                                         //初始化服务
	Online(context.Context) error                        //设备上线通知
	Offline(context.Context) error                       //设备下线通知
	ReportProperties(context.Context, Metadata) error    //上报属性
	ReportEvent(context.Context, string, Metadata) error //上报事件
    ReportUserMessage(context.Context, []byte) error     //上报自定义数据
}
```