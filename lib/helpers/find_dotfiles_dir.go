package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FindDotfilesDir() string {
	dotfilesDir, err := findDotfilesDir()
	if err == nil {
		log.Printf("Using dotfiles directory: %s\n", dotfilesDir)
	} else {
		log.Fatalf("Error finding dotfiles directory: %v\n", err)
		os.Exit(1)
	}
	return dotfilesDir
}

func findDotfilesDir() (string, error) {
	// 1. Try $DOOT_DIR if defined
	if dootDir := os.Getenv("DOOT_DIR"); dootDir != "" {
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
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
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
