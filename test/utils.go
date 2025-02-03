package test

import (
	"os"
	"path/filepath"
	"regexp"
)

func MatchRegex(s string, regex string) bool {
	regexp, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	return regexp.MatchString(s)
}

func FileExists(pathParts ...string) bool {
	_, err := os.Stat(filepath.Join(pathParts...))
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
