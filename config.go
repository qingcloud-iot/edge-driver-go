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

type config struct {
	metadata map[string]interface{}
	services []string
	deviceId string
	thingId  string
}

func newConfig(token string) (DeviceConfig, error) {
	var (
		deviceId string
		thingId  string
		conf     DeviceConfig
		err      error
	)
	if deviceId, thingId, err = parseToken(token); err != nil {
		return nil, err
	}
	conf = &config{
		deviceId: deviceId,
		thingId:  thingId,
	}
	return conf, nil
}
func (c *config) update() error {
	return nil
}
func (c *config) DeviceId() string {
	return c.deviceId
}
func (c *config) ThingId() string {
	return c.thingId
}
func (c *config) Services() []string {
	return c.services
}
func (c *config) Metadata() map[string]interface{} {
	return c.metadata
}
