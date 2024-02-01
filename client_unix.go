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
	dir := "/tmp"
	subDirs := []string{
		"",
		"app/com.discordapp.Discord",
		".flatpak/dev.vencord.Vesktop/xdg-run",
		"snap.discord",
	}

	for _, name := range []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"} {
		if value := os.Getenv(name); value != "" {
			dir = value
			break
		}
	}

	for _, sd := range subDirs {
		for i := 0; i < 9; i++ {
			conn, err := net.DialTimeout(
				"unix",
				filepath.Join(dir, sd, "discord-ipc-"+strconv.Itoa(i)),
				time.Second*5,
			)
			// socket exists, but it has a error.
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			if err == nil {
				return conn, nil
			}
		}
	}

	return nil, os.ErrNotExist
}
