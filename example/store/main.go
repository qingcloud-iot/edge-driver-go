package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"os"
	"time"
)

func main() {
	os.Setenv("EDGE_DEVICE_ID", "iott-k1y5aPNXqY")
	os.Setenv("EDGE_THING_ID", "iotd-2bc3ba1a-b325-4b13-a665-011d2fec2724")
	for {
		err := edge_driver_go.ReportEdgeProperties(context.Background(), edge_driver_go.Metadata{"A1": "A1", "B1": "B1", "C1": "C1"})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(3 * time.Second)
	}
}
