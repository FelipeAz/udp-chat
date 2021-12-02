package main

import (
	"log"
	"udp-chat/infra/logger"
	"udp-chat/infra/redis"
	"udp-chat/internal/app/chat/server"
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
	cli := server.NewServer(cache, loggerService)
	cli.Listen(":8080")
}
