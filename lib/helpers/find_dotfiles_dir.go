package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/constants"
	"github.com/pol-rivero/doot/lib/log"
)

func FindDotfilesDir() string {
	dotfilesDir, err := findDotfilesDir()
	if err != nil {
		log.Fatal("Error finding dotfiles directory: %v\n", err)
	}
	if !filepath.IsAbs(dotfilesDir) {
		log.Fatal("Dotfiles directory must be an absolute path: %s\n", dotfilesDir)
	}
	log.Info("Using dotfiles directory: %s\n", dotfilesDir)
	return dotfilesDir
}

func findDotfilesDir() (string, error) {
	// 1. Try $DOOT_DIR if defined
	if dootDir := os.Getenv(constants.ENV_DOOT_DIR); dootDir != "" {
		_, err := os.Stat(dootDir)
		if err == nil {
			return dootDir, nil
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error retrieving home directory: %v", err)
	}

	// 2. Try $XDG_DATA_HOME/dotfiles (or ~/.local/share/dotfiles)
	xdgDataHome := os.Getenv(constants.ENV_XDG_DATA_HOME)
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}
	dotfilesDir := filepath.Join(xdgDataHome, "dotfiles")
	if _, err = os.Stat(dotfilesDir); err == nil {
		return dotfilesDir, nil
	}

	// 3. Try ~/.dotfiles
	dotfilesDir = filepath.Join(homeDir, ".dotfiles")
	if _, err = os.Stat(dotfilesDir); err == nil {
		return dotfilesDir, nil
	}

	return "", fmt.Errorf("none of the candidate dotfiles directories exist")
}
