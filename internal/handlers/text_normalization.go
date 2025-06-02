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
	nonEmptyLines := 0
	for _, l := range lines {
		trimmedLine := strings.TrimSpace(l)
		if len(trimmedLine) > 0 {
			nonEmptyLines++
		}
		if len([]rune(trimmedLine)) == 1 {
			singleCharLines++
		}
	}
	// If majority of non-empty lines are single characters, collapse
	// Use nonEmptyLines for the comparison, ensuring it's not zero to avoid division by zero
	if nonEmptyLines > 0 && singleCharLines > nonEmptyLines/2 {
		// Reconstruct string from single, printable characters found on their own lines
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
