package client

import (
	"context"
	"io"
	"net"
	"time"
)

const maxBufferSize = 1024
const timeout = 5

func ChatClient(ctx context.Context, address string, reader io.Reader) (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}

	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	doneChan := make(chan error, 1)

	go func() {
		_, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}

		deadline := time.Now().Add(timeout * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
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
