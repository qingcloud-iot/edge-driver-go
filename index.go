package edge_driver_go

import "context"

type Metadata map[string]interface{}

//服务调用
type CallService func(ctx context.Context, name string, params Metadata) (Metadata, error)

//边端设备sdk接口
type Client interface {
	//RegisterAndOnlineByCode() error
	Init() error                              //驱动初始化
	Online() error                            //设备上线通知
	Offline() error                           //设备下线通知
	ReportProperties() error                  //上报属性
	ReportEvents() error                      //上报事件
	GetDriverInfo() (Metadata, error)         //获取驱动配置
	SetProperties(Metadata) error             //设置设备属性
	GetProperties([]string) (Metadata, error) //获取设备属性
	Close() error                             //销毁驱动
}

// edge sdk初始化
func NewClient(opt ...ServerOption) Client {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return nil
}
