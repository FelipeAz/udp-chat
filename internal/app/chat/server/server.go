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

type Client struct {
	Username string
	UserId   string
	Addr     net.Addr
}

type Server struct {
	Message messages.MessageInterface
	Logger  logger.LogInterface
	Clients []*Client
}

func NewServer(message messages.MessageInterface, log logger.LogInterface) Server {
	return Server{
		Message: message,
		Logger:  log,
	}
}

func (s *Server) Listen(port string) {
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

func (s *Server) serve(ctx context.Context, conn net.PacketConn) (err error) {
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
			bmsg := bytes.Trim(buffer, "\x00")
			msgObj, err := s.Message.UnmarshalMessage(bmsg)
			if err != nil {
				s.Logger.Error(err)
				log.Fatal(err)
			}

			// Create new User
			if msgObj.NewClient {
				s.addClient(msgObj.Username, msgObj.UserId, clientAddr)
				s.showLastMessages(conn, clientAddr)
			}

			// Store last message
			if !msgObj.NewClient {
				err = s.Message.Store(&msgObj)
				if err != nil {
					s.Logger.Error(err)
					return
				}
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
			s.broadcast(conn, reply)
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func (s *Server) addClient(username, userid string, address net.Addr) {
	cli := Client{
		Username: username,
		UserId:   userid,
		Addr:     address,
	}
	s.Clients = append(s.Clients, &cli)
}

func (s *Server) broadcast(conn net.PacketConn, bmsg []byte) {
	for _, client := range s.Clients {
		_, err := conn.WriteTo(bmsg, client.Addr)
		if err != nil {
			s.Logger.Error(err)
		}
	}
}

func (s *Server) showLastMessages(conn net.PacketConn, addr net.Addr) {
	msgs, err := s.Message.Get()
	if err != nil {
		fmt.Println(err)
	}
	resp := ""
	for _, msg := range msgs {
		resp += msg.GetMessageFormated() + "\n"
	}

	_, err = conn.WriteTo([]byte(resp), addr)
	if err != nil {
		fmt.Println(err)
	}
}
