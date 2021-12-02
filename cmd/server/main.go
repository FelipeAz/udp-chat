package main

import (
	"context"
	"log"
	"sync"
	"udp-chat/internal/app/chat/server"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := server.ChatServer(ctx, ":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
