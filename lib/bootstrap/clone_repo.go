package bootstrap

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils"
)

func CloneRepoOrExit(repo string, dotfilesDir AbsolutePath) {
	if dirExists(dotfilesDir) {
		if dirExists(dotfilesDir.Join(".git")) {
			log.Warning("Bootstrap directory (%s) seems to have already been cloned. Skipping.", dotfilesDir)
		} else {
			log.Fatal("Bootstrap directory (%s) already exists, but it's not a git repository. Aborting.", dotfilesDir)
		}
		return
	}

	repoDirName := filepath.Base(dotfilesDir.Str())
	gitUrl := getGitUrl(repo)
	log.Info("Cloning repository %s into %s (%s)", gitUrl, dotfilesDir, repoDirName)
	err := utils.RunCommand(dotfilesDir.Parent(), "git", "clone", gitUrl, repoDirName)
	if err != nil {
		log.Fatal("Failed to clone repository: %v", err)
	}
}

func dirExists(path AbsolutePath) bool {
	fileInfo, err := os.Stat(path.Str())
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("Error checking if directory %s exists: %v", path, err)
	}

	if !fileInfo.IsDir() {
		log.Fatal("Expected %s to be a directory, but it's not", path)
	}
	return true
}

func getGitUrl(repo string) string {
	if isShortGithubRepo(repo) {
		return "https://github.com/" + repo + ".git"
	}
	return repo
}

func isShortGithubRepo(repo string) bool {
	regex := regexp.MustCompile(`^[\w.-]+/[\w.-]+$`) // <user>/<repo>
	return regex.MatchString(repo)
}
