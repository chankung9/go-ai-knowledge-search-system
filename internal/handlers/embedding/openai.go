package embedding

type OpenAIConfig struct {
	APIKey string
	Model  string // e.g. "text-embedding-ada-002"
}

var _ EmbeddingAPI = (*OpenAIEmbeddingClient)(nil)

type OpenAIEmbeddingClient struct {
	config OpenAIConfig
}

func NewOpenAIEmbeddingClient(cfg OpenAIConfig) *OpenAIEmbeddingClient {
	return &OpenAIEmbeddingClient{config: cfg}
}

// Embed implements EmbeddingAPI for OpenAI
func (c *OpenAIEmbeddingClient) Embed(texts []string) ([][]float32, error) {
	// Implementation goes here (see previous example)
	return nil, nil
}
