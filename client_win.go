//go:build windows

package drpc

import (
	"errors"
	"net"
	"strconv"
	"time"

	npipe "gopkg.in/natefinch/npipe.v2"
)

func connect() (net.Conn, error) {
	for i := 0; i < 10; i++ {
		conn, err := npipe.DialTimeout(
			`\\.\pipe\discord-ipc-`+strconv.Itoa(i),
			time.Second*5,
		)
		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("drpc: connection failed")
}
