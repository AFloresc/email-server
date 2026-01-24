package security

import "strings"

var attackPatterns = []string{
	"<script", "</script", "javascript:",
	"onerror=", "onload=", "alert(",
	"select ", "insert ", "update ", "delete ", "drop ",
	"union ", " or 1=1", "--", ";--", "' or '1'='1",
	"../", "..\\", "%00", "%3c", "%3e",
	"$(", "`", "|", "&&", "||",
}

func containsAttackPattern(input string) bool {
	lower := strings.ToLower(input)
	for _, pattern := range attackPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}
