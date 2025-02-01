package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/constants"
)

func SetUp(t *testing.T) {
	tempDootDir := t.TempDir()
	tempCacheDir := t.TempDir()
	tempHomeDir := t.TempDir()
	os.Setenv(constants.ENV_DOOT_DIR, tempDootDir)
	os.Setenv(constants.ENV_DOOT_CACHE_DIR, tempCacheDir)
	os.Setenv("HOME", tempHomeDir)
}

func SetUpFiles(t *testing.T, setUpDir []FsNode) {
	SetUp(t)
	for _, node := range setUpDir {
		createNode(os.Getenv(constants.ENV_DOOT_DIR), node)
	}
}

func sourceDir() string {
	return os.Getenv(constants.ENV_DOOT_DIR)
}

func cacheDir() string {
	return os.Getenv(constants.ENV_DOOT_CACHE_DIR)
}

func cacheFile() string {
	return filepath.Join(cacheDir(), "doot-cache.bin")
}

func homeDir() string {
	return os.Getenv("HOME")
}

type FsNode interface {
	IsDir() bool
	GetName() string
	GetChildren() []FsNode
}

type FsFile struct {
	Name string
}

func File(name string) FsFile {
	return FsFile{Name: name}
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

func (f FsFile) IsDir() bool {
	return false
}

func (f FsDir) GetName() string {
	return f.Name
}

func (f FsDir) GetChildren() []FsNode {
	return f.Children
}

func (f FsDir) IsDir() bool {
	return true
}

func createNode(parentDir string, node FsNode) {
	if node.IsDir() {
		createDir(parentDir, node)
	} else {
		createFile(parentDir, node)
	}
}

func createDir(parentDir string, dir FsNode) {
	dirPath := filepath.Join(parentDir, dir.GetName())
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		panic(err)
	}
	for _, child := range dir.GetChildren() {
		createNode(dirPath, child)
	}
}

func createFile(parentDir string, file FsNode) {
	filePath := filepath.Join(parentDir, file.GetName())
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString("dummy text for file " + filePath)
	if err != nil {
		panic(err)
	}
}
