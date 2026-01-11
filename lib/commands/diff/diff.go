package diff

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
	. "github.com/pol-rivero/doot/lib/types"
)

func Diff(staged bool) {
	dotfilesDir := common.FindDotfilesDir()

	args := []string{"diff", "--color=always"}
	if staged {
		args = append(args, "--staged")
	}

	err := utils.RunCommand(dotfilesDir, "git", args...)
	if err != nil {
		log.Fatal("git diff failed: %v", err)
	}
}

func Status() {
	dotfilesDir := common.FindDotfilesDir()

	err := utils.RunCommand(dotfilesDir, "git", "status", "--short")
	if err != nil {
		log.Fatal("git status failed: %v", err)
	}
}

func GetChangedFiles(dotfilesDir AbsolutePath) []string {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dotfilesDir.Str()

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("git status failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var files []string
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		file := strings.TrimSpace(line[2:])
		if file != "" {
			files = append(files, file)
		}
	}
	return files
}

func ShowFileDiff(dotfilesDir AbsolutePath, file string) error {
	return utils.RunCommand(dotfilesDir, "git", "diff", "--color=always", "--", file)
}
