package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	style := flag.String("style", "none", "Guide line style: none, tianzige, mizige")
	ui := flag.Bool("ui", false, "Launch native desktop GUI")
	web := flag.Bool("web", false, "Launch browser-based web UI")
	port := flag.Int("port", 8080, "Port for web UI server")

	flag.Parse()

	// Auto-detect macOS .app bundle launch → default to GUI
	exe, _ := os.Executable()
	if !*ui && !*web && strings.Contains(exe, ".app/Contents/MacOS/") {
		launchGUI()
		return
	}

	// Native GUI mode
	if *ui {
		launchGUI()
		return
	}

	// Browser-based web UI mode
	if *web {
		startServer(*port)
		return
	}

	// CLI mode
	outPath := *output
	if outPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		outPath = filepath.Join(home, "Documents", "berries", "print-format",
			fmt.Sprintf("printme-%.1fcm.pdf", *boxW))
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating directory: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	cfg := GridConfig{
		BoxW:   *boxW,
		BoxH:   *boxH,
		HGap:   *hGap,
		VGap:   *vGap,
		Margin: *marginLR,
		Pages:  *pages,
		Header: *header,
		Footer: *footer,
		Style:  GridStyle(*style),
	}

	info, err := GeneratePDF(cfg, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PDF saved → %s\n", outPath)
	fmt.Printf("Grid: %d cols × %d rows  |  box %.1f×%.1f cm  |  hgap %.1f cm  vgap %.1f cm  |  style %s  |  %d page(s)\n",
		info.Cols, info.Rows, *boxW, *boxH, *hGap, *vGap, *style, info.Pages)
}
