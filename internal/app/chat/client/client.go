package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	error_messages "udp-chat/internal/app/chat/client/constants"
	"udp-chat/internal/app/chat/entity"
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

func NewClient(username string, log logger.LogInterface) Client {
	return Client{
		Username: username,
		UserId:   uuid.NewString(),
		Logger:   log,
	}
}

func (c Client) Listen(port string) {
	ctx := context.Background()
	for {
		fmt.Printf("Type a message: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			msg := entity.Message{
				Id:       uuid.NewString(),
				UserId:   c.UserId,
				Username: c.Username,
				Text:     scanner.Text(),
				Date:     time.Time{},
			}

			bmsg, err := json.Marshal(msg)
			if err != nil {
				c.Logger.Error(err)
				return
			}

			smsg := string(bmsg)
			err = c.ConnectClient(ctx, port, strings.NewReader(smsg))
			if err != nil {
				c.Logger.Error(err)
				log.Fatal(err)
			}
		}
	}
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
	defer closeConn(conn)

	serverResp := make([]byte, 512)
	doneChan := make(chan error, 1)
	go func() {
		for {
			// Send the client input to the server
			_, err := io.Copy(conn, reader)
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
			resp := bytes.NewBuffer(bytes.Trim(serverResp, "\x00")).String()
			fmt.Println(resp)

			doneChan <- nil
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
