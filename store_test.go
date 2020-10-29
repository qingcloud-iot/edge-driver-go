package edge_driver_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetValue(t *testing.T) {
	err := SetValue("test", []byte("xxxxxxx"))
	assert.Nil(t, err)
}
