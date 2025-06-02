package pdf

import (
	"regexp"
	"strings"
)

var (
	whitespaceRegex   = regexp.MustCompile(`\s+`)
	controlCharRegex = regexp.MustCompile(`[\x00-\x08\x0B-\x1F\x7F]`)
)

// PreprocessText cleans up text for embedding and search.
func PreprocessText(input string) string {
	// Replace multiple whitespaces with a single space
	cleaned := whitespaceRegex.ReplaceAllString(input, " ")
	// Remove non-printable/control characters (except newlines)
	cleaned = controlCharRegex.ReplaceAllString(cleaned, "")
	// Trim leading/trailing whitespace
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

// ChunkText splits text into chunks by double newline (paragraphs).
// You can adjust chunking logic as needed (e.g., by sentence or token count).
func ChunkText(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	paragraphs := strings.Split(input, "\n\n")
	var chunks []string
	for _, para := range paragraphs {
		para = PreprocessText(para)
		if para != "" {
			chunks = append(chunks, para)
		}
	}
	return chunks
}
