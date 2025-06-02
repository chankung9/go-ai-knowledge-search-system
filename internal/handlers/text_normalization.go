package handlers

import (
	"strings"
	"unicode"
)

// NormalizePDFText attempts to "fix" extracted text that is split into single-character lines.
// If most lines are single characters, it joins all characters into a single string.
// Otherwise, it returns the original text.
func NormalizePDFText(raw string) string {
	lines := strings.Split(raw, "\n")
	singleCharLines := 0
	totalChars := 0
	for _, l := range lines {
		totalChars += len([]rune(l))
		if len([]rune(strings.TrimSpace(l))) <= 1 && len(strings.TrimSpace(l)) > 0 {
			singleCharLines++
		}
	}
	// If majority of non-empty lines are single characters, collapse
	if singleCharLines > len(lines)/2 {
		// Optionally, remove any non-printable characters
		var builder strings.Builder
		for _, l := range lines {
			c := []rune(strings.TrimSpace(l))
			if len(c) == 1 && unicode.IsPrint(c[0]) {
				builder.WriteRune(c[0])
			}
		}
		return builder.String()
	}
	return raw
}
