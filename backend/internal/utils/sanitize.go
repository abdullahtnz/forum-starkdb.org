package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var sanitizer = bluemonday.UGCPolicy()

func SanitizeHTML(input string) string {
	return sanitizer.Sanitize(input)
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func FilterBadWords(content string, badWords []string) string {
	lowerContent := strings.ToLower(content)
	for _, word := range badWords {
		lowerWord := strings.ToLower(word)
		idx := 0
		for {
			i := strings.Index(lowerContent[idx:], lowerWord)
			if i == -1 {
				break
			}
			absIdx := idx + i
			replacement := strings.Repeat("*", len(word))
			content = content[:absIdx] + replacement + content[absIdx+len(word):]
			lowerContent = lowerContent[:absIdx] + replacement + lowerContent[absIdx+len(word):]
			idx = absIdx + len(replacement)
		}
	}
	return content
}
