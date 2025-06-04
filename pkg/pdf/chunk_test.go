package pdf

import (
	"reflect"
	"testing"
	"time"
)

func TestPreprocessText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalizes whitespace and trims",
			input:    "  Hello,   world!   \nNew\tline.  ",
			expected: "Hello, world! New line.",
		},
		{
			name:     "removes control characters",
			input:    "Hello\x00World\x1F!",
			expected: "HelloWorld!",
		},
		{
			name:     "already clean text",
			input:    "Sample text.",
			expected: "Sample text.",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreprocessText(tt.input)
			if got != tt.expected {
				t.Errorf("PreprocessText() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestChunkText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "splits paragraphs",
			input:    "Para1. This is the first paragraph.\n\nPara2. This is the second.",
			expected: []string{"Para1. This is the first paragraph.", "Para2. This is the second."},
		},
		{
			name:     "handles extra whitespace",
			input:    "First.\n\n  \n\nSecond.",
			expected: []string{"First.", "Second."},
		},
		{
			name:     "single paragraph",
			input:    "Just one paragraph.",
			expected: []string{"Just one paragraph."},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChunkText(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ChunkText() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestChunkDocumentText(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		text       string
		page       int
		section    string
		wantChunks int
		wantTexts  []string
	}{
		{
			name:       "multi-paragraph",
			documentID: "doc-1",
			text:       "Header1\n\nThis is the first chunk.\n\nThis is the second chunk.",
			page:       5,
			section:    "Header1",
			wantChunks: 3,
			wantTexts:  []string{"Header1", "This is the first chunk.", "This is the second chunk."},
		},
		{
			name:       "single paragraph",
			documentID: "doc-2",
			text:       "Only one chunk here.",
			page:       1,
			section:    "Intro",
			wantChunks: 1,
			wantTexts:  []string{"Only one chunk here."},
		},
		{
			name:       "empty text",
			documentID: "doc-3",
			text:       "",
			page:       0,
			section:    "",
			wantChunks: 0,
			wantTexts:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := ChunkDocumentText(tt.documentID, tt.text, tt.page, tt.section)
			if len(chunks) != tt.wantChunks {
				t.Errorf("got %d chunks, want %d", len(chunks), tt.wantChunks)
			}
			for i, want := range tt.wantTexts {
				if i >= len(chunks) {
					t.Fatalf("missing chunk %d", i)
				}
				if chunks[i].Text != want {
					t.Errorf("chunk %d text = %q, want %q", i, chunks[i].Text, want)
				}
				if chunks[i].DocumentID != tt.documentID {
					t.Errorf("chunk %d DocumentID = %q, want %q", i, chunks[i].DocumentID, tt.documentID)
				}
				if chunks[i].PageNumber != tt.page {
					t.Errorf("chunk %d PageNumber = %d, want %d", i, chunks[i].PageNumber, tt.page)
				}
				if chunks[i].Section != tt.section {
					t.Errorf("chunk %d Section = %q, want %q", i, chunks[i].Section, tt.section)
				}
				if chunks[i].Metadata["page_number"] != tt.page {
					t.Errorf("chunk %d Metadata[page_number] = %v, want %d", i, chunks[i].Metadata["page_number"], tt.page)
				}
				if chunks[i].Metadata["section"] != tt.section {
					t.Errorf("chunk %d Metadata[section] = %v, want %q", i, chunks[i].Metadata["section"], tt.section)
				}
				if chunks[i].ID == "" {
					t.Errorf("chunk %d should have a non-empty ID", i)
				}
				if time.Since(chunks[i].CreatedAt) > time.Minute {
					t.Errorf("chunk %d CreatedAt is too old: %v", i, chunks[i].CreatedAt)
				}
			}
		})
	}
}
