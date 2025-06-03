package vector

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteVectorStore implements VectorStore using sqlite3.
type SQLiteVectorStore struct {
	DB *sql.DB
}

var _ VectorStore = (*SQLiteVectorStore)(nil)

// NewSQLiteVectorStore opens or creates a new sqlite DB for vectors.
func NewSQLiteVectorStore(dbPath string) (*SQLiteVectorStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS embeddings (
		id TEXT PRIMARY KEY,
		vector BLOB NOT NULL,
		text TEXT,
		metadata TEXT
	);`
	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}
	return &SQLiteVectorStore{DB: db}, nil
}

// Insert inserts or replaces a vector record in the database.
func (s *SQLiteVectorStore) Insert(record VectorRecord) error {
	vecBytes, err := float32SliceToBytes(record.Vector)
	if err != nil {
		return err
	}
	metaJSON, err := json.Marshal(record.Metadata)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(
		`INSERT OR REPLACE INTO embeddings (id, vector, text, metadata) VALUES (?, ?, ?, ?)`,
		record.ID, vecBytes, record.Text, string(metaJSON),
	)
	return err
}

// Query retrieves the top-k most similar vectors to the query vector using cosine similarity.
func (s *SQLiteVectorStore) Query(query Vector, k int) ([]SimilarityResult, error) {
	rows, err := s.DB.Query(`SELECT id, vector, text, metadata FROM embeddings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SimilarityResult
	for rows.Next() {
		var id, text, metaJSON string
		var vecBytes []byte
		if err := rows.Scan(&id, &vecBytes, &text, &metaJSON); err != nil {
			return nil, err
		}
		vec, err := bytesToFloat32Slice(vecBytes)
		if err != nil {
			return nil, err
		}
		score := cosineSimilarity(query, vec)

		metadata := map[string]string{}
		if metaJSON != "" {
			json.Unmarshal([]byte(metaJSON), &metadata)
		}

		results = append(results, SimilarityResult{
			Record: VectorRecord{
				ID:       id,
				Vector:   vec,
				Text:     text,
				Metadata: metadata,
			},
			Score: score,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })

	if len(results) > k {
		results = results[:k]
	}
	return results, nil
}

// Helper: serialize []float32 to []byte (little endian)
func float32SliceToBytes(vec []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, f := range vec {
		if err := binary.Write(buf, binary.LittleEndian, f); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// Helper: deserialize []byte to []float32 (little endian)
func bytesToFloat32Slice(b []byte) ([]float32, error) {
	if len(b)%4 != 0 {
		return nil, errors.New("invalid []float32 blob length")
	}
	vec := make([]float32, len(b)/4)
	for i := range vec {
		vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4 : i*4+4]))
	}
	return vec, nil
}

// Helper: cosine similarity
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
