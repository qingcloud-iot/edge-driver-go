package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"time"
)

func main() {
	opts := edge_driver_go.SetName("hexing")
	client := edge_driver_go.NewClient(opts)
	time.Sleep(2 * time.Second)
	err := client.Online(context.Background(), "iotd-b85f3264-f58e-49c6-aa4a-75de98c9c214")
	if err != nil {
		fmt.Println(err)
	}

	select {}
}
