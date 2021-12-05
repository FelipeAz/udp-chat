package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"udp-chat/infra/logger"
	"udp-chat/internal/app/chat/client"
)

const (
	ServiceName = "Client"
)

func main() {
	var username string

	fmt.Printf("Enter your Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}
	userId := uuid.NewString()

	loggerService := logger.NewLogger("log/client", ServiceName)

	cli := client.NewClient(username, userId, loggerService)
	cli.Listen("0.0.0.0:8000")
}
