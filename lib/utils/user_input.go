package utils

import (
	"fmt"
	"strings"
)

const MOCK_NO_INPUT string = "__NoInput__"

var USER_INPUT_MOCK_RESPONSE string = ""

func RequestInput(options string, format string, args ...interface{}) rune {
	suffix := fmt.Sprintf(" [%s] ", addSlashes(options))
	fmt.Printf(format+suffix, args...)
	defaultResponse := ensureLower(getFirstUpperRune(options))
	responseStr := getUserInput()

	var responseRune rune
	if responseStr == "" {
		fmt.Printf("> %c\n", defaultResponse)
		responseRune = defaultResponse
	} else {
		responseRune = ensureLower(getFirstRune(responseStr))
	}

	acceptedResponses := strings.ToLower(options)
	if !strings.ContainsRune(acceptedResponses, responseRune) {
		fmt.Printf("Invalid response: '%c', defaulting to '%c'\n", responseRune, defaultResponse)
		responseRune = defaultResponse
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
	if USER_INPUT_MOCK_RESPONSE != "" {
		if USER_INPUT_MOCK_RESPONSE == MOCK_NO_INPUT {
			panic("Test was not expecting user input")
		}
		return USER_INPUT_MOCK_RESPONSE
	}
	var response string
	fmt.Scanf("%s", &response)
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
