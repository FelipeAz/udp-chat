package main

import (
	"udp-chat/infra/logger"
	"udp-chat/internal/app/chat/client"
)

const (
	ServiceName = "Client"
)

func main() {
	loggerService := logger.NewLogger("log/client", ServiceName)

	cli := client.NewClient(loggerService)
	cli.Listen("0.0.0.0:8000")
}
