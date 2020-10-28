package edge_driver_go

type valueData struct {
	value interface{}
	time  int64
}
type metadata struct {
	values map[string]valueData
}

func NewMetadata() MetadataMsg {
	return &metadata{values: make(map[string]valueData, 0)}
}
func (m *metadata) Add(key string, value interface{}, time int64) {
	val := valueData{
		value: value,
		time:  time,
	}
	m.values[key] = val
	return
}
func (m *metadata) Del(key string) {
	delete(m.values, key)
}
func (m *metadata) Data() Metadata {
	result := Metadata{}
	for key, _ := range m.values {
		result[key] = m.values[key]
	}
	return result
}
