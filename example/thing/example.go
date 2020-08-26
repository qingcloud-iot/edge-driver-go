package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"math/rand"
	"sync"
	"time"
)

func main() {
	tokens := []string{
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5MzI4MDEyLCJpYXQiOjE1OTc3OTIwMTIsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJmSW0iLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLWQzNTllYjdlLWU4ZTUtNDAzYi1hZTRmLWU4MmUxMDczZjBlMiIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1lbmQtdXNlci1zeXN0ZW0iLCJ0eXAiOiJJRCJ9.hOr5Dfmd_SKZkBIdBtwcL8kPu3nt4fWlTllVU8v6fQ7YDjPAfh5XyblmvoG5RdB5ZILEym7zgDXXotwRQBWEoG5ic1q6KnMhFc6dUU3TgYbm86RF5GnuQZwwc1f_cWteIjOGLIHPYRAAVd36nMFoVlJSUFXIGWXjChAY3vUrPp4",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5Nzg1NjIxLCJpYXQiOjE1OTgyNDk2MjEsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJnOHEiLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLTIwNDQwM2Q0LTU3YzQtNDIyYy04MTg4LTViNWQ1ZDY5Mjg2NSIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1ieHdjQWdRbEs4IiwidHlwIjoiSUQifQ.bQVnlETr8buXvUiY9ea7SAVOjHQgpIJPTd06G1P76bjQcMlVvcH1m5peuXwKYLeIoOe8-lkFGCQ3HyMt5mI8U0mGJOSQiYPohFMXej6hxjeovtu8xKEoFqQwprHCI7mtwe5zO1rbm6gwPOGGFp3VCXXXaLBHrdnI0RCQCsaSOO0",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5Nzg1NjcxLCJpYXQiOjE1OTgyNDk2NzEsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJnREIiLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLTkwYjE0OGJhLTgwNGMtNGMxMy04Mzc1LTIzYWU1NjJkNmRlNSIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1ieHdjQWdRbEs4IiwidHlwIjoiSUQifQ.tA9xnTWcZyZZVc7lJSBXgXj3PVO22vWWTsmPa57WWzcS8OOQwKTmg96yIqD7O__ZGvQEC_hiHZ3L5UWql7FbMKp-b3b3Ze6KWhjOJPJqgmz149Y5x679OVQPG3MaXgg4ykIEleJJPhMg41NmY9ZbRljgMJVJR_TWAglTWzC0CNs",
	}
	var wg sync.WaitGroup
	wg.Add(len(tokens))
	for _, v := range tokens {
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
			client := edge_driver_go.NewEndClient(token,
				opts...)
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
				//err = client.Offline(context.Background())
				//if err != nil {
				//	fmt.Println(err)
				//}
				//time.Sleep(2 * time.Second)
			}
		}(v)
	}
	wg.Wait()
}
