package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/storage"
)

// QueryChunksHandler returns chunks by metadata key/value (?key=page_number&value=1)
func QueryChunksHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}
	chunks := storage.QueryChunksByMetadata(key, value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunks)
}
