package security

import (
	"regexp"
	"strings"
	"unicode"
)

func sanitize(input string) string {
	cleaned := strings.TrimSpace(input)

	tagRegex := regexp.MustCompile(`<.*?>`)
	cleaned = tagRegex.ReplaceAllString(cleaned, "")

	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, cleaned)

	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	return cleaned
}
