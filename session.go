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
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ins       *session
	once       sync.Once
	driverId   string
	driverName string
)

func getSessionIns() *session {
	os.Setenv("DRIVER_ID", "test_driver_id")
	if driverId == "" {
		driverId = "edge.driver." + os.Getenv("DRIVER_ID")
		driverName = "edge.driver." + os.Getenv("DRIVER_ID")
	}
	once.Do(func() {
		_ins = &session{
			client: nil,
			status: hubNotConnected,
			logger: newLogger(),
		}
		_ins.init()
	})
	return _ins
}
func init() {
}

//module api
type session struct {
	lock           sync.RWMutex
	subDevices     map[string]Client //sub device
	client         mqtt.Client       //hub client
	metadataClient *http.Client
	deviceId       string
	thingId        string
	topics         []string
	status         uint32           //0:not connected, 1:connected
	connectLost    ConnectLost      //connect lost callback
	configChange   ConfigChangeFunc //config change
	logger         Logger
	//messageArrived messageArrived			//message callback
}

func (s *session) init() {
	var (
		err    error
		result *edgeDevInfo
	)
	s.metadataClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},
	}
	result, err = s.getEdgeInfo()
	for err != nil {
		result, err = s.getEdgeInfo()
		if err != nil {
			s.logger.Warn("sdk get edge info fail,", err.Error())
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	s.deviceId = result.Id
	s.thingId = result.ThingId
	options := mqtt.NewClientOptions()
	options.AddBroker(hubBroker).
		SetClientID(driverId).
		SetUsername(driverName).
		SetPassword(driverName).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			//heartbeat lost
			atomic.StoreUint32(&s.status, hubNotConnected)
			if s.connectLost != nil {
				s.connectLost(err)
			}
			if s.logger != nil {
				s.logger.Info("connect lost")
			}
		}).
		SetOnConnectHandler(func(client mqtt.Client) {
			atomic.StoreUint32(&s.status, hubConnected)
			client.Subscribe(configChange, byte(0), func(client mqtt.Client, i mqtt.Message) {
				var msg message
				t, err := msg.parseConfigType(i.Topic())
				if err != nil {
					if s.logger != nil {
						s.logger.Warn("connect lost")
					}
					return
				}
				if s.configChange != nil {
					s.configChange(t, i.Payload())
				}
			})
			if s.logger != nil {
				s.logger.Info("connect success")
			}
		})
	client := mqtt.NewClient(options)
	s.connect(client) //reconnected
	s.client = client
}
func (s *session) getDeviceId() string {
	return s.deviceId
}
func (s *session) getThingId() string {
	return s.thingId
}
func (s *session) connect(client mqtt.Client) {
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			if s.logger != nil {
				s.logger.Info("connect retry...")
			}
			time.Sleep(3 * time.Second)
			continue
		} else {
			atomic.StoreUint32(&s.status, hubConnected)
			return
		}
	}
}
func (s *session) subscribe(topic string, call messageArrived) error {
	if atomic.LoadUint32(&s.status) == 0 {
		return notConnected
	}
	token := s.client.Subscribe(topic, byte(0), func(client mqtt.Client, message mqtt.Message) {
		call(message.Topic(), message.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (s *session) subscribes(topics []string, call messageArrived) error {
	if atomic.LoadUint32(&s.status) == 0 {
		return notConnected
	}
	filters := make(map[string]byte)
	for _, v := range topics {
		filters[v] = byte(0)
	}
	token := s.client.SubscribeMultiple(filters, func(client mqtt.Client, message mqtt.Message) {
		call(message.Topic(), message.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (s *session) setConnectLost(connectLost ConnectLost) {
	s.connectLost = connectLost
}
func (s *session) setConfigChange(configChange ConfigChangeFunc) {
	s.configChange = configChange
}

func (s *session) publish(topic string, payload []byte) error {
	if atomic.LoadUint32(&s.status) == 0 {
		return notConnected
	}
	if token := s.client.Publish(topic, byte(0), false, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (s *session) getEdgeInfo() (*edgeDevInfo, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		response *edgeDevInfo
		result   map[string]string
	)
	response = &edgeDevInfo{}
	resp, err = s.metadataClient.Get(edgeInfoRequest)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		s.logger.Error(string(content))
		return response, err
	}
	for k, v := range result {
		switch {
		case strings.Contains(k, "device_id"):
			response.Id = v
		case strings.Contains(k, "edge_id"):
		case strings.Contains(k, "thing_id"):
			response.Id = v
		case strings.Contains(k, "edge_version"):
		case strings.Contains(k, "user_id"):
		}
	}
	return response, err
}
func (s *session) getConfig() ([]*SubDeviceInfo, error) {
	var (
		err        error
		resp       *http.Response
		content    []byte
		result     driverResult
		response   []*SubDeviceInfo
		subDevices map[string]device
		//temp map[string]string
	)
	//temp = make(map[string]string)
	resp, err = s.metadataClient.Get(edgeDriverRequest + os.Getenv("DRIVER_ID"))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if err = json.Unmarshal(content, &result); err != nil {
		s.logger.Error("getConfig Unmarshal:", err.Error())
		return response, err
	}
	if subDevices, err = s.getSubDevices(); err != nil {
		s.logger.Error("getSubDevices Unmarshal:", err.Error())
		return response, err
	}
	for _, v := range result.Channels {
		if val, ok := subDevices[v.SubDeviceId]; ok {
			channelConfig := make(map[string]interface{})
			if err = json.Unmarshal([]byte(v.ChannelCfg), &channelConfig); err != nil {
				continue
			}
			deviceConfig := make(map[string]interface{})
			if err = json.Unmarshal([]byte(v.SubDeviceCfg), &deviceConfig); err != nil {
				continue
			}
			dev := &SubDeviceInfo{
				Token:       val.TokenContent,
				TokenStatus: val.TokenStatus,
				DeviceId:    val.DeviceId,
				Ext:         deviceConfig,
				ChannelCfg:  channelConfig,
			}
			response = append(response, dev)
		}
	}
	return response, err
}
func (s *session) getSubDevices() (map[string]device, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		response map[string]device
		result   map[string]string
	)
	response = make(map[string]device, 0)
	resp, err = s.metadataClient.Get(subDeviceRequest)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	s.logger.Info("[getSubDevices] ", string(content))
	err = json.Unmarshal(content, &result)
	if err != nil {
		s.logger.Error("[getSubDevices] Unmarshal err:", err.Error())
		return response, err
	}
	for _, v := range result {
		dev := device{}
		err = json.Unmarshal([]byte(v), &dev)
		if err != nil {
			s.logger.Error("[getSubDevices] result Unmarshal err:", err.Error())
			return response, err
		}
		response[dev.DeviceId] = dev
	}
	return response, err
}
func (s *session) getModel(id string) ([]byte, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		response []byte
	)
	resp, err = s.metadataClient.Get(subDeviceRequest + id)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	s.logger.Info(string(content))
	//todo need fix
	err = json.Unmarshal(content, &response)
	return content, err
}
func (s *session) getDriver() (string, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		result   *driverResult
		response string
	)
	//response = Metadata{}
	resp, err = s.metadataClient.Get(edgeDriverRequest + os.Getenv("DRIVER_ID"))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	result = &driverResult{}
	err = json.Unmarshal(content, result)
	if err != nil {
		s.logger.Error("getDriver:", string(content), err.Error())
		return response, err
	}
	return result.DriverCfg, err
}
func (s *session) disconnect() {
	if s.client != nil {
		s.client.Disconnect(250)
		atomic.StoreUint32(&s.status, hubConnected)
		s.connectLost = nil
	}
	if s.metadataClient != nil {
		s.metadataClient.CloseIdleConnections()
	}
}
