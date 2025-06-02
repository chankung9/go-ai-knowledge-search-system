package handlers

import (
	"bytes"
	"fmt"
	"strings"

	"rsc.io/pdf"
)

// ExtractTextFromPDF extracts embedded text from a PDF file using rsc.io/pdf.
func ExtractTextFromPDF(pdfPath string) (string, error) {
	doc, err := pdf.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	var buf bytes.Buffer
	numPages := doc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := doc.Page(i)
		if page.V.IsNull() {
			continue
		}
		content := page.Content()
		for _, text := range content.Text {
			buf.WriteString(text.S)
			buf.WriteRune('\n')
		}
	}
	return strings.TrimSpace(buf.String()), nil
}
