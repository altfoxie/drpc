//go:build !windows

package drpc

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func connect() (net.Conn, error) {
	temp := "/tmp"
	for _, name := range []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"} {
		if value := os.Getenv(name); value != "" {
			temp = value
			break
		}
	}

	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout(
			"unix",
			filepath.Join(temp, "discord-ipc-"+strconv.Itoa(i)),
			time.Second*5,
		)
		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("drpc: connection failed")
}
