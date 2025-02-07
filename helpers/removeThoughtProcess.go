package helpers

import (
	"regexp"
	"strings"
)

func RemoveThoughtProcess(response string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleanedText := re.ReplaceAllString(response, "")
	
	return strings.TrimSpace(cleanedText)
}