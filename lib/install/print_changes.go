package install

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func printChanges(added []AbsolutePath, removed []AbsolutePath) {
	if len(added) == 0 && len(removed) == 0 {
		log.Printlnf("No changes made")
		return
	}

	const SHOW_LIMIT = 5
	homePrefix := getHome() + string(filepath.Separator)

	for _, target := range orderAndLimitSlice(added, SHOW_LIMIT) {
		log.Printlnf(color.GreenString("+ %s"), strings.TrimPrefix(target.Str(), homePrefix))
	}
	if len(added) > SHOW_LIMIT {
		boldGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
		log.Printlnf(boldGreen("+ %d more"), len(added)-SHOW_LIMIT)
	}

	for _, target := range orderAndLimitSlice(removed, SHOW_LIMIT) {
		log.Printlnf(color.RedString("- %s"), strings.TrimPrefix(target.Str(), homePrefix))
	}
	if len(removed) > SHOW_LIMIT {
		boldRed := color.New(color.FgRed, color.Bold).SprintFunc()
		log.Printlnf(boldRed("- %d more"), len(removed)-SHOW_LIMIT)
	}
}

func orderAndLimitSlice(slice []AbsolutePath, limit int) []AbsolutePath {
	slices.SortFunc(slice, func(i, j AbsolutePath) int {
		return strings.Compare(i.Str(), j.Str())
	})
	if len(slice) > limit {
		return slice[:limit]
	}
	return slice
}

func getHome() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving home directory: %v", err)
	}
	return homedir
}
