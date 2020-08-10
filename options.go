package edge_driver_go

import uuid "github.com/satori/go.uuid"

var defaultServerOptions = options{
	name:            "driver-" + uuid.NewV4().String(),
	edgeServices:    []string{},
	edgeServiceCall: nil,
	endServiceCall:  nil,
	broker:          "tcp://127.0.0.1:1883",
	logger:          newLogger(),
}

type options struct {
	name            string          `json:"name"`        //driver name
	broker          string          `json:"broker"`      //hub address
	metaBroker      string          `json:"meta_broker"` //meta service address
	edgeServices    []string        `json:"services"`    //edge service define
	edgeServiceCall EdgeCallService //service call func
	endServiceCall  CallService     //service call func
	userServiceCall UserCallService //user service call func
	logger          Logger          //logger
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

//set driver name
func SetName(name string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.name = name
	})
}

//set hub client
func SetBroker(url string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.broker = url
	})
}

//register service
func SetRegisterServices(services []string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.edgeServices = services
	})
}

//set edge service callback
func SetEdgeCallService(call CallService) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.edgeServiceCall = call
	})
}

//set user service callback
func SetUserCallService(call UserCallService) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.userServiceCall = call
	})
}

//set user service callback
func SetURL(url string) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.metaBroker = url
	})
}

//set logger
func SetLogger(logger Logger) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.Logger = logger
	})
}
