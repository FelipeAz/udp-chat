package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"udp-chat/internal/app/chat/client"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			var msg string

			_, err := fmt.Scan(&msg)
			if err != nil {
				log.Fatal(err)
			}

			err = client.ChatClient(ctx, ":8080", strings.NewReader(msg))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()
}
