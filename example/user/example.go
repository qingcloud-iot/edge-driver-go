package main

import (
	"context"
	"fmt"
	edge_driver_go "github.com/qingcloud-iot/edge-driver-go"
	"sync"
	"time"
)

func main() {
	tokens := []string{
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5Nzg2Nzc5LCJpYXQiOjE1OTgyNTA3NzksImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJnSFciLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLTNkYmI1YjU2LWExZDEtNGYzOS04NTE0LWI1NTY2MGI0YjIzMiIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1lbmQtdXNlci1zeXN0ZW0iLCJ0eXAiOiJJRCJ9.jzxTzEj2f1NjW0b5h97Eu7a27leA8eVCqZPM0ywO7_oY6H6wVToG7iPF7lJvkejpo1wVRWjNqYH-8SNZwebZzkZA4b1Uj80WVZ8z68PpT4KpR2kS1rLbfklIVtJKzDWujnj-Fa3lIuh9kLzkTNJDq3ypJOigzcOdZHpCQpN9F7U",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiZXhwIjoxNjI5MzI4MDEyLCJpYXQiOjE1OTc3OTIwMTIsImlzcyI6InN0cyIsImp0aSI6InM5aVM5UWdxS1E1bGxqZmxBdjJmSW0iLCJuYmYiOjAsIm9yZ2kiOiJpb3RkLWQzNTllYjdlLWU4ZTUtNDAzYi1hZTRmLWU4MmUxMDczZjBlMiIsIm93dXIiOiJ1c3ItQjBleFduMWciLCJzdWIiOiJzdHMiLCJ0aGlkIjoiaW90dC1lbmQtdXNlci1zeXN0ZW0iLCJ0eXAiOiJJRCJ9.hOr5Dfmd_SKZkBIdBtwcL8kPu3nt4fWlTllVU8v6fQ7YDjPAfh5XyblmvoG5RdB5ZILEym7zgDXXotwRQBWEoG5ic1q6KnMhFc6dUU3TgYbm86RF5GnuQZwwc1f_cWteIjOGLIHPYRAAVd36nMFoVlJSUFXIGWXjChAY3vUrPp4",
	}
	var wg sync.WaitGroup
	wg.Add(len(tokens))
	for _, v := range tokens {
		go func(token string) {
			defer func() {
				wg.Done()
			}()
			var opts []edge_driver_go.ServerOption
			opt := edge_driver_go.SetUserServiceCall(func(data []byte) (bytes []byte, e error) {
				return []byte("success"), nil
			})
			opts = append(opts, opt)
			client := edge_driver_go.NewEndClient(token,
				opts...)
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
		}(v)
	}
	wg.Wait()
}
