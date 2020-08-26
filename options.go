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

var defaultServerOptions = options{
	//edgeServiceCall: nil,
	endServiceCall:  nil,
	userServiceCall: nil,
	setServiceCall:  nil,
	getServiceCall:  nil,
	logger:          newLogger(),
}

type options struct {
	//module Module
	//edgeServiceCall OnEdgeServiceCall 			//service call func
	endServiceCall  OnEndServiceCall  //service call func
	userServiceCall OnUserServiceCall //user service call func
	setServiceCall  OnSetServiceCall  //set service call func
	getServiceCall  OnGetServiceCall  //get service call func
	logger          Logger            //logger
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

func SetSetServiceCall(call OnSetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.setServiceCall = call
	})
}
func SetGetServiceCall(call OnGetServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.getServiceCall = call
	})
}
func SetEndServiceCall(call OnEndServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.endServiceCall = call
	})
}
func SetUserServiceCall(call OnUserServiceCall) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.userServiceCall = call
	})
}

//set logger
func SetLogger(logger Logger) ServerOption {
	return newFuncServerOption(func(i *options) {
		i.logger = logger
	})
}
