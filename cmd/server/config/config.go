package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// AppConfig holds global application configuration.
type AppConfig struct {
	Embedding EmbeddingConfig
}

// EmbeddingConfig holds configuration relevant to embedding providers.
type EmbeddingConfig struct {
	Provider       string `envconfig:"EMBEDDING_PROVIDER" validate:"required,oneof=openai ollama"`
	OpenAIAPIKey   string `envconfig:"OPENAI_API_KEY" validate:"required_if=Provider openai"`
	OpenAIModel    string `envconfig:"OPENAI_EMBEDDING_MODEL" default:"text-embedding-ada-002" validate:"required_if=Provider openai"`
	OllamaEndpoint string `envconfig:"OLLAMA_EMBEDDING_ENDPOINT" default:"http://localhost:11434/api/embeddings" validate:"required_if=Provider ollama"`
	OllamaModel    string `envconfig:"OLLAMA_EMBEDDING_MODEL" default:"nomic-embed-text" validate:"required_if=Provider ollama"`
}

// LoadConfig loads the application config from environment variables and
// optionally from a .env file in cmd/server/.env using viper.
// It validates config using go-playground/validator and panics if not valid.
func LoadConfig() (*AppConfig, error) {
	envDir := "cmd/server"
	envFile := ".env"
	envPath := filepath.Join(envDir, envFile)

	v := viper.New()
	v.SetConfigFile(envPath)
	v.SetConfigType("env")
	_ = v.BindEnv("EMBEDDING_PROVIDER")
	_ = v.BindEnv("OPENAI_API_KEY")
	_ = v.BindEnv("OPENAI_EMBEDDING_MODEL")
	_ = v.BindEnv("OLLAMA_EMBEDDING_ENDPOINT")
	_ = v.BindEnv("OLLAMA_EMBEDDING_MODEL")

	if _, err := os.Stat(envPath); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// For each key, if the .env value is not empty and env system is empty, set it.
		for _, key := range []string{
			"EMBEDDING_PROVIDER",
			"OPENAI_API_KEY",
			"OPENAI_EMBEDDING_MODEL",
			"OLLAMA_EMBEDDING_ENDPOINT",
			"OLLAMA_EMBEDDING_MODEL",
		} {
			envVal := os.Getenv(key)
			fileVal := v.GetString(key)
			if envVal == "" && fileVal != "" {
				os.Setenv(key, fileVal)
			}
			// If envVal is already set (not ""), keep it (system env wins).
		}
	}

	var cfg AppConfig
	if err := envconfig.Process("", &cfg.Embedding); err != nil {
		return nil, err
	}

	validate.MustValid(&cfg.Embedding)

	return &cfg, nil
}
