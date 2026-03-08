# 🟥 print-square

A simple, cross-platform tool to generate printable grid-box PDFs on A4 paper — perfect for **Chinese character writing practice**.

Supports three grid styles used in Chinese calligraphy education:

| Style | Name | Description |
|-------|------|-------------|
| `none` | Plain | Empty boxes |
| `tianzige` | 田字格 (Tián Zì Gé) | Cross guides — divides each box into 4 quadrants |
| `mizige` | 米字格 (Mǐ Zì Gé) | Star guides — cross + diagonals (resembles 米) |

## Quick Start

### Option A: Web UI (recommended for non-technical users)

```bash
make ui
```

This opens a browser with a simple interface where you can adjust all settings and download your PDF.

### Option B: Command Line

```bash
# Plain 1.2cm boxes, 2 pages
make print

# Tian Zi Ge style, 1.5cm boxes, 5 pages
make print ARGS='-bw 1.5 -bh 1.5 -style tianzige -pages 5 -header "Practice" -footer "I can do this"'

# Mi Zi Ge style, custom output
make print ARGS='-bw 2 -style mizige -o ./my-grid.pdf'
```

## Installation

### From source

Requires [Go 1.21+](https://go.dev/dl/).

```bash
git clone https://github.com/nicober/print-square.git
cd print-square
make build
```

The binary is built to `./build/print-square`.

### Cross-compile for all platforms

```bash
make build-all
```

Produces binaries in `./build/` for:
- macOS (Intel + Apple Silicon)
- Linux (amd64 + arm64)
- Windows (amd64)

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-bw` | `1.2` | Box width in cm |
| `-bh` | `1.2` | Box height in cm |
| `-hgap` | `0` | Horizontal gap between boxes (cm) |
| `-vgap` | `0.5` | Vertical gap between rows (cm) |
| `-margin` | `0.5` | Left/right page margin (cm) |
| `-pages` | `2` | Number of pages |
| `-style` | `none` | Guide style: `none`, `tianzige`, `mizige` |
| `-header` | | Header text (centered) |
| `-footer` | | Footer text (left-aligned, page number on right) |
| `-o` | `~/Documents/berries/print-format/printme-<size>.pdf` | Output file path |
| `-ui` | `false` | Launch web UI instead of CLI |
| `-port` | `8080` | Port for web UI server |

## Web UI

Run `print-square -ui` (or `make ui`) to launch a local web server with a clean form interface. Your browser opens automatically. No internet connection required — the UI is fully embedded in the binary.

The UI lets you:
- Set box dimensions and gaps
- Pick a grid style with visual previews
- Add header/footer text
- Set page count
- Download the generated PDF instantly

## Project Structure

```
print-square/
├── main.go          # CLI entry point & flag parsing
├── generate.go      # PDF generation engine
├── server.go        # Web UI HTTP server
├── ui/
│   └── index.html   # Embedded web interface
├── Makefile         # Build, run, cross-compile
├── LICENSE          # MIT
└── README.md
```

## License

MIT — see [LICENSE](LICENSE).
