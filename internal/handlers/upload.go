package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/embedding"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/pdf"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/storage"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/vector"
	"github.com/google/uuid"
)

// UploadHandlerCfg holds dependencies for the upload handler.
type UploadHandlerCfg struct {
	Embedder    embedding.EmbeddingAPI // e.g., OpenAI client etc.
	VectorStore vector.VectorStore     // Your SQLiteVectorStore or other
}

// UploadHandler returns an http.HandlerFunc with dependency injection via closure.
func UploadHandler(cfg UploadHandlerCfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.ServeFile(w, r, "upload.html")
			return
		}

		file, header, err := r.FormFile("pdf")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		safeFilename := filepath.Base(header.Filename)
		out, err := os.CreateTemp("", "upload-*-"+safeFilename)
		if err != nil {
			http.Error(w, "Unable to create temporary file", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		tempFile := out.Name()
		defer os.Remove(tempFile)

		if _, err = io.Copy(out, file); err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}

		text, err := ExtractTextFromPDF(tempFile)
		if err != nil {
			log.Printf("PDF extraction error: %v", err)
			http.Error(w, "The uploaded file is not a valid PDF or could not be processed. Please check your file and try again.", http.StatusBadRequest)
			return
		}

		normalizedText := NormalizePDFText(text)

		documentID := uuid.NewString()
		page := 0     // Integrate true page number if you extract by page
		section := "" // Implement section detection if needed
		chunks := pdf.ChunkDocumentText(documentID, normalizedText, page, section)

		stored := 0
		for _, chunk := range chunks {
			// Save chunk meta (optional)
			storage.SaveChunk(chunk)

			// --- Embedding and vector storage integration ---
			embeddings, err := cfg.Embedder.Embed([]string{chunk.Text})
			if err != nil {
				log.Printf("Embedding error for chunk %s: %v", chunk.ID, err)
				continue // skip this chunk or handle error as needed
			}
			if len(embeddings) == 0 || len(embeddings[0]) == 0 {
				log.Printf("Embedding API returned empty result for chunk %s, input text: %q", chunk.ID, chunk.Text)
				continue
			}
			record := vector.ChunkToVectorRecord(*chunk, embeddings[0])
			if err := cfg.VectorStore.Insert(record); err != nil {
				log.Printf("Vector store insert error for chunk %s: %v", chunk.ID, err)
				continue
			}
			stored++
			log.Printf("Saved chunk & vector: %s (DocID: %s)", chunk.ID, chunk.DocumentID)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("Chunks processed: %d (vectors stored: %d)\n", len(chunks), stored)))
	}
}
