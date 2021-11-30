package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"
)

const maxBufferSize = 1024
const timeout = 5

func ChatServer(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer func(conn net.PacketConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			buffer = make([]byte, maxBufferSize)
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			msg := bytes.NewBuffer(buffer)
			fmt.Printf("%s: %s\n", addr.String(), msg.String())

			deadline := time.Now().Add(timeout * time.Second)
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			n, err = conn.WriteTo(buffer[:n], addr)
			if err != nil {
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
