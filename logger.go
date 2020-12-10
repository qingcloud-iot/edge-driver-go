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
	"os"
)

// Logger specifies logging API.
type Logger interface {
	// Debug logs any object in JSON format on debug level.
	Debug(a ...interface{})
	// Info logs any object in JSON format on info level.
	Info(a ...interface{})
	// Warn logs any object in JSON format on warning level.
	Warn(a ...interface{})
	// Error logs any object in JSON format on error level.
	Error(a ...interface{})
}
type logger struct {
	open bool
}

func newLogger() Logger {
	open := os.Getenv("EDGE_ENABLE_LOG")
	if open == "true" {
		return &logger{open: true}
	}
	return &logger{}
}
func (l *logger) Debug(a ...interface{}) {
	if l.open {
		fmt.Fprintln(os.Stdout, a...)
	}
}
func (l *logger) Info(a ...interface{}) {
	if l.open {
		fmt.Fprintln(os.Stdout, a...)
	}
}
func (l *logger) Warn(a ...interface{}) {
	if l.open {
		fmt.Fprintln(os.Stdout, a...)
	}
}
func (l *logger) Error(a ...interface{}) {
	if l.open {
		fmt.Fprintln(os.Stderr, a...)
	}
}
