package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chankung9/go-ai-knowledge-search-system/pkg/storage"
)

func ListChunksHandler(w http.ResponseWriter, r *http.Request) {
	chunks := storage.ListAllChunks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunks)
}
