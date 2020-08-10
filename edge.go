package edge_driver_go

import (
	"context"
	"github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
	"sync/atomic"
	"time"
)

type edgeDriver struct {
	ctx      context.Context
	cancel   context.CancelFunc
	validate validate
	client   mqtt.Client //hub client
	status   uint32      //0:not connected, 1:connected
	url      string      //meta address
	logger   Logger
}

// edge sdk init
func NewClient(opt ...ServerOption) Client {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	edge := &edgeDriver{
		validate: newDataValidate(),
		logger:   opts.Logger,
		ctx:      ctx,
		cancel:   cancel,
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(opts.Broker)
	options.SetClientID(opts.Name + uuid.NewV4().String())
	options.SetUsername(opts.Name + uuid.NewV4().String())
	options.SetPassword(opts.Name)
	options.SetCleanSession(true)
	options.SetAutoReconnect(true)
	options.SetKeepAlive(30 * time.Second)
	options.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		if edge.logger != nil {
			edge.logger.Warn("edge connect lost")
		}
		//heartbeat lost
		atomic.StoreUint32(&edge.status, hubNotConnected)
	})
	options.SetOnConnectHandler(func(client mqtt.Client) {
		if edge.logger != nil {
			edge.logger.Warn("edge reconnect success")
		}
		atomic.StoreUint32(&edge.status, hubConnected)
		for _, v := range opts.Services {
			if token := edge.client.Subscribe(v, byte(0), func(client mqtt.Client, i mqtt.Message) {
				edge.call(i.Topic(), i.Payload())
			}); token.Wait() && token.Error() != nil {
				if edge.logger != nil {
					edge.logger.Warn("edge sub error")
				}
			}
		}
	})
	client := mqtt.NewClient(options)
	go edge.connect(client) //reconnected
	edge.client = client
	return edge
}
func (e *edgeDriver) call(topic string, payload []byte) {
	if e.logger != nil {
		e.logger.Warn(topic, payload)
	}
	return
}
func (e *edgeDriver) getSubDevice(device string) (string, error) {
	return "", nil
}
func (e *edgeDriver) connect(client mqtt.Client) {
	for {
		select {
		case <-e.ctx.Done():
			return
		default:
		}
		if token := client.Connect(); !token.Wait() && token.Error() != nil {
			if e.logger != nil {
				e.logger.Warn("edge connect retry......")
			}
			time.Sleep(3 * time.Second)
			continue
		} else {
			if e.logger != nil {
				e.logger.Info("edge connect success")
			}
			atomic.StoreUint32(&e.status, hubConnected)
			return
		}
	}
}

func (e *edgeDriver) GetEdgeDeviceConfig(ctx context.Context) error {
	return nil
}
func (e *edgeDriver) GetEndDeviceConfig(ctx context.Context) error {
	return nil
}
func (e *edgeDriver) Online(ctx context.Context, deviceId string) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildStatusTopic(deviceId, thingId)
	data = msg.buildHeartbeatMsg(deviceId, thingId, online)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) Offline(ctx context.Context, deviceId string) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildStatusTopic(deviceId, thingId)
	data = msg.buildHeartbeatMsg(deviceId, thingId, offline)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) ReportProperties(ctx context.Context, deviceId string, params Metadata) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if err = e.validate.validateProperties(ctx, deviceId, params); err != nil {
		return err
	}
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildPropertyTopic(deviceId, thingId)
	data = msg.buildPropertyMsg(deviceId, thingId, params)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) ReportEvents(ctx context.Context, deviceId string, eventName string, params Metadata) error {
	var (
		topic   string
		msg     message
		data    []byte
		thingId string
		err     error
	)
	if err = e.validate.validateEvents(ctx, deviceId, eventName, params); err != nil {
		return err
	}
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	if thingId, err = e.getSubDevice(deviceId); err != nil {
		return err
	}
	topic = msg.buildEventTopic(deviceId, thingId, eventName)
	data = msg.buildEventMsg(deviceId, thingId, eventName, params)
	if token := e.client.Publish(topic, byte(0), false, data); token.Wait() && token.Error() != nil {
		if token.Error() != nil && e.logger != nil {
			e.logger.Warn(token.Error())
		}
		return pubMessageError
	}
	return nil
}
func (e *edgeDriver) GetDriverInfo(ctx context.Context) (interface{}, error) {
	return nil, nil
}
func (e *edgeDriver) SetProperties(ctx context.Context, params Metadata) error {
	if atomic.LoadUint32(&e.status) == hubNotConnected {
		return notConnected
	}
	return nil
}
func (e *edgeDriver) GetProperties(ctx context.Context, properties []string) (Metadata, error) {
	return nil, nil
}
func (e *edgeDriver) Close() error {
	atomic.StoreUint32(&e.status, hubNotConnected)
	if e.client != nil {
		e.client.Disconnect(250)
	}
	e.cancel()
	return nil
}
