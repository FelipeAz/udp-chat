package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	error_messages "udp-chat/internal/app/chat/server/constants"
	"udp-chat/internal/app/chat/server/messages"
	"udp-chat/internal/logger"
)

const maxBufferSize = 1024
const timeout = 5

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
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := s.ConnectServer(ctx, port)
		if err != nil {
			s.Logger.Error(err)
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

func (s Server) ConnectServer(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		s.Logger.Error(err)
		log.Fatal(err)
	}
	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		for {
			buffer := make([]byte, maxBufferSize)
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				s.Logger.Warn(error_messages.FailedToReadFromBuffer)
				doneChan <- err
				return
			}

			msg := bytes.NewBuffer(bytes.Trim(buffer, "\x00")).String()
			id, err := s.Message.StoreMessage(msg)
			if err != nil {
				s.Logger.Error(err)
				return
			}

			fmt.Printf("%s: %s\n", id, msg)

			deadline := time.Now().Add(timeout * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				s.Logger.Warn(error_messages.FailedToWriteDeadline)
				doneChan <- err
				return
			}

			n, err = conn.WriteTo(buffer[:n], addr)
			if err != nil {
				s.Logger.Warn(error_messages.FailedToWriteToBuffer)
				doneChan <- err
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
