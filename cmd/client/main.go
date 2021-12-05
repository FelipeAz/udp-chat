package main

import (
	"bufio"
	"fmt"
	"os"
	"udp-chat/infra/logger"
	"udp-chat/internal/app/chat/client"
)

func main() {
	var username string
	loggerService := logger.NewLogger("log")

	fmt.Printf("Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}

	cli := client.NewClient(username, loggerService)
	cli.Listen(":8000")
}
