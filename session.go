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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
			client:  nil,
			status:  hubNotConnected,
			logger:  newLogger(),
			endList: make([]*endClient, 0),
		}
		_ins.init()
	})
	return _ins
}

type desc struct {
	clientDeviceId string
	clientThingId  string
}

//module api
type session struct {
	client         mqtt.Client //hub client
	metadataClient *http.Client
	driverId       string
	version        string
	deviceId       string
	thingId        string
	endList []*endClient
	status         uint32           //0:not connected, 1:connected
	connectLost    ConnectLost      //connect lost callback
	configChange   ConfigChangeFunc //config change
	logger         Logger
}

func (s *session) init() {
	var (
		hubAddress string
	)
	if os.Getenv("EDGE_HUB_HOST") == "" || os.Getenv("EDGE_HUB_PORT") == "" {
		hubAddress = hubBroker
	} else {
		hubAddress = fmt.Sprintf("tcp://%s:%s", os.Getenv("EDGE_HUB_HOST"), os.Getenv("EDGE_HUB_PORT"))
	}
	if s.driverId == "" {
		if os.Getenv("EDGE_APP_ID") == "" {
			panic(errors.New("driver id is not set,sdk can't run!"))
		} else {
			s.driverId = os.Getenv("EDGE_APP_ID")
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
	s.deviceId = os.Getenv("EDGE_DEVICE_ID")
	s.thingId = os.Getenv("EDGE_THING_ID")
	if s.deviceId == "" || s.thingId == "" {
		panic("edge device id or thing id is not set!")
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(hubAddress).
		SetClientID("edge.go." + s.driverId).
		SetUsername("edge.go." + s.driverId).
		SetPassword("edge.go." + s.driverId).
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
			for _, e := range s.endList {
				clientDeviceId := e.config.DeviceId()
				clientThingId := e.config.ThingId()
				var msg message
				if isUserDevice(clientDeviceId) {
					err := s.subscribe(msg.buildUserServiceTopic(clientDeviceId, clientThingId), e.userCall)
					if err != nil {
						if s.logger != nil {
							s.logger.Warn(fmt.Sprintf("subscribe user service topic failed: %v", err))
						}
					}
				} else {
					//end service
					err := getSessionIns().subscribe(msg.buildSetTopic(clientDeviceId, clientThingId), e.endCall)
					if err != nil {
						if s.logger != nil {
							s.logger.Warn(fmt.Sprintf("subscribe set property topic failed: %v", err))
						}
					}
					err = getSessionIns().subscribe(msg.buildGetTopic(clientDeviceId, clientThingId), e.getCall)
					if err != nil {
						if s.logger != nil {
							s.logger.Warn(fmt.Sprintf("subscribe get property topic failed: %v", err))
						}
					}
					err = getSessionIns().subscribe(fmt.Sprintf(deviceService, clientThingId, clientDeviceId, "+"), e.endCall)
					if err != nil {
						if s.logger != nil {
							s.logger.Warn(fmt.Sprintf("subscribe device service topic failed: %v", err))
						}
					}
				}
			}

			client.Subscribe(fmt.Sprintf(configChange, s.driverId), byte(0), func(client mqtt.Client, i mqtt.Message) {
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
		})
	client := mqtt.NewClient(options)
	s.connect(hubAddress, client) //reconnected
	s.client = client
}

func (s *session) contains(l []*endClient, e *endClient) bool {
	for _, a := range l {
		if a == e {
			return true
		}
	}
	return false
}

func (s *session) registerEndClient(e *endClient) error {
	if ! s.contains(s.endList, e) {
		s.endList = append(s.endList, e)
		s.logger.Info("[sdk] register end device,", e.config.DeviceId(), e.config.ThingId())
	}
	return nil
}

func (s *session) getDriverVersion() string {
	if s.version == "" {
		resp, err := s.getDriverInfo()
		if err != nil {
			return ""
		} else {
			s.version = resp.Version
			return s.version
		}
	} else {
		return s.version
	}
}

func (s *session) getDriverId() string {
	return s.driverId
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
	s.logger.Info("[sdk] subscribe topic:", topic)
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
		request  string
	)
	response = &edgeDevInfo{}
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(edgeInfoRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(edgeInfoRequest, val)
	}
	resp, err = s.metadataClient.Get(request)
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
		temp    *device
		request string
	)
	//temp = make(map[string]string)
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(edgeDriverRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(edgeDriverRequest, val)
	}
	resp, err = s.metadataClient.Get(request + s.driverId)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	//s.logger.Info(string(content))
	if err = json.Unmarshal(content, &result); err != nil {
		s.logger.Error("[sdk] getConfig Unmarshal:", err.Error())
		return response, err
	}
	for _, v := range result.Channels {
		temp, err = s.getSubDevice(v.SubDeviceId)
		if err != nil {
			if s.logger != nil {
				s.logger.Warn("[sdk] getSubDevice error:", err.Error())
			}
			continue
		}
		channelConfig := make(map[string]interface{})
		if v.ChannelCfg != "" {
			if err = json.Unmarshal([]byte(v.ChannelCfg), &channelConfig); err != nil {
				s.logger.Warn("[sdk] channel cfg decode error,", err.Error())
			}
		}
		deviceConfig := make(map[string]interface{})
		if v.SubDeviceCfg != "" {
			if err = json.Unmarshal([]byte(v.SubDeviceCfg), &deviceConfig); err != nil {
				s.logger.Warn("[sdk] sub device cfg decode error,", err.Error())
			}
		}
		//connectConfig := make(map[string]interface{})
		//if temp.ConnectInfo != "" {
		//	if err = json.Unmarshal([]byte(temp.ConnectInfo), &connectConfig); err != nil {
		//		s.logger.Warn("[sdk] sub device connect info decode error,", err.Error())
		//	}
		//}
		dev := &SubDeviceInfo{
			Token:       temp.TokenContent,
			TokenStatus: TokenStatus(temp.TokenStatus),
			DeviceId:    v.SubDeviceId,
			Ext:         deviceConfig,
			ChannelCfg:  channelConfig,
			ConnectInfo: temp.ConnectInfo,
		}
		if s.logger != nil {
			s.logger.Info(fmt.Sprintf("[sdk] getSubDevice deviceId:%s,ext:%+v,cfg:%+v", dev.DeviceId, dev.Ext, dev.ChannelCfg))
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
		request  string
	)
	response = &device{}
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		// request = fmt.Sprintf("%s/test/data/childDevice/", metadataBroker)
		request = fmt.Sprintf(subDeviceRequest, metadataBroker)
	} else {
		// request = fmt.Sprintf("%s/test/data/childDevice/", val)
		request = fmt.Sprintf(subDeviceRequest, val)
	}
	// resp, err = s.metadataClient.Get(request + id + "/get")
	resp, err = s.metadataClient.Get(request + id)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	// 单值情况
	err = json.Unmarshal(content, response)
	if err != nil {
		return nil, err
	}
	if response.TokenContent != "" {
		return response, err
	}

	// 多值情况
	kv := make(map[string]string)
	err = json.Unmarshal(content, &kv)
	if err != nil {
		return nil, err
	}
	for _, v := range kv {
		d := &device{}
		err = json.Unmarshal([]byte(v), d)
		if err != nil {
			return nil, err
		}
		if d.DeviceId == id {
			return d, nil
		}
	}
	return response, err
}
func (s *session) getModel(id string) (*ThingModel, error) {
	var (
		err      error
		resp     *http.Response
		content  []byte
		response *ThingModel
		temp     device
		request  string
	)
	response = &ThingModel{
		Properties: make(map[string]*Property, 0),
	}
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(subDeviceRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(subDeviceRequest, val)
	}
	resp, err = s.metadataClient.Get(request + id)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	//s.logger.Info(string(content))
	//todo need fix
	err = json.Unmarshal(content, &temp)
	if err != nil {
		s.logger.Error("json unmarshal error", string(content))
		return response, err
	}

	if temp.TokenContent != "" {
		// 单值情况
		for _, v := range temp.Properties {
			p := &Property{
				Name:       v.Name,
				Identifier: v.Identifier,
				Type:       DefineType(v.Type).String(),
				Define:     make(map[string]interface{}),
				Ext:        make(map[string]interface{}),
			}
			if v.Define != nil {
				if err = json.Unmarshal(v.Define, &p.Define); err != nil {
					continue
				}
			}
			if v.Ext != nil {
				if err = json.Unmarshal(v.Ext, &p.Ext); err != nil {
					continue
				}
			}
			response.Properties[v.Identifier] = p
		}
		return response, err
	} else {
		// 多值情况
		kv := make(map[string]string)
		err = json.Unmarshal(content, &kv)
		if err != nil {
			return nil, err
		}
		for _, v := range kv {
			d := &device{}
			err = json.Unmarshal([]byte(v), d)
			if err != nil {
				return nil, err
			}
			if d.DeviceId == id {
				for _, v := range d.Properties {
					p := &Property{
						Name:       v.Name,
						Identifier: v.Identifier,
						Type:       DefineType(v.Type).String(),
						Define:     make(map[string]interface{}),
						Ext:        make(map[string]interface{}),
					}
					if v.Define != nil {
						if err = json.Unmarshal(v.Define, &p.Define); err != nil {
							continue
						}
					}
					if v.Ext != nil {
						if err = json.Unmarshal(v.Ext, &p.Ext); err != nil {
							continue
						}
					}
					response.Properties[v.Identifier] = p
				}
				return response, err
			}
		}
	}
	return nil, errors.New("no ssuch thing model")
}
func (s *session) getDriver() (string, error) {
	resp, err := s.getDriverInfo()
	if err != nil {
		return "", err
	} else {
		return resp.DriverCfg, nil
	}
}
func (s *session) getDriverInfo() (*driverResult, error) {
	var (
		err     error
		resp    *http.Response
		content []byte
		result  *driverResult
		request string
	)
	//response = Metadata{}
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(edgeDriverRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(edgeDriverRequest, val)
	}
	resp, err = s.metadataClient.Get(request + s.driverId)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	result = &driverResult{}
	err = json.Unmarshal(content, result)
	if err != nil {
		s.logger.Error("[sdk] getDriver:", string(content), err.Error())
		return result, err
	}
	return result, nil
}

// support json
func (s *session) setValue(key string, value []byte) error {
	var (
		err     error
		resp    *http.Response
		request string
	)
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(storeRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(storeRequest, val)
	}
	//response = Metadata{}
	resp, err = s.metadataClient.Post(request+key, "application/json", bytes.NewBuffer(value))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
func (s *session) getValue(key string) ([]byte, error) {
	var (
		err     error
		resp    *http.Response
		request string
		content []byte
	)
	//response = Metadata{}
	if val := os.Getenv("EDGE_META_ADDRESS"); val == "" {
		request = fmt.Sprintf(storeRequest, metadataBroker)
	} else {
		request = fmt.Sprintf(storeRequest, val)
	}
	resp, err = s.metadataClient.Get(request + key)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	return content, nil
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
