package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	boxW := flag.Float64("bw", 1.2, "Box width in cm")
	boxH := flag.Float64("bh", 1.2, "Box height in cm")
	hGap := flag.Float64("hgap", 0, "Horizontal gap between boxes in cm")
	vGap := flag.Float64("vgap", 0.5, "Vertical gap between rows in cm")
	output := flag.String("o", "", "Output PDF path (default: ~/Documents/berries/print-format/printme-<size>.pdf)")
	header := flag.String("header", "", "TS Printables")
	footer := flag.String("footer", "", "I can do this")
	pages := flag.Int("pages", 2, "Number of pages to generate")
	marginLR := flag.Float64("margin", 0.5, "Left/right page margin in cm")

	flag.Parse()

	// cm → mm
	bw := *boxW * 10
	bh := *boxH * 10
	hg := *hGap * 10
	vg := *vGap * 10
	margin := *marginLR * 10

	// Resolve output path
	outPath := *output
	if outPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		outPath = filepath.Join(home, "Documents/berries/print-format",
			fmt.Sprintf("printme-%.1fcm.pdf", *boxW))
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating directory: %v\n", err)
		os.Exit(1)
	}

	// A4 in mm
	const (
		pageW = 210.0
		pageH = 297.0
	)

	topMargin := 8.0
	bottomMargin := 8.0
	headerH := 0.0
	footerH := 0.0

	if *header != "" {
		headerH = 10.0
		topMargin = 5.0
	}
	if *footer != "" || *pages > 1 {
		footerH = 10.0
		bottomMargin = 5.0
	}

	drawableW := pageW - 2*margin
	drawableH := pageH - topMargin - headerH - bottomMargin - footerH

	// How many boxes fit
	colCount := int((drawableW + hg) / (bw + hg))
	if colCount < 1 {
		colCount = 1
	}
	rowCount := int((drawableH + vg) / (bh + vg))
	if rowCount < 1 {
		rowCount = 1
	}

	// Center the grid horizontally
	gridW := float64(colCount)*bw + float64(colCount-1)*hg
	startX := (pageW - gridW) / 2
	startY := topMargin + headerH

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)

	for p := 1; p <= *pages; p++ {
		pdf.AddPage()

		// Header
		if *header != "" {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(80, 80, 80)
			pdf.SetXY(margin, topMargin)
			pdf.CellFormat(pageW-2*margin, 8, *header, "", 0, "C", false, 0, "")
		}

		// Grid
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetLineWidth(0.25)

		for r := 0; r < rowCount; r++ {
			y := startY + float64(r)*(bh+vg)
			for c := 0; c < colCount; c++ {
				x := startX + float64(c)*(bw+hg)
				pdf.Rect(x, y, bw, bh, "D")
			}
		}

		// Footer line: footer text (left/center) + page number (right)
		footerY := pageH - bottomMargin - footerH + 2
		if *footer != "" {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(80, 80, 80)
			pdf.SetXY(margin, footerY)
			pdf.CellFormat(pageW-2*margin, 8, *footer, "", 0, "L", false, 0, "")
		}
		if *pages > 1 {
			pdf.SetFont("Helvetica", "", 8)
			pdf.SetTextColor(120, 120, 120)
			pdf.SetXY(margin, footerY)
			pdf.CellFormat(pageW-2*margin, 8,
				fmt.Sprintf("Page %d of %d", p, *pages),
				"", 0, "R", false, 0, "")
		}
	}

	if err := pdf.OutputFileAndClose(outPath); err != nil {
		fmt.Fprintf(os.Stderr, "error writing PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PDF saved → %s\n", outPath)
	fmt.Printf("Grid: %d cols × %d rows  |  box %.1f×%.1f cm  |  hgap %.1f cm  vgap %.1f cm  |  %d page(s)\n",
		colCount, rowCount, *boxW, *boxH, *hGap, *vGap, *pages)
}
