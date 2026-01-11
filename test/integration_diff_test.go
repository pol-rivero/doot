package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pol-rivero/doot/lib/commands/diff"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/stretchr/testify/assert"
)

func TestDiff_GetChangedFiles_NoChanges(t *testing.T) {
	setUpGitRepo(t)

	files := diff.GetChangedFiles(sourceDirPath())
	assert.Empty(t, files)
}

func TestDiff_GetChangedFiles_UnstagedChanges(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	files := diff.GetChangedFiles(sourceDirPath())
	assert.ElementsMatch(t, []string{"file1"}, files)
}

func TestDiff_GetChangedFiles_StagedChanges(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)
	runGitCommand(t, sourceDir(), "add", "file1")

	files := diff.GetChangedFiles(sourceDirPath())
	assert.ElementsMatch(t, []string{"file1"}, files)
}

func TestDiff_GetChangedFiles_NewFile(t *testing.T) {
	setUpGitRepo(t)

	createFile(sourceDir(), FsFile{Name: "newfile", Content: "new content"})

	files := diff.GetChangedFiles(sourceDirPath())
	assert.ElementsMatch(t, []string{"newfile"}, files)
}

func TestDiff_GetChangedFiles_DeletedFile(t *testing.T) {
	setUpGitRepo(t)

	os.Remove(filepath.Join(sourceDir(), "file1"))

	files := diff.GetChangedFiles(sourceDirPath())
	assert.ElementsMatch(t, []string{"file1"}, files)
}

func TestDiff_GetChangedFiles_MultipleChanges(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified"), 0644)
	createFile(sourceDir(), FsFile{Name: "newfile", Content: "new"})
	os.Remove(filepath.Join(sourceDir(), "dir1/file3"))

	files := diff.GetChangedFiles(sourceDirPath())
	assert.Len(t, files, 3)
	assert.Contains(t, files, "file1")
	assert.Contains(t, files, "newfile")
	assert.Contains(t, files, "dir1/file3")
}

func TestDiff_GetChangedFiles_RenamedFile(t *testing.T) {
	setUpGitRepo(t)

	os.Rename(filepath.Join(sourceDir(), "file1"), filepath.Join(sourceDir(), "file1_renamed"))
	runGitCommand(t, sourceDir(), "add", "-A")

	files := diff.GetChangedFiles(sourceDirPath())
	assert.Len(t, files, 1)
}

func TestDiff_Diff_NoError(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	assert.NotPanics(t, func() {
		diff.Diff(false)
	})
}

func TestDiff_Diff_Staged(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)
	runGitCommand(t, sourceDir(), "add", "file1")

	assert.NotPanics(t, func() {
		diff.Diff(true)
	})
}

func TestDiff_Status_NoError(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	assert.NotPanics(t, func() {
		diff.Status()
	})
}

func TestDiff_ShowFileDiff_NoError(t *testing.T) {
	setUpGitRepo(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	err := diff.ShowFileDiff(sourceDirPath(), "file1")
	assert.NoError(t, err)
}

func setUpGitRepo(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ImplicitDot = false
	SetUpFiles(t, false, []FsNode{
		Dir("doot", []FsNode{
			ConfigFile(cfg),
		}),
		File("file1"),
		File("file2.txt"),
		Dir("dir1", []FsNode{
			File("file3"),
		}),
	})

	runGitCommand(t, sourceDir(), "init")
	runGitCommand(t, sourceDir(), "config", "user.email", "test@test.com")
	runGitCommand(t, sourceDir(), "config", "user.name", "Test User")
	runGitCommand(t, sourceDir(), "add", "-A")
	runGitCommand(t, sourceDir(), "commit", "-m", "Initial commit")
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, output)
	}
}
