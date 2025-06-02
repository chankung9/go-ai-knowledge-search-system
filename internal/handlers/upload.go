package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// UploadHandler handles PDF uploads and text extraction.
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
	// Get the path of the created temporary file for later use and cleanup
	tempFile := out.Name()
	defer os.Remove(tempFile) // Ensure the temporary file is removed after processing

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	text, err := ExtractTextFromPDF(tempFile)
	if err != nil {
		// Log the technical error for debugging
		log.Printf("PDF extraction error: %v", err)
		// Show a friendly error to the user
		http.Error(w, "The uploaded file is not a valid PDF or could not be processed. Please check your file and try again.", http.StatusBadRequest)
		return
	}

	// Normalize the extracted text before further processing or response
	normalizedText := NormalizePDFText(text)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(normalizedText))
}
