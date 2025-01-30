package lib

import (
	"log"

	"github.com/pol-rivero/doot/lib/config"
	"github.com/pol-rivero/doot/lib/helpers"
)

func Install() {
	dotfilesDir := helpers.FindDotfilesDir()
	config := config.FromDotfilesDir(dotfilesDir)

	ignoreDootCrypt := !gitCryptIsInitialized()
	filter := helpers.CreateFilter(&config, ignoreDootCrypt)
	fileList := helpers.ScanDirectory(dotfilesDir, filter)

	log.Printf("Found %d files to symlink\n", len(fileList))
}

func gitCryptIsInitialized() bool {
	// TODO
	return false
}
