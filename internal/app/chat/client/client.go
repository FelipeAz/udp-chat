package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	error_messages "udp-chat/internal/app/chat/client/constants"
	"udp-chat/internal/app/chat/messages/model"
	"udp-chat/internal/logger"
)

const (
	timeout = 5
)

type Client struct {
	Username string
	UserId   string
	Logger   logger.LogInterface
}

func NewClient(username, userId string, log logger.LogInterface) Client {
	return Client{
		Username: username,
		UserId:   userId,
		Logger:   log,
	}
}

func (c Client) Listen(port string) {
	ctx := context.Background()
	err := c.ConnectClient(ctx, port)
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}
}

func (c Client) ConnectClient(ctx context.Context, address string) (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}
	defer closeConn(conn)

	msgId := 1
	serverResp := make([]byte, 512)
	doneChan := make(chan error, 1)
	go func() {
		for {
			// scanner.Scan locks process until the user types a message
			fmt.Printf("Type a message: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				msg := model.NewMessage(msgId, c.Username, c.UserId, scanner.Text())
				bmsg, err := msg.ToBytes()
				if err != nil {
					c.Logger.Error(err)
					return
				}

				// Send the client input to the server
				_, err = io.Copy(conn, strings.NewReader(string(bmsg)))
				if err != nil {
					c.Logger.Warn(error_messages.FailedToCopyFromReader)
					doneChan <- err
					return
				}

				// set a connection deadline
				deadline := time.Now().Add(timeout * time.Second)
				err = conn.SetReadDeadline(deadline)
				if err != nil {
					c.Logger.Warn(error_messages.FailedToSetReaderDeadline)
					doneChan <- err
					return
				}

				// Read Response from server
				_, err = conn.Read(serverResp)
				if err != nil {
					c.Logger.Error(err)
					doneChan <- err
					return
				}
				msgId++
				//resp := bytes.NewBuffer(bytes.Trim(serverResp, "\x00")).String()
				//fmt.Println(resp)
			}
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func closeConn(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
