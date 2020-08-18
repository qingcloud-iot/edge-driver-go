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

import "fmt"

type topicCodec interface {
	build(args ...string) string
	parse(string) (string, error)
}
type codec interface {
	build(params Metadata, args ...string) interface{}
	parse(payload []byte) (interface{}, error)
}

type propertyCodec struct {
}

func (p propertyCodec) build(args ...string) string {
	if len(args) != 2 {
		return ""
	}
	return fmt.Sprintf(devicePropertiesReport, args[1], args[0])
}
func (p propertyCodec) parse(string) (string, error) {
	return "", nil
}

type eventCodec struct {
}

func (p eventCodec) build(args ...string) string {
	if len(args) != 3 {
		return ""
	}
	return fmt.Sprintf(deviceEventsReport, args[1], args[0], args[2])
}
func (p eventCodec) parse(string) (string, error) {
	return "", nil
}

type statusCodec struct {
}

func (p statusCodec) build(args ...string) string {
	if len(args) != 2 {
		return ""
	}
	return fmt.Sprintf(deviceStatusReport, args[1], args[0])
}
func (p statusCodec) parse(string) (string, error) {
	return "", nil
}
