package edge_driver_go

var defaultServerOptions = options{}

type options struct {
	Name     string   //驱动名称
	Broker   string   `json:"broker"`   //hub地址
	Services []string `json:"services"` //服务定义
}

type ServerOption interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncServerOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

//设置驱动名称
func SetName(name string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Name = name
	})
}

//设置连接hub地址
func SetBroker(url string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Broker = url
	})
}

//设置服务调用方法
func SetServices(services []string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Services = services
	})
}
