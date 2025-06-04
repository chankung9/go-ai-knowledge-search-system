package pdf

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DynamicMetadata provides a flexible way to attach metadata to a chunk.
type DynamicMetadata map[string]interface{}

// Chunk represents a chunk of document text with metadata.
type Chunk struct {
	ID         string          `json:"id"`
	DocumentID string          `json:"document_id"`
	PageNumber int             `json:"page_number,omitempty"`
	Section    string          `json:"section,omitempty"`
	Text       string          `json:"text"`
	CreatedAt  time.Time       `json:"created_at"`
	Metadata   DynamicMetadata `json:"metadata,omitempty"`
}

var (
	whitespaceRegex  = regexp.MustCompile(`\s+`)
	controlCharRegex = regexp.MustCompile(`[\x00-\x08\x0B-\x1F\x7F]`)
)

// PreprocessText cleans up text for embedding and search.
func PreprocessText(input string) string {
	cleaned := whitespaceRegex.ReplaceAllString(input, " ")
	cleaned = controlCharRegex.ReplaceAllString(cleaned, "")
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

// ChunkDocumentText splits the document text into chunks and returns a slice of Chunk structs.
// You can pass documentID, and optionally page or section if available (otherwise 0/empty).
func ChunkDocumentText(documentID string, text string, page int, section string) []*Chunk {
	rawChunks := ChunkText(text)
	chunks := make([]*Chunk, 0, len(rawChunks))
	now := time.Now()
	for _, c := range rawChunks {
		chunks = append(chunks, &Chunk{
			ID:         uuid.NewString(),
			DocumentID: documentID,
			PageNumber: page,
			Section:    section,
			Text:       c,
			CreatedAt:  now,
			Metadata:   DynamicMetadata{"page_number": page, "section": section},
		})
	}
	return chunks
}
