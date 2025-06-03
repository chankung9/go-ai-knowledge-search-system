package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chankung9/go-ai-knowledge-search-system/internal/handlers"
)

func main() {
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/query-chunks", handlers.QueryChunksHandler)
	http.HandleFunc("/list-chunks", handlers.ListChunksHandler)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
