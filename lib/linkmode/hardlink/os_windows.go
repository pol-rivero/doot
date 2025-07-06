//go:build windows

package linkmode_hardlink

import (
	"github.com/pol-rivero/doot/lib/common/log"
)

type NlinkType int
type HardlinkId struct{}

func osStat(_ string) (*OsStatResult, error) {
	log.Fatal("use_hardlinks is not supported on Windows.")
	panic("log.Fatal should not return")
}
