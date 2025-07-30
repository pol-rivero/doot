package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/config"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

func SetUp(t *testing.T) {
	utils.USER_INPUT_MOCK_RESPONSE = utils.MOCK_NO_INPUT
	tempDootDir := t.TempDir()
	tempCacheDir := t.TempDir()
	tempHomeDir := t.TempDir()
	os.Setenv(common.ENV_DOOT_DIR, tempDootDir)
	os.Setenv(common.ENV_DOOT_CACHE_DIR, tempCacheDir)
	os.Setenv("HOME", tempHomeDir)
}

func SetUpFiles(t *testing.T, setUpDir []FsNode) {
	SetUp(t)
	for _, node := range setUpDir {
		createNode(sourceDir(), node)
	}
}

func sourceDir() string {
	return os.Getenv(common.ENV_DOOT_DIR)
}

func sourceDirPath() AbsolutePath {
	return NewAbsolutePath(sourceDir())
}

func cacheDir() string {
	return os.Getenv(common.ENV_DOOT_CACHE_DIR)
}

func cacheFile() string {
	return filepath.Join(cacheDir(), "doot-cache.bin")
}

func homeDir() string {
	// This is not the real home dir, it's the temp dir set in SetUp
	return os.Getenv("HOME")
}

type FsNode interface {
	GetName() string
	GetChildren() []FsNode
}

type FsFile struct {
	Name    string
	Content string
}

func File(name string) FsFile {
	return FsFile{
		Name:    name,
		Content: "dummy text for file " + name,
	}
}

func ConfigFile(config config.Config) FsFile {
	if config.TargetDir == homeDir() {
		config.TargetDir = "$HOME"
	}
	configBytes, err := toml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return FsFile{
		Name:    "config.toml",
		Content: string(configBytes),
	}
}

type FsDir struct {
	Name     string
	Children []FsNode
}

func Dir(name string, children []FsNode) FsDir {
	return FsDir{Name: name, Children: children}
}

func (f FsFile) GetName() string {
	return f.Name
}

func (f FsFile) GetChildren() []FsNode {
	return nil
}

func (f FsDir) GetName() string {
	return f.Name
}

func (f FsDir) GetChildren() []FsNode {
	return f.Children
}

func createNode(parentDir string, node FsNode) {
	switch n := node.(type) {
	case FsDir:
		createDir(parentDir, n)
	case FsFile:
		createFile(parentDir, n)
	default:
		panic("unknown FsNode type")
	}
}

func createDir(parentDir string, dir FsDir) {
	dirPath := filepath.Join(parentDir, dir.GetName())
	err := os.Mkdir(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	for _, child := range dir.GetChildren() {
		createNode(dirPath, child)
	}
}

func createFile(parentDir string, file FsFile) {
	filePath := filepath.Join(parentDir, file.GetName())
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(file.Content)
	if err != nil {
		panic(err)
	}
}

func createSymlink(parentDir string, name string, target string) {
	symlinkPath := filepath.Join(parentDir, name)
	err := os.Symlink(target, symlinkPath)
	if err != nil {
		panic(err)
	}
}

func replaceWithSymlink(parentDir string, name string, target string) {
	os.Remove(filepath.Join(parentDir, name))
	createSymlink(parentDir, name, target)
}

func createHardlink(parentDir string, name string, otherFile string) {
	hardlinkPath := filepath.Join(parentDir, name)
	err := os.Link(otherFile, hardlinkPath)
	if err != nil {
		panic(err)
	}
}
