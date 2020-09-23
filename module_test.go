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
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRegisterEdgeService(t *testing.T) {
	err := RegisterEdgeService("xxxx", func(args Metadata) (reply *Reply, e error) {
		return
	})
	assert.Nil(t, err)
	err = ReportEdgeProperties(context.Background(), Metadata{"int32": rand.Int()})
	assert.Nil(t, err)
	err = ReportEdgeEvent(context.Background(), "event", Metadata{"int32": rand.Int()})
	assert.Nil(t, err)
	time.Sleep(3 * time.Second)
}

func TestGetConfig(t *testing.T) {
	res, err := GetConfig()
	assert.Nil(t, err)
	for _, v := range res {
		t.Log(v)
	}
}
func TestGetDriverInfo(t *testing.T) {
	res, err := GetDriverInfo()
	assert.Nil(t, err)
	t.Log(res)
}
func TestDiscovery(t *testing.T) {
	err := ReportDiscovery(context.Background(), "onvif", Metadata{"name": "hello world"})
	assert.Nil(t, err)
}
