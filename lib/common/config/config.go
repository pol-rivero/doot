package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
)

type Config struct {
	// The target directory for the symlinks.
	TargetDir string `toml:"target_dir"`

	// Files and directories to ignore. Each entry is a glob pattern relative to the dotfiles directory.
	ExcludeFiles []string `toml:"exclude_files"`

	// Files and directories that are always symlinked, even if they start with a dot or match a pattern in `exclude_files`. Each entry is a glob pattern relative to the dotfiles directory.
	IncludeFiles []string `toml:"include_files"`

	// If set to true, files and directories in the root of the dotfiles directory will be prefixed with a dot. For example, `<dotfiles dir>/config/foo` will be symlinked to `~/.config/foo`.
	// This is useful if you don't want to have hidden files in the root of the dotfiles directory.
	ImplicitDot bool `toml:"implicit_dot"`

	// Top-level files and directories that won't be prefixed with a dot if `implicit_dot` is set to true. Each entry is the name of a file or directory in the root of the dotfiles directory.
	ImplicitDotIgnore []string `toml:"implicit_dot_ignore"`

	// Key-value pairs of "host name" -> "host-specific directory".
	// In the example below, <dotfiles dir>/laptop-dots/.zshrc will be symlinked to ~/.zshrc, taking precedence over <dotfiles dir>/.zshrc, if the hostname is "my-laptop".
	// If `implicit_dot` is set to true, the host-specific directories also count as top-level. For example, <dotfiles dir>/laptop-dots/config/foo will be symlinked to ~/.config/foo.
	Hosts map[string]string `toml:"hosts"`
}

func DefaultConfig() Config {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving home directory: %v", err)
	}
	return Config{
		TargetDir:         homedir,
		ExcludeFiles:      []string{"**/.*", "LICENSE", "README.md"},
		IncludeFiles:      []string{},
		ImplicitDot:       false,
		ImplicitDotIgnore: []string{},
		Hosts:             map[string]string{},
	}
}

func FromFile(path AbsolutePath) Config {
	config := DefaultConfig()
	fileContents, err := os.ReadFile(path.Str())
	if err != nil {
		log.Info("Config file not found or unaccessible, using default config")
		return config
	}
	err = toml.Unmarshal(fileContents, &config)
	if err != nil {
		log.Error("Error parsing config file: %v", err)
	}
	verifyConfig(&config)
	return config
}

func FromDotfilesDir(dotfilesDir AbsolutePath) Config {
	return FromFile(dotfilesDir.Join("doot").Join("config.toml"))
}

func verifyConfig(config *Config) {
	config.TargetDir = filepath.Clean(os.ExpandEnv(config.TargetDir))
	if !filepath.IsAbs(config.TargetDir) {
		log.Fatal("Invalid config: 'target_dir = %s', must be an absolute path", config.TargetDir)
	}
	for _, implicitDotIgnore := range config.ImplicitDotIgnore {
		if strings.ContainsRune(implicitDotIgnore, filepath.Separator) {
			topLevelDir := RelativePath(implicitDotIgnore).TopLevelDir()
			log.Fatal("Invalid config. 'implicit_dot_ignore -> %s' must be a top-level file or directory. Consider adding '%s' instead",
				implicitDotIgnore, topLevelDir)
		}
	}
}
