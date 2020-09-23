# edge-driver-go

## edge 设备接入　SDK　GO版本
边缘端sdk是为边缘端设备接入提供基础能力，所有接入到edge的设备都是需要通过驱动进行接入

## 设备驱动分为三个部分（用户实现）
- 连接管理　设备和edge设备建立连接，我们不限制建立通信连接的协议，可根据客户业务灵活选择
- 数据转换（上行）  设备接入驱动需要将设备数据转换成符合Qingcloud　IoT物模型规范的数据格式（可选）
- 数据和命令处理（下行）　驱动可以处理云端对于设备的操作请求，并且将结果返回云端（可选）

其中,数据转换和数据和命令处理部分为可选,一种是转换成Qingcloud　IoT物模型规范格式数据，另外一种是将数据直接透传不做解析，直接上传云端

## 设备驱动SDK接口
调试（环境变量）
- ENV_EDGE_DEVICE_ID 边设备id
- ENV_EDGE_THING_ID 边设备模型id
- ENV_EDGE_HUB_HOST,EDGE_HUB_PORT 默认为本地地址（tcp://127.0.0.1:1883），调试过程中可以修改,方便调试
- ENV_EDGE_META_ADDRESS 默认为本地地址（http://127.0.0.1:9611），调试过程中可以修改,方便调试

### 驱动配置管理接口
```go
/*
 * 边端上报探测发现的设备信息
 *
 * 阻塞接口
 *
 * ctx:         @ctx, 接口超时控制上下文
 * eventId:     @deviceType, 事件标识符.
 * params:      @Metadata, 属性数组.
 * err:         @err 成功返回nil,  失败返回错误信息.
 *
 */
func ReportDiscovery(ctx context.Context,deviceType string, meta Metadata) error
/*
 * 边端获取本配置信息(包括子设备属性　token等)
 *
 * 阻塞接口
 *
 * config:		 @config 子设备配置信息.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 *
 */
func GetConfig() (config []*SubDeviceInfo, err error)
/*
 * 边端获取本驱动信息
 *
 * 阻塞接口.
 *
 * info:		 @info 本驱动信息.
 * err:			 @err 成功返回nil,  失败返回错误信息.
 */
func GetDriverInfo() (info string,err error)
/*
 * 获取模型详细信息（包括扩展描述）
 *
 * id:          @id 模型id.
 *
 * 阻塞接口.
 * info:        @info 模型详细信息
 * err:         @err 成功返回nil,  失败返回错误信息.
 */
func GetDeviceModel(id string) (info string, err error)
/*
 * 边端hub离线通知
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
 * call:        @call, 子设备属性设置接口.
 * err:         @err 成功返回nil,  失败返回错误信息.
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
 * call:        @call, 子设备属性获取接口.
 * err:         @err 成功返回nil,  失败返回错误信息.
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
 * call:        @call, 子设备服务调用接口.
 * err:         @err 成功返回nil,  失败返回错误信息.
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
 * call:        @call, 子设备自定义格式设备服务调用接口.
 * err:         @err 成功返回nil,  失败返回错误信息.
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
     * ctx:         @ctx, 接口超时控制上下文
     *
     * 阻塞接口.
     * err:         @err 成功返回nil,  失败返回错误信息.
     */
	Online(ctx context.Context) error                        //设备上线通知
    /*
     * 子设备下线
     *
     * ctx:         @ctx, 接口超时控制上下文
     *
     * 阻塞接口.
     * err:         @err 成功返回nil,  失败返回错误信息.
     */
	Offline(ctx context.Context) error                       //设备下线通知
    /*
     * 子设备上报属性, 设备具有的属性的设备能力描述在设备物模型规定.
     *
     * 上报属性, 可以上报一个, 也可以多个一起上报.
     *
     * ctx:         @ctx, 接口超时控制上下文
     * params:      @Metadata, 属性数组.
     *
     * 阻塞接口.
     * err:         @err 成功返回nil,  失败返回错误信息.
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
     * ctx:         @ctx, 接口超时控制上下文
     * data:        @data, 设备自定义数据.
     *
     * 阻塞接口.
     * err:         @err 成功返回nil,  失败返回错误信息.
     */
    ReportUserMessage(ctx context.Context,data []byte) error     //上报自定义数据
    /*
     * 子设备上报设备扩展信息.
     *
     *
     * ctx:         @ctx, 接口超时控制上下文
     * data:        @data, 设备数据.
     *
     * 阻塞接口.
     * err:         @err 成功返回nil,  失败返回错误信息.
     */
    ReportDeviceInfo(ctx context.Context, params Metadata) error			//上报设备数据
}
```

## 示例

下面是一个驱动sdk示例代码

```go
func main() {
    //获取驱动下分配子设备信息
	subs, err := edge_driver_go.GetConfig()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(subs))
	for _, v := range subs {
        //设备token可用
		if v.TokenStatus == edge_driver_go.Enable {
			go func(token string) {
				defer func() {
					wg.Done()
				}()
				var opts []edge_driver_go.ServerOption
				opt := edge_driver_go.SetEndServiceCall(func(name string, args edge_driver_go.Metadata) (reply *edge_driver_go.Reply, e error) {
					fmt.Println(name, args)
					return
				})
				opts = append(opts, opt)
				opt = edge_driver_go.SetGetServiceCall(func(args []string) (metadata edge_driver_go.Metadata, e error) {
					fmt.Println(args)
					return
				})
				opts = append(opts, opt)
				opt = edge_driver_go.SetSetServiceCall(func(args edge_driver_go.Metadata) error {
					fmt.Println(args)
					return nil
				})
				opts = append(opts, opt)
				client, err := edge_driver_go.NewEndClient(token,
					opts...)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				for {
                    //驱动上报子设备上线
					err := client.Online(context.Background())
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(2 * time.Second)
                    //云端定义端设备属性模型（temp）
					err = client.ReportProperties(context.Background(), edge_driver_go.Metadata{"temp": rand.Float32()})
					if err != nil {
						fmt.Println(err)
					}
                    //云端定义端设备事件模型（temperatureEvent）
					err = client.ReportEvent(context.Background(), "temperatureEvent", edge_driver_go.Metadata{"temperature": rand.Float32(), "reason": true})
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(2 * time.Second)
				}
			}(v.Token)
		}
	}
	wg.Wait()
}

```