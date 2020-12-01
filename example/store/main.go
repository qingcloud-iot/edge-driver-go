package main

import (
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"time"
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
	edge_driver_go.SetBroadcastCall(func(payload []byte) {
		fmt.Println(string(payload))
	})
	err = edge_driver_go.BroadcastReport([]byte("hello world!"))
	if err != nil {
		fmt.Print(err)
	}
	time.Sleep(3 * time.Second)
}
