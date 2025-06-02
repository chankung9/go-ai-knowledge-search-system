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

## PDF Extraction and Testing

This project robustly extracts and normalizes text from PDF uploads. The normalization removes excess whitespace and ensures extracted text is suitable for knowledge search.

### Test Coverage

- ✅ Valid PDF extraction and normalization
- ✅ Negative cases (invalid/corrupt PDFs)
- ✅ Multi-line/paragraph PDF extraction

### Running Tests

To run all tests:

```sh
go test -v ./...
```

Continuous Integration runs these tests automatically on pull requests.