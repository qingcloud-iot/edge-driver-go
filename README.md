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
/*
 * 边端获取本配置信息(包括子设备属性　token等)
 *
 * 阻塞接口
 *
 * config:		 @config 本配置信息.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 *
 */
func GetConfig() (config []byte, err error)
/*
* 获取模型相关的详细信息
*
* 阻塞接口.
*
* config:		 @config 模型相关的详细信息.
* err:			 @err 成功返回nil,  失败返回错误信息.
*/
func GetModel(id string) (config []byte,err error) {
	return getSessionIns().getModel(id)
}

/*
 * 边端获取本驱动信息
 *
 * 阻塞接口.
 *
 * info:		 @info 本驱动信息.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func GetDriverInfo() (info []byte,err error)
/*
 * 边端离线回调通知(异常)
 *
 * call:       @call, 离线回调.
 */
func SetConnectLost(call ConnectLost) 
/*
 * 边端配置变更回调通知
 *
 * call:       @call, 离线回调.
 */
func SetConfigChange(call ConfigChangeFunc) 

```
### 边设备模块接口
```go

/*
 * 边端注册服务, 设备注册的服务在设备能力描述在设备物模型规定.
 *
 * 注册服务, 接口中实现服务逻辑
 *
 * serviceId:    @serviceId, 服务标识符.
 * call:    	 @call, 服务回调接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func RegisterEdgeService(serviceId string,call OnEdgeServiceCall)
/*
 * 边端上报属性, 设备具有的属性在设备能力描述在设备物模型规定.
 *
 * 上报属性, 可以上报一个, 也可以多个一起上报.
 *
 * ctx:          @ctx, 接口超时控制上下文
 * params:       @Metadata, 属性数组.
 *
 * 阻塞接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func ReportEdgeProperties(ctx context.Context,params Metadata) error 
/*
 * 边端上报事件, 设备具有的事件在设备能力描述在设备物模型规定.
 *
 * 上报事件, 单个事件上报.
 *
 * ctx:          @ctx, 接口超时控制上下文
 * eventId:      @eventId, 事件标识符.
 * params:       @Metadata, 属性数组.
 *
 * 阻塞接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func ReportEdgeEvent(ctx context.Context,eventId string, params Metadata) error
```
### 子设备模块管理接口
```go
/*
 * 子设备属性设置接口
 *
 * 设置子设备属性设置回调
 *
 * call:    	 @call, 子设备属性设置接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func SetSetServiceCall(call OnSetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.setServiceCall = call
	})
}
/*
 * 子设备属性获取接口
 *
 * 设置子设备属性获取回调
 *
 * call:    	 @call, 子设备属性获取接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func SetGetServiceCall(call OnGetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.getServiceCall = call
	})
}
/*
 * 子设备服务调用接口设置
 *
 * 设置子设备子设备服务调用接口
 *
 * call:    	 @call, 子设备服务调用接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func SetEndServiceCall(call OnEndServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.endServiceCall = call
	})
}
/*
 * 子设备自定义格式设备服务调用接口设置
 *
 * 设置子设备自定义格式设备服务调用接口设置
 *
 * call:    	 @call, 子设备自定义格式设备服务调用接口.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func SetUserServiceCall(call OnUserServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.userServiceCall = call
	})
}
//子设备sdk接口
type Client interface {
    /*
     * 子设备上线
     *
     * ctx:          @ctx, 接口超时控制上下文
     *
     * 阻塞接口.
     * err:			 @err 成功返回nil,  失败返回错误信息.
     */
	Online(ctx context.Context) error                        //设备上线通知
    /*
     * 子设备下线
     *
     * ctx:          @ctx, 接口超时控制上下文
     *
     * 阻塞接口.
     * err:			 @err 成功返回nil,  失败返回错误信息.
     */
	Offline(ctx context.Context) error                       //设备下线通知
    /*
     * 子设备上报属性, 设备具有的属性的设备能力描述在设备物模型规定.
     *
     * 上报属性, 可以上报一个, 也可以多个一起上报.
     *
     * ctx:          @ctx, 接口超时控制上下文
     * params:       @Metadata, 属性数组.
     *
     * 阻塞接口.
     * err:			 @err 成功返回nil,  失败返回错误信息.
     */
	ReportProperties(ctx context.Context,params Metadata) error    //上报属性
    /*
     * 子设备上报事件, 设备具有的事件的设备能力描述在设备物模型规定.
     *
     * 上报事件, 单个事件上报.
     *
     * ctx:          @ctx, 接口超时控制上下文
     * eventId:      @eventId, 事件标识符.
     * params:       @Metadata, 属性数组.
     *
     * 阻塞接口.
     * err:			 @err 成功返回nil,  失败返回错误信息.
     */
	ReportEvent(ctx context.Context,eventId string,params Metadata) error //上报事件
    /*
     * 子设备上报自定义数据.
     *
     *
     * ctx:       @ctx, 接口超时控制上下文
     * data:      @data, 设备自定义数据.
     *
     * 阻塞接口.
     * err:			 @err 成功返回nil,  失败返回错误信息.
     */
    ReportUserMessage(ctx context.Context,data []byte) error     //上报自定义数据
}
```

## example