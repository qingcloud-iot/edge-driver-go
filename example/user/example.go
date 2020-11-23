package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"sync"
	"time"
)
func main() {
	subs, err := edge_driver_go.GetConfig()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(subs))
	for _, v := range subs {
		if v.TokenStatus == edge_driver_go.Enable {
			go func(token string) {
				defer func() {
					wg.Done()
				}()
				var opts []edge_driver_go.ServerOption
				opt := edge_driver_go.SetUserServiceCall(func(data []byte) (bytes []byte, e error) {
					return []byte("success"), nil
				})
				opts = append(opts, opt)
				client, err := edge_driver_go.NewEndClient(token,
					opts...)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				for {
					err := client.Online(context.Background())
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(2 * time.Second)
					err = client.ReportUserMessage(context.Background(), []byte("Qingcloud IoT"))
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(2 * time.Second)
					//err = client.Offline(context.Background())
					//if err != nil {
					//	fmt.Println(err)
					//}
					//time.Sleep(2 * time.Second)
				}
			}(v.Token)
		}
	}
	wg.Wait()
}
