package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/embedding"
)

func main() {
	// Example config: use environment variables or hardcode for demonstration
	cfg := embedding.OpenAIConfig{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  os.Getenv("OPENAI_EMBEDDING_MODEL"), // e.g., "text-embedding-ada-002"
	}

	client := embedding.NewOpenAIEmbeddingClient(cfg)

	// Example input to embed
	texts := []string{
		"Hello world. This is the first chunk.",
		"And this is the second chunk.",
	}

	embeddings, err := client.Embed(texts)
	if err != nil {
		log.Fatalf("embedding failed: %v", err)
	}

	for i, emb := range embeddings {
		fmt.Printf("Embedding for chunk %d (first 8 dims): %v\n", i, emb[:8])
	}
}
