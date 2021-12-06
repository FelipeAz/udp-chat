package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	error_messages "udp-chat/internal/app/chat/client/constants"
	client_model "udp-chat/internal/app/chat/client/model"
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

func NewClient(log logger.LogInterface) Client {
	return Client{
		Logger: log,
	}
}

func (c *Client) Listen(port string) {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
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

	ctx := context.Background()
	err = c.serve(ctx, conn)
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}
}

func (c *Client) serve(ctx context.Context, conn *net.UDPConn) (err error) {
	doneChan := make(chan error, 1)
	c.registerClient(conn)
	go c.listenServer(conn)
	go c.writeServer(conn)

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func (c *Client) registerClient(conn *net.UDPConn) {
	var register client_model.Register
	register.NewClient = true
	register.UserId = uuid.NewString()

	fmt.Printf("Enter your Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		register.Username = scanner.Text()
	}

	b, err := register.GetBytes()
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}

	_, err = conn.Write(b)
	if err != nil {
		c.Logger.Error(err)
		log.Fatal(err)
	}

	c.Username = register.Username
	c.UserId = register.UserId
}

func (c *Client) listenServer(conn *net.UDPConn) {
	buffer := make([]byte, 2048)
	for {
		// Read Response from server
		_, err := conn.Read(buffer)
		if err != nil {
			c.Logger.Error(err)
		}

		b := bytes.Trim(buffer, "\x00")
		if len(b) > 0 {
			resp := bytes.NewBuffer(bytes.Trim(buffer, "\x00")).String()
			fmt.Println(resp)
			buffer = make([]byte, 2048)
		}
	}
}

func (c *Client) writeServer(conn *net.UDPConn) {
	msgId := 1
	for {
		// scanner.Scan locks process until the user types a message
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
			}

			// set a connection deadline
			deadline := time.Now().Add(timeout * time.Second)
			err = conn.SetReadDeadline(deadline)
			if err != nil {
				c.Logger.Warn(error_messages.FailedToSetReaderDeadline)
			}

			msgId++
		}
	}
}

func closeConn(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
