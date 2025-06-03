package storage

import (
	"log"
	"sync"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/pdf"
)

var (
	chunkStore = make(map[string]*pdf.Chunk)
	mu         sync.RWMutex
)

func SaveChunk(chunk *pdf.Chunk) {
	mu.Lock()
	defer mu.Unlock()
	chunkStore[chunk.ID] = chunk
}

func GetChunkByID(id string) (*pdf.Chunk, bool) {
	mu.RLock()
	defer mu.RUnlock()
	chunk, ok := chunkStore[id]
	return chunk, ok
}

func QueryChunksByMetadata(key string, value interface{}) []*pdf.Chunk {
	mu.RLock()
	defer mu.RUnlock()
	var results []*pdf.Chunk
	for _, chunk := range chunkStore {
		if v, ok := chunk.Metadata[key]; ok && v == value {
			log.Println(v)
			results = append(results, chunk)
		}
	}
	return results
}

func ListAllChunks() []*pdf.Chunk {
	mu.RLock()
	defer mu.RUnlock()
	var all []*pdf.Chunk
	for _, chunk := range chunkStore {
		all = append(all, chunk)
	}
	return all
}
