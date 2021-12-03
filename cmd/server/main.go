package main

import (
	"log"
	"udp-chat/infra/logger"
	"udp-chat/infra/redis"
	"udp-chat/internal/app/chat/server"
	"udp-chat/internal/app/chat/server/messages"
)

func main() {
	//cache, err := redis.NewCache(
	//	os.Getenv("REDIS_HOST"),
	//	os.Getenv("REDIS_PORT"),
	//	os.Getenv("REDIS_EXPIRE"),
	//)
	cache, err := redis.NewCache(
		"localhost",
		"6380",
		"1200",
	)
	if err != nil {
		log.Fatal(err)
	}
	loggerService := logger.NewLogger("log")

	// os.Getenv("QUEUE_CACHE_LENGTH")
	message := messages.NewMessage(cache, loggerService, 20)
	cli := server.NewServer(message, loggerService)
	cli.Listen(":8080")
}
