/**
 * @Author: hexing
 * @Description:
 * @File:  module_test.go
 * @Version: 1.0.0
 * @Date: 20-8-19 上午7:28
 */

package edge_driver_go

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "getConfig",
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDriverInfo(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "getConfig",
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDriverInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDriverInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDriverInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetConfigChange(t *testing.T) {
	type args struct {
		call ConfigChangeFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "getConfig",
			args: args{func(config []byte) {
				t.Log(config)
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetConfigChange(tt.args.call)
		})
	}
}

func TestSetConnectLost(t *testing.T) {
	type args struct {
		call ConnectLost
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "getConfig",
			args: args{func(err error) {
				t.Log(err.Error())
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetConnectLost(tt.args.call)
		})
	}
}
