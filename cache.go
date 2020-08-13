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

import "context"

type config struct {
	metadata map[string]interface{}
	services []string
}

func NewConfig(services []string, metadata map[string]interface{}) configMessage {
	return &config{
		metadata: metadata,
		services: services,
	}
}
func (c *config) GetDeviceId() string {
	return ""
}
func (c *config) GetThingId() string {
	return ""
}
func (c *config) GetServices() []string {
	return c.services
}
func (c *config) GetMetadata() map[string]interface{} {
	return c.metadata
}

type cache struct {
}

func newCache(url string) configCache {
	return &cache{}
}
func (c *cache) run() {
}

func (c *cache) GetEndDevicesConfig(ctx context.Context) ([]configMessage, error) {
	return nil, nil
}
func (c *cache) GetEndDeviceConfig(ctx context.Context, id string) (configMessage, error) {
	return NewConfig([]string{}, map[string]interface{}{}), nil
}
