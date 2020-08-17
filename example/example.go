package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"math/rand"
	"time"
)

func main() {
	client := edge_driver_go.NewEdgeClient("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5MTg5NDM2LCJpYXQiOjE1OTc2NTM0MzYsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJlYk8iLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLWY0N2ZjOWMzLWZjNmYtNDhmYi1hZmNmLWMwNWUyMTJkNGYzYyIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1ieHdjQWdRbEs4IiwidHlwIjoiSUQifQ.gbUZHOxcurzZhELfRmUvCb0rxbDbf9tOfpMyfr1eDN6DsOZD9f8BRmxZ6l00L0FM_ntqaiByqClJLg54Cx0_bNpLJkbCOuW9_6wBGg_X-QDJUeVGnfeLmcncAP28T0-maqU5o5FQKp6Z4mD5HTfay99QSyEoVBfxJVXwU-fv0ws")
	time.Sleep(2 * time.Second)
	for {
		err := client.Online(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
		err = client.ReportProperties(context.Background(), edge_driver_go.Metadata{"temp": rand.Float32()})
		if err != nil {
			fmt.Println(err)
		}
		err = client.ReportEvent(context.Background(), "temperatureEvent", edge_driver_go.Metadata{"temperature": rand.Float32(), "reason": true})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
		err = client.Offline(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}
	select {}
}
