package main

import (
	"log"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	lines := []string{
		"Hello, PDF World!",
		"This is the second line.",
		"",
		"And here is a new paragraph.",
		"With another line.",
	}

	for _, line := range lines {
		pdf.Cell(0, 10, line)
		pdf.Ln(10)
	}

	if err := pdf.OutputFileAndClose("multiline.pdf"); err != nil {
		log.Fatalf("Could not create PDF: %v", err)
	}
}
