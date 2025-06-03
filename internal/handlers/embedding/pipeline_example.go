package embedding

import (
	"fmt"
	"os"
	"strings"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/pdf"
)

func EmbeddingClientFromEnv() (EmbeddingAPI, error) {
	provider := strings.ToLower(os.Getenv("EMBEDDING_PROVIDER")) // "openai" or "ollama"
	switch provider {
	case "openai":
		return NewOpenAIEmbeddingClient(OpenAIConfig{
			APIKey: os.Getenv("OPENAI_API_KEY"),
			Model:  os.Getenv("OPENAI_EMBEDDING_MODEL"), // e.g. "text-embedding-ada-002"
		}), nil
	case "ollama":
		return NewOllamaEmbeddingClient(OllamaConfig{
			Endpoint: os.Getenv("OLLAMA_EMBEDDING_ENDPOINT"), // e.g. "http://localhost:11434/api/embeddings"
			Model:    os.Getenv("OLLAMA_EMBEDDING_MODEL"),    // e.g. "nomic-embed-text"
		}), nil
	default:
		return nil, fmt.Errorf("unknown EMBEDDING_PROVIDER: %q", provider)
	}
}

// Example pipeline usage: after chunking PDF, embed and (optionally) store vectors
func ProcessDocumentChunks(documentID string, text string, page int, section string) error {
	// Chunking step (replace with your actual chunking function!)
	chunks := pdf.ChunkDocumentText(documentID, text, page, section)

	// Choose embedding client by env/config
	embedder, err := EmbeddingClientFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create embedding client: %w", err)
	}

	// Prepare slice of chunk texts for embedding
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Text
	}

	// Generate embeddings
	embeddings, err := embedder.Embed(texts)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}

	// Attach embeddings back to chunks, or upsert to vector DB here
	for i, chunk := range chunks {
		chunk.Metadata["embedding"] = embeddings[i]
		fmt.Printf("Chunk %d: Embedding %v\n", i, embeddings[i][:8]) // print first 8 dims as preview
	}

	// Optionally: Insert into vector DB here

	return nil
}
