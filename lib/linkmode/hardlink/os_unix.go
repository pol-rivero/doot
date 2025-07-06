//go:build darwin || linux || freebsd || openbsd || dragonfly || netbsd

package linkmode_hardlink

import (
	"fmt"
	"os"
	"syscall"
)

type HardlinkId struct {
	Inode uint64
	Dev   uint64
}

func osStat(path string) (*OsStatResult, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to cast info.Sys() to *syscall.Stat_t for path: %s", path)
	}
	return &OsStatResult{
		numLinks: uint64(stat.Nlink),
		hardlinkId: HardlinkId{
			Inode: stat.Ino,
			Dev:   uint64(stat.Dev),
		},
	}, nil
}
