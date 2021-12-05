package main

import (
	"log"
	"os"
	"strconv"
	"udp-chat/infra/logger"
	"udp-chat/infra/redis"
	"udp-chat/internal/app/chat/server"
	"udp-chat/internal/app/chat/server/messages"
)

const (
	ServiceName = "Server"
)

func main() {
	maxSize, err := strconv.Atoi(os.Getenv("QUEUE_CACHE_LENGTH"))
	if err != nil {
		log.Fatal(err)
	}

	cache, err := redis.NewCache(
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_EXPIRE"),
	)
	if err != nil {
		log.Fatal(err)
	}

	loggerService := logger.NewLogger(os.Getenv("SERVER_LOG_FILE_PATH"), ServiceName)
	message := messages.NewMessage(cache, loggerService, maxSize)

	cli := server.NewServer(message, loggerService)
	cli.Listen(os.Getenv("CHAT_SERVER_PORT"))
}
