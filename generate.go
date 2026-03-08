package main

import (
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf"
)

const (
	AppName    = "Print Square"
	AppVersion = "1.0.0"
	AppID      = "com.tsprintables.printsquare"
	AppAuthor  = "TS Printables"
)

type GridStyle string

const (
	StyleNone     GridStyle = "none"
	StyleTianZiGe GridStyle = "tianzige"
	StyleMiZiGe   GridStyle = "mizige"
)

type GridConfig struct {
	BoxW   float64   // cm
	BoxH   float64   // cm
	HGap   float64   // cm
	VGap   float64   // cm
	Margin float64   // cm
	Pages  int
	Header string
	Footer string
	Style  GridStyle
}

type GridInfo struct {
	Cols  int
	Rows  int
	Pages int
}

func GeneratePDF(cfg GridConfig, w io.Writer) (GridInfo, error) {
	// cm → mm
	bw := cfg.BoxW * 10
	bh := cfg.BoxH * 10
	hg := cfg.HGap * 10
	vg := cfg.VGap * 10
	margin := cfg.Margin * 10

	const (
		pageW = 210.0
		pageH = 297.0
	)

	topMargin := 8.0
	bottomMargin := 8.0
	headerH := 0.0
	footerH := 0.0

	if cfg.Header != "" {
		headerH = 10.0
		topMargin = 5.0
	}
	if cfg.Footer != "" || cfg.Pages > 1 {
		footerH = 10.0
		bottomMargin = 5.0
	}

	drawableW := pageW - 2*margin
	drawableH := pageH - topMargin - headerH - bottomMargin - footerH

	colCount := int((drawableW + hg) / (bw + hg))
	if colCount < 1 {
		colCount = 1
	}
	rowCount := int((drawableH + vg) / (bh + vg))
	if rowCount < 1 {
		rowCount = 1
	}

	gridW := float64(colCount)*bw + float64(colCount-1)*hg
	startX := (pageW - gridW) / 2
	startY := topMargin + headerH

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)

	// PDF metadata
	pdf.SetTitle(AppName+" — Practice Grid", true)
	pdf.SetAuthor(AppAuthor, true)
	pdf.SetSubject("Grid practice paper for Chinese character writing", true)
	pdf.SetCreator(AppName+" v"+AppVersion, true)
	pdf.SetKeywords("grid chinese practice tianzige mizige calligraphy writing", true)

	pages := cfg.Pages
	if pages < 1 {
		pages = 1
	}

	for p := 1; p <= pages; p++ {
		pdf.AddPage()

		// Header
		if cfg.Header != "" {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(80, 80, 80)
			pdf.SetXY(margin, topMargin)
			pdf.CellFormat(pageW-2*margin, 8, cfg.Header, "", 0, "C", false, 0, "")
		}

		// Grid
		for r := 0; r < rowCount; r++ {
			y := startY + float64(r)*(bh+vg)
			for c := 0; c < colCount; c++ {
				x := startX + float64(c)*(bw+hg)

				// Guide lines first (behind the border)
				drawGuideLines(pdf, cfg.Style, x, y, bw, bh)

				// Box border (solid, black)
				pdf.SetDrawColor(0, 0, 0)
				pdf.SetLineWidth(0.25)
				pdf.Rect(x, y, bw, bh, "D")
			}
		}

		// Footer line: text (left) + page number (right)
		footerY := pageH - bottomMargin - footerH + 2
		if cfg.Footer != "" {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(80, 80, 80)
			pdf.SetXY(margin, footerY)
			pdf.CellFormat(pageW-2*margin, 8, cfg.Footer, "", 0, "L", false, 0, "")
		}
		if pages > 1 {
			pdf.SetFont("Helvetica", "", 8)
			pdf.SetTextColor(120, 120, 120)
			pdf.SetXY(margin, footerY)
			pdf.CellFormat(pageW-2*margin, 8,
				fmt.Sprintf("Page %d of %d", p, pages),
				"", 0, "R", false, 0, "")
		}
	}

	if err := pdf.Output(w); err != nil {
		return GridInfo{}, err
	}

	return GridInfo{Cols: colCount, Rows: rowCount, Pages: pages}, nil
}

func drawGuideLines(pdf *gofpdf.Fpdf, style GridStyle, x, y, w, h float64) {
	if style == StyleNone {
		return
	}

	// Light gray, thin dashed lines
	pdf.SetDrawColor(190, 190, 190)
	pdf.SetLineWidth(0.15)
	pdf.SetDashPattern([]float64{1.2, 1.2}, 0)

	cx := x + w/2
	cy := y + h/2

	// Tianzige & Mizige both have a cross
	if style == StyleTianZiGe || style == StyleMiZiGe {
		pdf.Line(x, cy, x+w, cy) // horizontal center
		pdf.Line(cx, y, cx, y+h) // vertical center
	}

	// Mizige adds diagonals
	if style == StyleMiZiGe {
		pdf.Line(x, y, x+w, y+h)   // top-left → bottom-right
		pdf.Line(x+w, y, x, y+h)   // top-right → bottom-left
	}

	// Reset dash
	pdf.SetDashPattern([]float64{}, 0)
}
