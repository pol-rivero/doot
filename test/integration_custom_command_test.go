package test

import (
	"testing"

	"github.com/pol-rivero/doot/lib"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/stretchr/testify/assert"
)

func TestCustomCmd_RunCommand(t *testing.T) {
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestCustomCmd(t, config)
	createCustomCommandFile("my-command", `#!/bin/bash
		echo "hello $PWD $@" >> output.txt`)

	lib.ExecuteRootCmd(nil, []string{"my-command", "arg1", "arg2"})
	beforeContent := readFile(sourceDir() + "/output.txt")
	assert.Equal(t, "hello "+sourceDir()+" arg1 arg2\n", beforeContent)
}

func TestCustomCmd_CommandFail(t *testing.T) {
	log.PanicInsteadOfExit = true
	config := config.DefaultConfig()
	config.ImplicitDot = false
	setUpFiles_TestCustomCmd(t, config)

	assert.Panics(t, func() {
		lib.ExecuteRootCmd(nil, []string{"some-command-not-exists", "arg1", "arg2"})
	})
}

func setUpFiles_TestCustomCmd(t *testing.T, config config.Config) {
	SetUpFiles(t, []FsNode{
		Dir("doot", []FsNode{
			ConfigFile(config),
			Dir("commands", []FsNode{}),
		}),
	})
}
