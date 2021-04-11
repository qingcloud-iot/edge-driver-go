package main

import (
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
)

func main() {
	err := edge_driver_go.SetValue("time", []byte("{\"key\":\"xxx\"}"))
	if err != nil {
		fmt.Print(err)
	}
	resp, err := edge_driver_go.GetValue("time")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(string(resp))
	}
	//edge_driver_go.SetBroadcastCall(func(payload []byte) {
	//	fmt.Println(string(payload))
	//})
	//err = edge_driver_go.BroadcastReport([]byte("hello world!"))
	//if err != nil {
	//	fmt.Print(err)
	//}
	//res, err := edge_driver_go.GetDeviceModel("iotd-8d8ac8d1-b5a0-4a67-b6df-dfe270303e8c")
	//if err != nil {
	//	panic(err)
	//}
	//for _,v := range res.Properties{
	//	fmt.Println(v)
	//}
	//for _,v := range res.Events{
	//	fmt.Println(v)
	//	for _,val := range v.Output {
	//		fmt.Println(val)
	//	}
	//}
	//for _,v := range res.Services{
	//	fmt.Println(v)
	//	for _,val := range v.Input {
	//		fmt.Println(val)
	//	}
	//	for _,val := range v.Output {
	//		fmt.Println(val)
	//	}
	//}
	//time.Sleep(3 * time.Second)
}
