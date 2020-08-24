package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"time"
)

func main() {
	var opts []edge_driver_go.ServerOption
	opt := edge_driver_go.SetUserServiceCall(func(data []byte) (bytes []byte, e error) {
		return []byte("success"), nil
	})
	opts = append(opts, opt)
	client := edge_driver_go.NewEdgeClient("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5MzI4MDEyLCJpYXQiOjE1OTc3OTIwMTIsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJmSW0iLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLWQzNTllYjdlLWU4ZTUtNDAzYi1hZTRmLWU4MmUxMDczZjBlMiIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1lbmQtdXNlci1zeXN0ZW0iLCJ0eXAiOiJJRCJ9.hOr5Dfmd_SKZkBIdBtwcL8kPu3nt4fWlTllVU8v6fQ7YDjPAfh5XyblmvoG5RdB5ZILEym7zgDXXotwRQBWEoG5ic1q6KnMhFc6dUU3TgYbm86RF5GnuQZwwc1f_cWteIjOGLIHPYRAAVd36nMFoVlJSUFXIGWXjChAY3vUrPp4",
		opts...)
	for {
		err := client.Online(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
		err = client.ReportUserMessage(context.Background(), []byte{0x1, 0x2, 0x3})
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
}
