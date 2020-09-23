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
