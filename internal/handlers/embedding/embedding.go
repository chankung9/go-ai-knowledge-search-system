package embedding

import (
	"fmt"
	"strings"

	"github.com/chankung9/go-ai-knowledge-search-system/cmd/server/config"
)

type EmbeddingAPI interface {
	Embed(texts []string) ([][]float32, error)
}

// Factory that creates an embedding client from config struct
func EmbeddingClientFromConfig(cfg config.EmbeddingConfig) (EmbeddingAPI, error) {
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
