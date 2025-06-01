package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chankung9/go-ai-knowledge-search-system/internal/handlers"
)

func main() {
	// Register HTTP handlers
	http.HandleFunc("/upload", handlers.UploadHandler)

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
