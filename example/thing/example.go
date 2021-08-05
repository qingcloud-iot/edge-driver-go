package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"math/rand"
	"sync"
	"time"
)

func init() {
}
func main() {
	subs, err := edge_driver_go.GetConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(len(subs))

	var wg sync.WaitGroup
	wg.Add(len(subs))
	for _, v := range subs {
		fmt.Println(v.TokenStatus)
		if v.TokenStatus == edge_driver_go.Enable {
			go func(token string) {
				defer func() {
					wg.Done()
				}()
				var opts []edge_driver_go.ServerOption
				opt := edge_driver_go.SetEndServiceCall(func(name string, args edge_driver_go.Metadata) (reply *edge_driver_go.Reply, e error) {
					fmt.Println(name, args)
					return
				})
				opts = append(opts, opt)
				opt = edge_driver_go.SetGetServiceCall(func(args []string) (metadata edge_driver_go.Metadata, e error) {
					fmt.Println(args)
					return
				})
				opts = append(opts, opt)
				opt = edge_driver_go.SetSetServiceCall(func(args edge_driver_go.Metadata) error {
					fmt.Println(args)
					return nil
				})
				opts = append(opts, opt)
				client, err := edge_driver_go.NewEndClient(token,
					opts...)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
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
					//消息体带上tag
					err = client.ReportPropertiesWithTags(context.Background(), edge_driver_go.Metadata{"temp": rand.Float32()}, edge_driver_go.Metadata{"sn": "1234567890"})
					if err != nil {
						fmt.Println(err)
					}
					//消息体带上tag和自定义时间戳（毫秒）
					msg := edge_driver_go.MetadataMsg{}
					msg["temp"] = edge_driver_go.ValueData{
						Value: rand.Float32(),
						Time:  1603866709111,
					}
					err = client.ReportPropertiesWithTagsEx(context.Background(), msg, edge_driver_go.Metadata{"sn": "1234567890"})
					if err != nil {
						fmt.Println(err)
					}
					err = client.ReportEvent(context.Background(), "temperatureEvent", edge_driver_go.Metadata{"temperature": rand.Float32(), "reason": true})
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(2 * time.Second)
				}
			}(v.Token)
		}
	}
	wg.Wait()
}
