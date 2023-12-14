//go:build !windows

package drpc

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func connect() (net.Conn, error) {
	temp := "/tmp"
	subDirs := []string{
		"",
		"app/com.discordapp.Discord",
		"snap.discord",
	}

	for _, name := range []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"} {
		if value := os.Getenv(name); value != "" {
			temp = value
			break
		}
	}

	for _, sd := range subDirs {
		for i := 0; i < 10; i++ {
			conn, err := net.DialTimeout(
				"unix",
				filepath.Join(temp, sd, "discord-ipc-"+strconv.Itoa(i)),
				time.Second*5,
			)
			if err == nil {
				return conn, nil
			}
		}
	}

	return nil, ErrConnFailed
}
