package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"
	error_messages "udp-chat/internal/app/chat/server/constants"
	"udp-chat/internal/app/chat/server/messages"
	"udp-chat/internal/logger"
)

const (
	maxBufferSize = 1024
	timeout       = 5
)

type Server struct {
	Message messages.MessageInterface
	Logger  logger.LogInterface
}

func NewServer(message messages.MessageInterface, log logger.LogInterface) Server {
	return Server{
		Message: message,
		Logger:  log,
	}
}

func (s Server) Listen(port string) {
	fmt.Println("CHAT IS READY FOR CONNECTION")
	ctx := context.Background()

	err := s.ConnectServer(ctx, port)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}
}

func (s Server) ConnectServer(ctx context.Context, address string) (err error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}

	doneChan := make(chan error, 1)
	go func() {
		for {
			// Receiving message from server
			// buffer contains the message received from client
			buffer := make([]byte, maxBufferSize)
			_, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				s.Logger.Warn(error_messages.FailedToReadFromBuffer)
				doneChan <- err
				return
			}

			// Storing message to cache
			msg := bytes.NewBuffer(bytes.Trim(buffer, "\x00")).String()
			msgObj, err := s.Message.Store(msg)
			if err != nil {
				s.Logger.Error(err)
				return
			}

			// Return message
			msgTime := msgObj.Date.Format("01-02-2006 03:04")
			response := fmt.Sprintf("%s %s: %s", msgTime, msgObj.Username, msgObj.Text)
			fmt.Println(response)

			// Response deadline
			deadline := time.Now().Add(timeout * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				s.Logger.Warn(error_messages.FailedToWriteDeadline)
				doneChan <- err
				return
			}

			// Writing message to client
			reply := []byte(response)
			_, err = conn.WriteToUDP(reply, addr)
			if err != nil {
				s.Logger.Error(err)
				return
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
