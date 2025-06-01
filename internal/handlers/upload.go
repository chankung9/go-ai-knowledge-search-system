package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./storage"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("[ERROR] Method not allowed: %s\n", r.Method)
		return
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to create upload directory: %v\n", err)
		return
	}

	// Parse multipart form with max 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		log.Printf("[ERROR] Error parsing multipart form: %v\n", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		log.Printf("[ERROR] Error retrieving the file: %v\n", err)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		log.Printf("[ERROR] Invalid file extension: %s\n", handler.Filename)
		return
	}

	dstPath := filepath.Join(uploadPath, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Unable to save the file: "+err.Error(), http.StatusInternalServerError)
		log.Printf("[ERROR] Unable to save the file: %v\n", err)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Failed to write file: "+err.Error(), http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to write file: %v\n", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
	log.Printf("[INFO] File uploaded successfully: %s\n", handler.Filename)
}
