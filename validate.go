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

//todo need fix
//validate device thing model
type validate interface {
	validateProperties(ctx context.Context, deviceId string, metadata Metadata) (Metadata, error)
	validatePropertiesEx(ctx context.Context, deviceId string, metadata MetadataMsg) (MetadataMsg, error)
	validateEvent(ctx context.Context, deviceId string, eventName string, metadata Metadata) error
	validateServiceInput(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error
	validateServiceOutput(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error
}

//validate device thing model
type dataValidate struct {
}

func newDataValidate() validate {
	return &dataValidate{}
}

func (v *dataValidate) validateProperties(ctx context.Context, deviceId string, metadata Metadata) (Metadata, error) {
	var (
		thing *ThingModel
		resp  Metadata
		err   error
	)
	if thing, err = getSessionIns().getModel(deviceId); err != nil {
		return resp, err
	}
	for k, _ := range metadata {
		if _, ok := thing.Properties[k]; !ok {
			continue
		}
		resp[k] = metadata[k]
	}
	return resp, nil
}
func (v *dataValidate) validatePropertiesEx(ctx context.Context, deviceId string, metadata MetadataMsg) (MetadataMsg, error) {
	var (
		thing *ThingModel
		resp  MetadataMsg
		err   error
	)
	if thing, err = getSessionIns().getModel(deviceId); err != nil {
		return resp, err
	}
	for k, _ := range metadata {
		if _, ok := thing.Properties[k]; !ok {
			continue
		}
		resp[k] = metadata[k]
	}
	return resp, nil
}
func (v *dataValidate) validateEvent(ctx context.Context, deviceId string, eventName string, metadata Metadata) error {
	return nil
}
func (v *dataValidate) validateServiceInput(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error {
	return nil
}
func (v *dataValidate) validateServiceOutput(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error {
	return nil
}
