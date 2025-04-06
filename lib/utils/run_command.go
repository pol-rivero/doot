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
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), "ORIGINAL_PWD="+getOriginalPwd())

	log.Info("Running command: '%s %s' (PWD: %s)", command, strings.Join(args, " "), pwd)
	return cmd.Run()
}

func RunCommandStr(pwd AbsolutePath, commandAndArgsStr string, extraArgs ...string) error {
	commandAndArgs := strings.Fields(commandAndArgsStr)
	command := commandAndArgs[0]
	args := append(commandAndArgs[1:], extraArgs...)
	return RunCommand(pwd, command, args...)
}

func getOriginalPwd() string {
	originalPwd, err := os.Getwd()
	if err != nil {
		log.Warning("Failed to get current working directory: %v, ORIGINAL_PWD won't be set", err)
		return ""
	}
	return originalPwd
}
