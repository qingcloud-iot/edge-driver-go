/**
 * @Author: hexing
 * @Description:
 * @File:  logger_test.go
 * @Version: 1.0.0
 * @Date: 20-8-10 上午11:36
 */

package edge_driver_go

import (
	"reflect"
	"testing"
)

func Test_logger_Debug(t *testing.T) {
	type args struct {
		a []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "logger",
			args: args{a: []interface{}{"xxxxxx"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Debug(tt.args.a)
		})
	}
}

func Test_logger_Error(t *testing.T) {
	type args struct {
		a []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "logger",
			args: args{a: []interface{}{"xxxxxx"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Error(tt.args.a)
		})
	}
}

func Test_logger_Info(t *testing.T) {
	type args struct {
		a []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "logger",
			args: args{a: []interface{}{"xxxxxx"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Info(tt.args.a)
		})
	}
}

func Test_logger_Warn(t *testing.T) {
	type args struct {
		a []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "logger",
			args: args{a: []interface{}{"xxxxxx"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Warn(tt.args.a)
		})
	}
}

func Test_newLogger(t *testing.T) {
	tests := []struct {
		name string
		want Logger
	}{
		{
			name: "logger",
			want: newLogger(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
