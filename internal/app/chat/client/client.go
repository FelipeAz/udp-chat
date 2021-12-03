package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	error_messages "udp-chat/internal/app/chat/client/constants"
	"udp-chat/internal/logger"
)

const timeout = 5

type Client struct {
	Logger logger.LogInterface
}

func NewClient(log logger.LogInterface) Client {
	return Client{
		Logger: log,
	}
}

func (c Client) Listen(port string) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx := context.Background()
	go func() {
		for {
			var msg string
			_, err := fmt.Scan(&msg)
			if err != nil {
				c.Logger.Error(err)
				log.Fatal(err)
			}
			err = c.ConnectClient(ctx, port, strings.NewReader(msg))
			if err != nil {
				c.Logger.Error(err)
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()
}

func (c Client) ConnectClient(ctx context.Context, address string, reader io.Reader) (err error) {
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
	defer conn.Close()

	doneChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(conn, reader)
		if err != nil {
			c.Logger.Warn(error_messages.FailedToCopyFromReader)
			doneChan <- err
			return
		}

		deadline := time.Now().Add(timeout * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			c.Logger.Warn(error_messages.FailedToSetReaderDeadline)
			doneChan <- err
			return
		}

		if err != nil {
			doneChan <- err
			return
		}

		doneChan <- nil
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
