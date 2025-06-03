package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAIConfig struct {
	APIKey string
	Model  string
}

type openAIEmbeddingClient struct {
	apiKey string
	model  string
}

var _ EmbeddingAPI = (*openAIEmbeddingClient)(nil)

func NewOpenAIEmbeddingClient(cfg OpenAIConfig) EmbeddingAPI {
	return &openAIEmbeddingClient{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

func (c *openAIEmbeddingClient) Embed(texts []string) ([][]float32, error) {
	url := "https://api.openai.com/v1/embeddings"
	type EmbeddingRequest struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}
	reqBody, _ := json.Marshal(EmbeddingRequest{
		Model: c.model,
		Input: texts,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	var parsed struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	result := make([][]float32, len(parsed.Data))
	for i, d := range parsed.Data {
		result[i] = d.Embedding
	}
	return result, nil
}
