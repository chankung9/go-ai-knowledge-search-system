package embedding

type EmbeddingAPI interface {
	Embed(texts []string) ([][]float32, error)
}
