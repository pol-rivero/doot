package test

import "regexp"

func MatchRegex(s string, regex string) bool {
	regexp, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	return regexp.MatchString(s)
}
