package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chankung9/go-ai-knowledge-search-system/internal/handlers"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/embedding"
	"github.com/chankung9/go-ai-knowledge-search-system/pkg/vector"
)

func main() {
	// === Initialize dependencies ===

	// 1. Initialize Embedding API client
	embedder, err := embedding.NewEmbeddingClient(embedding.EmbeddingConfig{
		Provider:       "openai",
		OpenAIAPIKey:   "your-openai-api-key",
		OpenAIModel:    "text-embedding-ada-002",
		OllamaEndpoint: "http://localhost:11434/api/embeddings",
		OllamaModel:    "nomic-embed-text",
	})
	if err != nil {
		panic("Failed to initialize embedding client: " + err.Error())
	}

	// 2. Initialize Vector Store (SQLite as example)
	vectorStore, err := vector.NewSQLiteVectorStore("chunks_vectors.db")
	if err != nil {
		panic("Failed to initialize vector store: " + err.Error())
	}

	// 3. Prepare UploadHandler config
	uploadCfg := handlers.UploadHandlerCfg{
		Embedder:    embedder,
		VectorStore: vectorStore,
	}

	// === Register handlers ===

	http.HandleFunc("/upload", handlers.UploadHandler(uploadCfg))
	http.HandleFunc("/query-chunks", handlers.QueryChunksHandler)
	http.HandleFunc("/list-chunks", handlers.ListChunksHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
