package edge_driver_go

import "context"

//validate device thing model
type validate interface {
	validateProperties(ctx context.Context, deviceId string, metadata Metadata) error
	validateEvent(ctx context.Context, deviceId string, eventName string, metadata Metadata) error
	validateService(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error
}

//validate device thing model
type dataValidate struct {
}

func newDataValidate() validate {
	return &dataValidate{}
}

func (v *dataValidate) validateProperties(ctx context.Context, deviceId string, metadata Metadata) error {
	return nil
}
func (v *dataValidate) validateEvent(ctx context.Context, deviceId string, eventName string, metadata Metadata) error {
	return nil
}
func (v *dataValidate) validateService(ctx context.Context, deviceId string, serviceName string, metadata Metadata) error {
	return nil
}
