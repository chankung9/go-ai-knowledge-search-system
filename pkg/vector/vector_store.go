package vector

// Vector represents a vector embedding (e.g., from OpenAI)
type Vector []float32

// VectorRecord represents a stored vector and its metadata.
type VectorRecord struct {
	ID       string            // Unique identifier (could be chunk/document ID)
	Vector   Vector            // Embedding vector
	Text     string            // Original chunk text
	Metadata map[string]string // Optional metadata (e.g., doc id, page, section)
}

// SimilarityResult is the output of a similarity search.
type SimilarityResult struct {
	Record VectorRecord
	Score  float32 // Similarity score (higher = more similar, e.g., cosine similarity)
}

// VectorStore is an interface for storing and searching vectors.
type VectorStore interface {
	// Insert inserts or updates a vector record.
	Insert(record VectorRecord) error

	// Query finds top-k most similar vectors to the query vector.
	Query(query Vector, k int) ([]SimilarityResult, error)

	// (Optional) Delete, Update, etc.
}
