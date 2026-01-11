package update

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pol-rivero/doot/lib/commands/diff"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/pol-rivero/doot/lib/utils"
)

func Update() {
	dotfilesDir := common.FindDotfilesDir()
	changedFiles := diff.GetChangedFiles(dotfilesDir)

	if len(changedFiles) == 0 {
		log.Printlnf("No changes to review.")
		return
	}

	log.Printlnf("Found %d file(s) with changes:\n", len(changedFiles))

	staged := 0
	skipped := 0

	for i, file := range changedFiles {
		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Printf("─── [%d/%d] %s ", i+1, len(changedFiles), file)
		fmt.Println(color.New(color.FgCyan).Sprint("───────────────────────────────────────"))

		diff.ShowFileDiff(dotfilesDir, file)
		fmt.Println()

		response := utils.RequestInput("ySaq", "Stage this change?")

		switch response {
		case 'y':
			err := utils.RunCommand(dotfilesDir, "git", "add", file)
			if err != nil {
				log.Error("Failed to stage %s: %v", file, err)
			} else {
				staged++
				color.Green("  ✓ Staged")
			}
		case 's':
			skipped++
			color.Yellow("  ○ Skipped")
		case 'a':
			for _, f := range changedFiles[i:] {
				err := utils.RunCommand(dotfilesDir, "git", "add", f)
				if err != nil {
					log.Error("Failed to stage %s: %v", f, err)
				} else {
					staged++
				}
			}
			color.Green("  ✓ Staged all remaining files")
			goto summary
		case 'q':
			color.Yellow("  Quit")
			goto summary
		}
		fmt.Println()
	}

summary:
	fmt.Println()
	log.Printlnf("Summary: %d staged, %d skipped", staged, skipped)

	if staged > 0 {
		response := utils.RequestInput("yN", "Commit staged changes?")
		if response == 'y' {
			fmt.Print("Commit message: ")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			if msg == "" {
				msg = "Update dotfiles"
			}
			err := utils.RunCommand(dotfilesDir, "git", "commit", "-m", msg)
			if err != nil {
				log.Error("Commit failed: %v", err)
			} else {
				color.Green("Changes committed!")
			}
		}
	}
}
