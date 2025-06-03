package embedding

type EmbeddingAPI interface {
	// Embed takes a slice of texts and returns the corresponding embeddings.
	Embed(texts []string) ([][]float32, error)
}
