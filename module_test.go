/**
 * @Author: hexing
 * @Description:
 * @File:  module_test.go
 * @Version: 1.0.0
 * @Date: 20-8-19 上午7:28
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
