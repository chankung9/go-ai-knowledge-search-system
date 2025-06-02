package pdf

import (
	"reflect"
	"testing"
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
