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

//get device config
func GetConfig() (Metadata, error) {
	return getSessionIns().getConfig()
}

//get driver config
func GetDriverInfo() (Metadata, error) {
	return getSessionIns().getDriver()
}

//set connect lost handle
func SetConnectLost(call ConnectLost) {
	getSessionIns().setConnectLost(call)
}

//set config change handle
func SetConfigChange(call ConfigChangeFunc) {
	getSessionIns().setConfigChange(call)
}
