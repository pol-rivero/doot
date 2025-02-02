package utils

import (
	"os"
	"os/exec"

	"github.com/pol-rivero/doot/lib/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func PrintDiff(leftFile AbsolutePath, rightFile AbsolutePath) {
	cmd := exec.Command("diff", "-u", leftFile.Str(), rightFile.Str())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to run diff: %s", err)
	}
}
