package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	eMocks "github.com/chankung9/go-ai-knowledge-search-system/pkg/embedding/mocks"
	vMocks "github.com/chankung9/go-ai-knowledge-search-system/pkg/vector/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Test helpers ---

type args struct {
	filePath     string
	expectedCode int
	expectedBody string
	normalize    bool
}

type fields struct {
	Embedder    *eMocks.EmbeddingAPI
	VectorStore *vMocks.VectorStore
}

type testCase struct {
	name          string
	args          args
	expectedError error
	prepare       func(f *fields, args args)
}

func removeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// Helper to create a multipart form file upload request
func newFileUploadRequest(uri, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestUploadHandler(t *testing.T) {
	tests := []testCase{
		{
			name: "GoodPDF",
			args: args{
				filePath:     "../../mocks/good.pdf",
				expectedCode: http.StatusOK,
				expectedBody: "Chunksprocessed:1(vectorsstored:1)",
				normalize:    true,
			},
			prepare: func(f *fields, a args) {
				f.Embedder.On("Embed", mock.Anything).Return([][]float32{{0.1, 0.2, 0.3}}, nil).Once()
				f.VectorStore.On("Insert", mock.AnythingOfType("vector.VectorRecord")).Return(nil).Once()
			},
		},
		{
			name: "BadPDF",
			args: args{
				filePath:     "../../mocks/bad.pdf",
				expectedCode: http.StatusBadRequest,
				expectedBody: "The uploaded file is not a valid PDF or could not be processed.",
				normalize:    false,
			},
			prepare: nil,
		},
		{
			name: "MultiLinePDF",
			args: args{
				filePath:     "../../mocks/multiline.pdf",
				expectedCode: http.StatusOK,
				expectedBody: "Chunksprocessed:1(vectorsstored:1)",
				normalize:    true,
			},
			prepare: func(f *fields, a args) {
				f.Embedder.On("Embed", mock.Anything).Return([][]float32{{0.5, 0.6, 0.7}}, nil).Once()
				f.VectorStore.On("Insert", mock.AnythingOfType("vector.VectorRecord")).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				Embedder:    eMocks.NewEmbeddingAPI(t), // mockery-generated constructor
				VectorStore: vMocks.NewVectorStore(t),  // mockery-generated constructor
			}
			if tt.prepare != nil {
				tt.prepare(&f, tt.args)
			}

			req, err := newFileUploadRequest("/upload", "pdf", tt.args.filePath)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			cfg := UploadHandlerCfg{
				Embedder:    f.Embedder,
				VectorStore: f.VectorStore,
			}
			handler := UploadHandler(cfg)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.args.expectedCode {
				t.Errorf("expected status %d, got %d", tt.args.expectedCode, rr.Code)
			}
			got := rr.Body.String()
			if tt.args.normalize {
				got = removeWhitespace(got)
				if got != tt.args.expectedBody {
					t.Errorf("expected %q, got: %q", tt.args.expectedBody, got)
				}
			} else {
				if !strings.Contains(got, tt.args.expectedBody) {
					t.Errorf("expected error message in response, got: %s", got)
				}
			}

			if tt.prepare != nil {
				f.Embedder.AssertExpectations(t)
				f.VectorStore.AssertExpectations(t)
			}
		})
	}
}
