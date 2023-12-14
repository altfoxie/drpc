//go:build !windows

package drpc

import (
	"errors"
	"fmt"
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
			// socket exists, but it has a error.
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("%w: %w", ErrConnFailed, err)
			}
			if err == nil {
				return conn, nil
			}
		}
	}

	return nil, errors.New("drpc: discord is not running")
}
