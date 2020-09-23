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

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	getSessionIns().init()
	getSessionIns().setConfigChange(func(tp string, config []byte) {
		fmt.Println("config change:", tp, string(config))
	})
	err := getSessionIns().publish("/iot/internal/notify/edgeDeviceChanged", []byte("hello world"))
	assert.Nil(t, err)
	err = getSessionIns().subscribe("/sys/1/2/thing/service/set/call", func(topic string, payload []byte) {
		fmt.Println("subscribe:", topic, string(payload))
	})
	assert.Nil(t, err)
	err = getSessionIns().publish("/sys/1/2/thing/service/set/call", []byte("call hello world"))
	assert.Nil(t, err)
	time.Sleep(3 * time.Second)
}
func TestRequestEdge(t *testing.T) {
	res, err := getSessionIns().getEdgeInfo()
	assert.Nil(t, err)
	fmt.Println(res)
}
func TestRequestDriver(t *testing.T) {
	res, err := getSessionIns().getDriver()
	assert.Nil(t, err)
	fmt.Println(res)
}

func TestGetModel(t *testing.T) {
	res, err := getSessionIns().getModel("iotd-0adf702f-8c1c-489e-bde0-01788ac674c3")
	assert.Nil(t, err)
	fmt.Println(res)
}
func TestGetEdgeInfo(t *testing.T) {
	res, err := getSessionIns().getEdgeInfo()
	assert.Nil(t, err)
	fmt.Println(res)
}
