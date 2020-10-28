package edge_driver_go

//设置key value
func SetValue(key string, value []byte) error {
	return getSessionIns().setValue(key, value)
}

//获取key value
func GetValue(key string) ([]byte, error) {
	return getSessionIns().getValue(key)
}
