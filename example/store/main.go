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
}
