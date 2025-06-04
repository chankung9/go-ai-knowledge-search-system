package vector

import (
	"strconv"
	"time"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/pdf"
)

func ChunkToVectorRecord(chunk pdf.Chunk, embedding []float32) VectorRecord {
	metadata := map[string]string{
		"document_id": chunk.DocumentID,
		"page_number": strconv.Itoa(chunk.PageNumber),
		"section":     chunk.Section,
		"created_at":  chunk.CreatedAt.Format(time.RFC3339),
	}
	// Optionally, handle chunk.Metadata (marshal as JSON, etc.)
	return VectorRecord{
		ID:       chunk.ID,
		Vector:   embedding,
		Text:     chunk.Text,
		Metadata: metadata,
	}
}
