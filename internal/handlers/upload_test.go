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
)

type args struct {
	filePath     string
	expectedCode int
	expectedBody string
	normalize    bool
}

type testCase struct {
	name          string
	args          args
	expectedError error
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
				expectedBody: "Chunkscreated:1",
				normalize:    true,
			},
			expectedError: nil,
		},
		{
			name: "BadPDF",
			args: args{
				filePath:     "../../mocks/bad.pdf",
				expectedCode: http.StatusBadRequest,
				expectedBody: "The uploaded file is not a valid PDF or could not be processed.",
				normalize:    false,
			},
			expectedError: nil,
		},
		{
			name: "MultiLinePDF",
			args: args{
				filePath:     "../../mocks/multiline.pdf",
				expectedCode: http.StatusOK,
				expectedBody: "Chunkscreated:1",
				normalize:    true,
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := newFileUploadRequest("/upload", "pdf", tc.args.filePath)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if err != tc.expectedError {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
				return
			}
			rr := httptest.NewRecorder()
			UploadHandler(rr, req)
			if rr.Code != tc.args.expectedCode {
				t.Errorf("expected status %d, got %d", tc.args.expectedCode, rr.Code)
			}
			got := rr.Body.String()
			if tc.args.normalize {
				got = removeWhitespace(got)
				if got != tc.args.expectedBody {
					t.Errorf("expected %q, got: %q", tc.args.expectedBody, got)
				}
			} else {
				if !strings.Contains(got, tc.args.expectedBody) {
					t.Errorf("expected error message in response, got: %s", got)
				}
			}
		})
	}
}
