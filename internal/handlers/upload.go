package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/pdf"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/storage"
	"github.com/google/uuid"
)

// UploadHandler handles PDF uploads, text extraction, and chunking.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	// Sanitize the filename to prevent path traversal
	safeFilename := filepath.Base(header.Filename)

	// Create a temporary file with a unique name
	out, err := os.CreateTemp("", "upload-*-"+safeFilename)
	if err != nil {
		http.Error(w, "Unable to create temporary file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	tempFile := out.Name()
	defer os.Remove(tempFile)

	_, err = io.Copy(out, file)
	if err != nil {
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

	// --- New logic: chunk the text and attach metadata ---
	documentID := uuid.NewString()
	page := 0     // Integrate true page number if you extract by page
	section := "" // Implement section detection if needed
	chunks := pdf.ChunkDocumentText(documentID, normalizedText, page, section)

	// For now, we just print/log the chunk metadata for verification
	for _, chunk := range chunks {
		log.Printf("ChunkID: %s, DocID: %s, Page: %v, Section: %v, Text: %.60s, Metadata: %+v",
			chunk.ID, chunk.DocumentID, chunk.PageNumber, chunk.Section, chunk.Text, chunk.Metadata)
		storage.SaveChunk(chunk)
		log.Printf("Saved chunk: %s (DocID: %s)", chunk.ID, chunk.DocumentID)
	}

	// Optionally, respond with chunk count and preview (for debugging)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(
		"Chunks created: " +
			// Show count and a snippet for verification
			// Optionally, you can marshal to JSON and return all chunks if you wish
			// Here we return only the count for now
			fmt.Sprintf("%d\n", len(chunks)),
	))
}
