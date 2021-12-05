package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"udp-chat/internal/app/chat/messages"
	error_messages "udp-chat/internal/app/chat/server/constants"
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

	conn, err := net.ListenPacket("udp", port)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}

	err = s.serve(ctx, conn)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}
}

func (s Server) serve(ctx context.Context, conn net.PacketConn) (err error) {
	doneChan := make(chan error, 1)
	go func() {
		for {
			// Receiving message from server
			// buffer contains the message received from client
			buffer := make([]byte, maxBufferSize)
			_, clientAddr, err := conn.ReadFrom(buffer)
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
			response := msgObj.GetMessageFormated()
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
			_, err = conn.WriteTo(reply, clientAddr)
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
