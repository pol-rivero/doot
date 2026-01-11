package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pol-rivero/doot/lib/commands/update"
	"github.com/pol-rivero/doot/lib/common/config"
	"github.com/pol-rivero/doot/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdate_NoChanges(t *testing.T) {
	setUpGitRepoForUpdate(t)

	utils.USER_INPUT_MOCK_RESPONSE = utils.MOCK_NO_INPUT

	assert.NotPanics(t, func() {
		update.Update()
	})
}

func TestUpdate_StageWithY(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	utils.USER_INPUT_MOCK_RESPONSE = "y"

	assert.NotPanics(t, func() {
		update.Update()
	})

	commits := getCommitCount(t, sourceDir())
	assert.Equal(t, 2, commits)
}

func TestUpdate_SkipWithS(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified content"), 0644)

	utils.USER_INPUT_MOCK_RESPONSE = "s"

	assert.NotPanics(t, func() {
		update.Update()
	})

	staged := getStagedFiles(t, sourceDir())
	assert.Empty(t, staged)
}

func TestUpdate_StageAllWithA(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified 1"), 0644)
	os.WriteFile(filepath.Join(sourceDir(), "file2.txt"), []byte("modified 2"), 0644)

	utils.USER_INPUT_MOCK_RESPONSE = "a"

	assert.NotPanics(t, func() {
		update.Update()
	})

	staged := getStagedFiles(t, sourceDir())
	assert.Contains(t, staged, "file1")
	assert.Contains(t, staged, "file2.txt")
}

func TestUpdate_QuitWithQ(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified 1"), 0644)
	os.WriteFile(filepath.Join(sourceDir(), "file2.txt"), []byte("modified 2"), 0644)

	utils.USER_INPUT_MOCK_RESPONSE = "q"

	assert.NotPanics(t, func() {
		update.Update()
	})

	staged := getStagedFiles(t, sourceDir())
	assert.Empty(t, staged)
}

func TestUpdate_NewFile(t *testing.T) {
	setUpGitRepoForUpdate(t)

	createFile(sourceDir(), FsFile{Name: "newfile", Content: "new content"})

	utils.USER_INPUT_MOCK_RESPONSE = "y"

	assert.NotPanics(t, func() {
		update.Update()
	})

	commits := getCommitCount(t, sourceDir())
	assert.Equal(t, 2, commits)
}

func TestUpdate_DeletedFile(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.Remove(filepath.Join(sourceDir(), "file1"))

	utils.USER_INPUT_MOCK_RESPONSE = "y"

	assert.NotPanics(t, func() {
		update.Update()
	})

	commits := getCommitCount(t, sourceDir())
	assert.Equal(t, 2, commits)
}

func TestUpdate_MultipleFilesStageFirst(t *testing.T) {
	setUpGitRepoForUpdate(t)

	os.WriteFile(filepath.Join(sourceDir(), "file1"), []byte("modified 1"), 0644)
	os.WriteFile(filepath.Join(sourceDir(), "file2.txt"), []byte("modified 2"), 0644)

	utils.USER_INPUT_MOCK_RESPONSE = "y"

	assert.NotPanics(t, func() {
		update.Update()
	})

	commits := getCommitCount(t, sourceDir())
	assert.Equal(t, 2, commits)
}

func setUpGitRepoForUpdate(t *testing.T) {
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

	runGitCommandForUpdate(t, sourceDir(), "init")
	runGitCommandForUpdate(t, sourceDir(), "config", "user.email", "test@test.com")
	runGitCommandForUpdate(t, sourceDir(), "config", "user.name", "Test User")
	runGitCommandForUpdate(t, sourceDir(), "add", "-A")
	runGitCommandForUpdate(t, sourceDir(), "commit", "-m", "Initial commit")
}

func runGitCommandForUpdate(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, output)
	}
}

func getStagedFiles(t *testing.T, dir string) []string {
	t.Helper()
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("git diff --cached failed: %v", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}
	}
	return strings.Split(output, "\n")
}

func getCommitCount(t *testing.T, dir string) int {
	t.Helper()
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("git rev-list failed: %v", err)
	}

	var count int
	_, err = fmt.Sscanf(strings.TrimSpace(out.String()), "%d", &count)
	if err != nil {
		t.Fatalf("failed to parse commit count: %v", err)
	}
	return count
}
