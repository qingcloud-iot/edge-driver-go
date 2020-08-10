package edge_driver_go

import uuid "github.com/satori/go.uuid"

var defaultServerOptions = options{
	Name:            "driver-" + uuid.NewV4().String(),
	Services:        []string{},
	ServiceCall:     nil,
	UserServiceCall: nil,
	Broker:          "tcp://127.0.0.1:1883",
	Logger:          newLogger(),
}

type options struct {
	Name            string          `json:"name"`        //driver name
	Broker          string          `json:"broker"`      //hub address
	MetaBroker      string          `json:"meta_broker"` //meta service address
	Services        []string        `json:"services"`    //service define
	ServiceCall     CallService     //service call func
	UserServiceCall UserCallService //user service call func
	Logger          Logger          //logger
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
func SetRegisterServices(services []string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Services = services
	})
}

//设置服务调用回调
func SetCallService(call CallService) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.ServiceCall = call
	})
}

//设置自定义格式回调
func SetUserCallService(call UserCallService) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.UserServiceCall = call
	})
}

//设置自定义格式回调
func SetLogger(logger Logger) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Logger = logger
	})
}
