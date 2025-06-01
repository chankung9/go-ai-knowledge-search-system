# Go AI Knowledge Search System

## Project Structure

```
go-ai-knowledge-search-system/
├── cmd/
│   └── server/
│       └── main.go             # Server entry point
├── internal/
│   └── handlers/
│       └── upload.go           # Upload endpoint handler
│       └── extract.go          # (Placeholder) PDF extraction logic
├── storage/                    # Uploaded files (gitignored)
├── go.mod
├── .gitignore
```

## Getting Started

1. Install dependencies:
   ```
   go mod tidy
   ```

2. Run the server:
   ```
   go run cmd/server/main.go
   ```

3. Upload a PDF:
   ```
   curl -F "file=@yourfile.pdf" http://localhost:8080/upload
   ```

## Next Steps

- Implement PDF text extraction in `internal/handlers/extract.go`.
- Add semantic search and vector storage.
- Build frontend interface.
