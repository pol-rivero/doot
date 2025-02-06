package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

func FindDotfilesDir() AbsolutePath {
	dotfilesDir, err := findDotfilesDir()
	if err != nil {
		log.Fatal("Error finding dotfiles directory: %v\n", err)
	}
	if !filepath.IsAbs(dotfilesDir) {
		log.Fatal("Dotfiles directory must be an absolute path: %s\n", dotfilesDir)
	}
	log.Info("Using dotfiles directory: %s\n", dotfilesDir)
	return NewAbsolutePath(dotfilesDir)
}

func findDotfilesDir() (string, error) {
	// 1. Try $DOOT_DIR if defined
	if dootDir := os.Getenv(ENV_DOOT_DIR); dootDir != "" {
		fileInfo, err := os.Stat(dootDir)
		if err == nil && fileInfo.IsDir() {
			return dootDir, nil
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error retrieving home directory: %v", err)
	}

	// 2. Try $XDG_DATA_HOME/dotfiles (or ~/.local/share/dotfiles)
	xdgDataHome := os.Getenv(ENV_XDG_DATA_HOME)
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}
	dotfilesDir := filepath.Join(xdgDataHome, "dotfiles")
	if fileInfo, err := os.Stat(dotfilesDir); err == nil && fileInfo.IsDir() {
		return dotfilesDir, nil
	}

	// 3. Try ~/.dotfiles
	dotfilesDir = filepath.Join(homeDir, ".dotfiles")
	if fileInfo, err := os.Stat(dotfilesDir); err == nil && fileInfo.IsDir() {
		return dotfilesDir, nil
	}

	err = fmt.Errorf("none of the candidate dotfiles directories exist:\n  - $DOOT_DIR = '%s'\n  - %s\n  - %s",
		os.Getenv(ENV_DOOT_DIR),
		filepath.Join(xdgDataHome, "dotfiles"),
		dotfilesDir)
	return "", err
}
