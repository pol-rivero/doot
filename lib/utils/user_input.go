package utils

import (
	"fmt"
	"strings"
)

func RequestInput(options string, format string, args ...interface{}) rune {
	suffix := fmt.Sprintf(" [%s] ", addSlashes(options))
	fmt.Printf(format+suffix, args...)
	defaultResponse := getFirstUpperRune(options)
	responseStr := getUserInput()

	var responseRune rune
	if responseStr == "" {
		fmt.Println(defaultResponse)
		responseRune = defaultResponse
	} else {
		responseRune = getFirstRune(responseStr)
	}
	responseRune = ensureLower(responseRune)

	acceptedResponses := strings.ToLower(options)
	if !strings.ContainsRune(acceptedResponses, responseRune) {
		fmt.Printf("\nInvalid response: '%c', defaulting to '%c'\n", responseRune, defaultResponse)
		responseRune = ensureLower(defaultResponse)
	}

	return responseRune
}

func addSlashes(s string) string {
	var sb strings.Builder
	for i, c := range s {
		if i > 0 {
			sb.WriteString("/")
		}
		sb.WriteRune(c)
	}
	return sb.String()
}

func getUserInput() string {
	var response string
	fmt.Scanf("%s", &response)
	response = strings.ToLower(response)
	return response
}

func getFirstUpperRune(s string) rune {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return c
		}
	}
	panic("No uppercase rune found in " + s)
}

func getFirstRune(s string) rune {
	for _, c := range s {
		return c
	}
	panic("No rune found in " + s)
}

func ensureLower(c rune) rune {
	const CASE_BIT = 'a' - 'A'
	return c | CASE_BIT
}
