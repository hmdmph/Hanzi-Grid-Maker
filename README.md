# 🟥 Hanzi Grid Maker - Xizi

A cross-platform tool to generate printable grid-box PDFs on A4 paper — built for **Chinese character writing practice**.

Available as a **native desktop app** (macOS `.app`, Windows `.exe`) or a **CLI** for power users.

## Grid Styles

| Style | Name | Description |
|-------|------|-------------|
| `none` | Plain | Empty boxes |
| `tianzige` | 田字格 (Tián Zì Gé) | Cross guides — divides each box into 4 quadrants |
| `mizige` | 米字格 (Mǐ Zì Gé) | Star guides — cross + diagonals (resembles 米) |

## Quick Start

### Option A: Native Desktop App (recommended)

Double-click **PrintSquare.app** (macOS) or **print-square-gui.exe** (Windows).

A native window opens with form fields, style selection, and a Generate button. No browser, no terminal needed.

**Build the app yourself:**

```bash
make build-gui      # build native GUI binary
make package-macos  # create macOS .app bundle
```

### Option B: Browser-Based Web UI

```bash
make web
```

Opens a browser with an embedded UI. Works without CGO — useful on headless machines or when you can't build the native GUI.

### Option C: Command Line (pro users)

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
make build       # CLI only (no CGO needed)
make build-gui   # native GUI (requires CGO)
```

### macOS App Bundle

```bash
make package-macos
# → ./build/PrintSquare.app
# Drag to /Applications to install
```

### Cross-compile CLI for all platforms

```bash
make build-all
```

Produces CLI binaries in `./build/` for:
- **macOS** (Intel + Apple Silicon)
- **Linux** (amd64 + arm64)
- **Windows** (amd64)

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
| `-ui` | `false` | Launch native desktop GUI |
| `-web` | `false` | Launch browser-based web UI |
| `-port` | `8080` | Port for web UI server |

## Architecture

The project uses **Go build tags** to keep the CLI small and dependency-free, while the native GUI is opt-in:

- **`make build`** → CLI binary. No CGO, no native deps. Cross-compiles to any OS.
- **`make build-gui`** → Native GUI binary. Uses [Fyne](https://fyne.io) toolkit (requires CGO). Produces a real desktop window — not a browser.
- The macOS `.app` bundle **auto-detects** that it's running from a bundle and launches the GUI without needing `-ui`.

### Project Structure

```
print-square/
├── main.go          # Entry point: CLI / GUI / web routing
├── generate.go      # PDF generation engine + metadata
├── gui.go           # Native Fyne GUI (build tag: gui)
├── gui_stub.go      # Stub when built without GUI
├── server.go        # Browser-based web UI server
├── ui/
│   └── index.html   # Embedded HTML for web UI
├── build/
│   └── darwin/
│       └── Info.plist  # macOS app bundle config
├── Makefile         # Build, package, run targets
├── LICENSE          # MIT
├── README.md
└── blog.md          # Origin story
```

## PDF Metadata

Generated PDFs include embedded metadata:
- **Title:** Hanzi Grid Maker — Xizi
- **Author:** TS Printables
- **Subject:** Grid practice paper for Chinese character writing
- **Creator:** Hanzi Grid Maker — Xizi v1.0.0
- **Keywords:** grid, chinese, practice, tianzige, mizige, calligraphy, writing

## License

MIT — see [LICENSE](LICENSE).
