package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func RunCommand(pwd AbsolutePath, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = pwd.Str()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("Running command: '%s %s' (PWD: %s)", command, strings.Join(args, " "), pwd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
