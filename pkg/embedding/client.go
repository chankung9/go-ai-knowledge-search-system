package embedding

import (
	"fmt"
	"strings"
)

type EmbeddingConfig struct {
	Provider       string `json:"provider"`        // e.g., "openai", "ollama"
	OpenAIAPIKey   string `json:"openai_api_key"`  // OpenAI API key
	OpenAIModel    string `json:"openai_model"`    // OpenAI embedding model, e.g., "text-embedding-ada-002"
	OllamaEndpoint string `json:"ollama_endpoint"` // Ollama API endpoint, e.g., "http://localhost:11434/api/embeddings"
	OllamaModel    string `json:"ollama_model"`    // Ollama embedding model, e.g., "nomic-embed-text"
}

// Factory that creates an embedding client from config struct
func NewEmbeddingClient(cfg EmbeddingConfig) (EmbeddingAPI, error) {
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "openai":
		return NewOpenAIEmbeddingClient(OpenAIConfig{
			APIKey: cfg.OpenAIAPIKey,
			Model:  cfg.OpenAIModel,
		}), nil
	case "ollama":
		return NewOllamaEmbeddingClient(OllamaConfig{
			Endpoint: cfg.OllamaEndpoint,
			Model:    cfg.OllamaModel,
		}), nil
	default:
		return nil, fmt.Errorf("unknown embedding provider: %q", provider)
	}
}
