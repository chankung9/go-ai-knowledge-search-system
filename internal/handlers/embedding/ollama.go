package embedding

type OllamaConfig struct {
	Endpoint string // e.g., "http://localhost:11434/api/embeddings"
	Model    string // e.g., "nomic-embed-text"
}

type OllamaEmbeddingClient struct {
	config OllamaConfig
}

func NewOllamaEmbeddingClient(cfg OllamaConfig) *OllamaEmbeddingClient {
	return &OllamaEmbeddingClient{config: cfg}
}

// Embed implements EmbeddingAPI for Ollama
func (c *OllamaEmbeddingClient) Embed(texts []string) ([][]float32, error) {
	// Implementation goes here
	return nil, nil
}
