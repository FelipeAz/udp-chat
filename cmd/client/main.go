package main

import (
	"bufio"
	"fmt"
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

	loggerService := logger.NewLogger("log/client", ServiceName)
	cli := client.NewClient(username, loggerService)
	cli.Listen(":8000")
}
