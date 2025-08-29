//go:build linux

package filesystememulation

import (
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"testing"
)

func PickDifferentFSDir(t *testing.T, other string) string {
	candidates := []string{
		"/dev/shm",
		"/run/shm",
		"/var/run/shm",
	}

	// Per-user runtime dir often a tmpfs.
	if u, _ := user.Current(); u != nil && u.Uid != "" {
		candidates = append(candidates, filepath.Join("/run/user", u.Uid))
	}

	for _, base := range candidates {
		if base == "" {
			continue
		}
		d, err := os.MkdirTemp(base, "dootfs-*")
		if err != nil {
			continue
		}
		// If it’s not a different device, discard and continue.
		if !sameDev(d, other) {
			return d
		}
		_ = os.RemoveAll(d)
	}
	t.Fatal("no writable different filesystem available")
	panic("unreachable")
}

func sameDev(a, b string) bool {
	return devOf(a) == devOf(b)
}

func devOf(p string) uint64 {
	fi, err := os.Stat(p)
	if err != nil {
		return 0
	}
	st := fi.Sys().(*syscall.Stat_t)
	return uint64(st.Dev)
}
