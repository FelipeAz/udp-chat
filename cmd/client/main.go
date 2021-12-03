package main

import (
	"udp-chat/infra/logger"
	"udp-chat/internal/app/chat/client"
)

func main() {
	loggerService := logger.NewLogger("log")
	cli := client.NewClient(loggerService)
	cli.Listen(":8080")
}
