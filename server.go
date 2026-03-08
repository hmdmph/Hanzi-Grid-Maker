package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

//go:embed ui/index.html
var indexHTML []byte

func startServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	http.HandleFunc("/icon.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(AppIconBytes)
	})

	http.HandleFunc("/generate", handleGenerate)

	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("🟥 print-square UI → %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")

	openBrowser(url)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	cfg := GridConfig{
		BoxW:   parseFloat(r.FormValue("bw"), 1.2),
		BoxH:   parseFloat(r.FormValue("bh"), 1.2),
		HGap:   parseFloat(r.FormValue("hgap"), 0),
		VGap:   parseFloat(r.FormValue("vgap"), 0.5),
		Margin: parseFloat(r.FormValue("margin"), 0.5),
		Pages:  parseInt(r.FormValue("pages"), 2),
		Header: r.FormValue("header"),
		Footer: r.FormValue("footer"),
		Style:  GridStyle(r.FormValue("style")),
	}

	var buf bytes.Buffer
	_, err := GeneratePDF(cfg, &buf)
	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("printme-%.1fcm.pdf", cfg.BoxW)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func parseFloat(s string, def float64) float64 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
