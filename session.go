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
	"errors"
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
	_ins  *session
	_once sync.Once
)

func getSessionIns() *session {
	_once.Do(func() {
		_ins = &session{
			client: nil,
			status: hubNotConnected,
			logger: newLogger(),
		}
		_ins.init()
	})
	return _ins
}

//module api
type session struct {
	//lock           sync.RWMutex
	//subDevices     map[string]Client //sub device
	client         mqtt.Client //hub client
	metadataClient *http.Client
	driverId       string
	deviceId       string
	thingId        string
	//topics         []string
	status       uint32           //0:not connected, 1:connected
	connectLost  ConnectLost      //connect lost callback
	configChange ConfigChangeFunc //config change
	logger       Logger
	//messageArrived messageArrived			//message callback
}

func (s *session) init() {
	var (
		err        error
		result     *edgeDevInfo
		hubAddress string
	)
	if os.Getenv("DRIVER_HUB_ADDRESS") == "" {
		hubAddress = hubBroker
	} else {
		hubAddress = os.Getenv("DRIVER_HUB_ADDRESS")
	}
	if s.driverId == "" {
		if os.Getenv("DRIVER_ID") == "" {
			panic(errors.New("driver id is not set,sdk can't run"))
		} else {
			s.driverId = os.Getenv("DRIVER_ID")
		}
	}
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
	options.AddBroker(hubAddress).
		SetClientID("edge.driver." + s.driverId).
		SetUsername("edge.driver." + s.driverId).
		SetPassword("edge.driver." + s.driverId).
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
	s.connect(hubAddress, client) //reconnected
	s.client = client
}
func (s *session) getDeviceId() string {
	return s.deviceId
}
func (s *session) getThingId() string {
	return s.thingId
}
func (s *session) connect(address string, client mqtt.Client) {
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			if s.logger != nil {
				s.logger.Info("[sdk] connect retry...,", address, token.Error().Error())
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
		s.logger.Error("[sdk] getEdgeInfo err:", err.Error(), string(content))
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
			response.ThingId = v
		case strings.Contains(k, "edge_version"):
		case strings.Contains(k, "user_id"):
		}
	}
	return response, err
}
func (s *session) getConfig() ([]*SubDeviceInfo, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		result   driverResult
		response []*SubDeviceInfo
		//subDevices map[string]device
		temp *device
	)
	//temp = make(map[string]string)
	resp, err = s.metadataClient.Get(edgeDriverRequest + s.driverId)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if err = json.Unmarshal(content, &result); err != nil {
		s.logger.Error("[sdk] getConfig Unmarshal:", err.Error())
		return response, err
	}
	//if subDevices, err = s.getSubDevices(); err != nil {
	//	s.logger.Error("getSubDevices Unmarshal:", err.Error())
	//	return response, err
	//}
	for _, v := range result.Channels {
		temp, err = s.getSubDevice(v.SubDeviceId)
		if err != nil {
			if s.logger != nil {
				s.logger.Warn("[sdk] getSubDevice error:", err.Error())
			}
			continue
		}
		channelConfig := make(map[string]interface{})
		if err = json.Unmarshal([]byte(v.ChannelCfg), &channelConfig); err != nil {
			continue
		}
		deviceConfig := make(map[string]interface{})
		if err = json.Unmarshal([]byte(v.SubDeviceCfg), &deviceConfig); err != nil {
			continue
		}
		dev := &SubDeviceInfo{
			Token:       temp.TokenContent,
			TokenStatus: TokenStatus(temp.TokenStatus),
			DeviceId:    temp.DeviceId,
			Ext:         deviceConfig,
			ChannelCfg:  channelConfig,
		}
		response = append(response, dev)
	}
	return response, err
}
func (s *session) getSubDevice(id string) (*device, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		response *device
	)
	response = &device{}
	resp, err = s.metadataClient.Get(subDeviceRequest + id)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	s.logger.Info("[sdk] getSubDevice ", string(content))
	err = json.Unmarshal(content, response)
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
	result = make(map[string]string)
	response = make(map[string]device)
	resp, err = s.metadataClient.Get(subDeviceRequest)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	s.logger.Info("[sdk] getSubDevices ", string(content))
	err = json.Unmarshal(content, &result)
	if err != nil {
		//s.logger.Error("[getSubDevices] Unmarshal err:", err.Error())
		return response, err
	}
	for _, v := range result {
		dev := device{}
		err = json.Unmarshal([]byte(v), &dev)
		if err != nil {
			s.logger.Error("[sdk] getSubDevices result Unmarshal err:", err.Error())
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
	resp, err = s.metadataClient.Get(edgeDriverRequest + s.driverId)
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
		s.logger.Error("[sdk] getDriver:", string(content), err.Error())
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
