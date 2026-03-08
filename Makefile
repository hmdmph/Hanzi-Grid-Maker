APP      := hanze-grid-maker
BUILD    := ./build
HEADER   := TS Printables
FOOTER   := I can do this
PAGES    := 2
export GOROOT ?= $(shell go env GOROOT 2>/dev/null || echo /usr/local/Cellar/go/1.26.1/libexec)
ARGS     := -header "$(HEADER)" -footer "$(FOOTER)" -pages "$(PAGES)"

CLI_PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64

.DEFAULT_GOAL := help
.PHONY: help build build-gui build-all package-macos run print ui web tidy clean

## Show help
help:
	@echo ""
	@echo "  ┌──────────────────────────────────────────────┐"
	@echo "  │           🟥  print-square  🟥                │"
	@echo "  │  Generate grid-box PDFs for A4 paper         │"
	@echo "  │  Chinese character practice (田字格 / 米字格)  │"
	@echo "  └──────────────────────────────────────────────┘"
	@echo ""
	@echo "  BUILD TARGETS:"
	@echo "    make build          CLI binary (no CGO needed)"
	@echo "    make build-gui      Native GUI binary (requires CGO)"
	@echo "    make build-all      Cross-compile CLI for macOS, Linux, Windows"
	@echo "    make package-macos  Create macOS .app bundle"
	@echo ""
	@echo "  RUN TARGETS:"
	@echo "    make run            Run CLI with defaults"
	@echo "    make print          Run CLI with HEADER/FOOTER/PAGES vars"
	@echo "    make ui             Launch native desktop GUI"
	@echo "    make web            Launch browser-based web UI"
	@echo "    make tidy           Download Go dependencies"
	@echo "    make clean          Remove build artifacts"
	@echo ""
	@echo "  EXAMPLES:"
	@echo "    # Default grid with your header & footer"
	@echo "    make print"
	@echo ""
	@echo "    # 5 pages of 1.5cm Tian Zi Ge boxes"
	@echo '    make print ARGS='"'"'-bw 1.5 -bh 1.5 -pages 5 -style tianzige -header "Practice" -footer "I can do this"'"'"''
	@echo ""
	@echo "    # Mi Zi Ge style, custom output path"
	@echo '    make print ARGS='"'"'-bw 2 -style mizige -vgap 0.3 -o ./my-grid.pdf'"'"''
	@echo ""
	@echo "    # Launch native GUI app"
	@echo "    make ui"
	@echo ""
	@echo "    # Launch browser-based UI on port 9090"
	@echo '    make web ARGS='"'"'-port 9090'"'"''
	@echo ""
	@echo "  CLI FLAGS:"
	@echo "    -bw      Box width in cm         (default: 1.2)"
	@echo "    -bh      Box height in cm        (default: 1.2)"
	@echo "    -hgap    Horizontal gap in cm    (default: 0)"
	@echo "    -vgap    Vertical gap in cm      (default: 0.5)"
	@echo "    -pages   Number of pages         (default: 2)"
	@echo "    -style   none | tianzige | mizige (default: none)"
	@echo "    -header  Header text             (default: none)"
	@echo "    -footer  Footer text             (default: none)"
	@echo "    -margin  Left/right margin in cm (default: 0.5)"
	@echo "    -o       Output PDF path         (default: ~/Documents/berries/print-format/printme-<size>.pdf)"
	@echo "    -ui      Launch native GUI       (flag, no value)"
	@echo "    -web     Launch browser UI       (flag, no value)"
	@echo "    -port    Web UI port             (default: 8080)"
	@echo ""

## ── Build ────────────────────────────────────────────────

## CLI binary (no CGO, no GUI)
build: tidy
	@mkdir -p $(BUILD)
	CGO_ENABLED=0 go build -o $(BUILD)/$(APP) .
	@echo "Built CLI → $(BUILD)/$(APP)"

## Native GUI binary (requires CGO for Fyne)
build-gui: tidy
	@mkdir -p $(BUILD)
	CGO_ENABLED=1 go build -tags gui -o $(BUILD)/$(APP)-gui .
	@echo "Built GUI → $(BUILD)/$(APP)-gui"

## Cross-compile CLI for all platforms
build-all: tidy
	@mkdir -p $(BUILD)
	@for platform in $(CLI_PLATFORMS); do \
		os=$$(echo $$platform | cut -d/ -f1); \
		arch=$$(echo $$platform | cut -d/ -f2); \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		out="$(BUILD)/$(APP)-$$os-$$arch$$ext"; \
		echo "Building $$out ..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o $$out . || exit 1; \
	done
	@echo ""
	@echo "All CLI binaries → $(BUILD)/"
	@ls -lh $(BUILD)/$(APP)-*

## macOS .app bundle (native GUI)
package-macos: build-gui
	@rm -rf $(BUILD)/PrintSquare.app
	@mkdir -p $(BUILD)/PrintSquare.app/Contents/MacOS
	@mkdir -p $(BUILD)/PrintSquare.app/Contents/Resources
	@cp $(BUILD)/$(APP)-gui $(BUILD)/PrintSquare.app/Contents/MacOS/print-square
	@cp build/darwin/Info.plist $(BUILD)/PrintSquare.app/Contents/Info.plist
	@echo "Packaged → $(BUILD)/PrintSquare.app"
	@echo "To install: drag PrintSquare.app to /Applications"

## ── Run ──────────────────────────────────────────────────

## Run CLI with default settings
run: build
	$(BUILD)/$(APP)

## Run CLI with configured args
print: build
	$(BUILD)/$(APP) $(ARGS)

## Launch native desktop GUI
ui: build-gui
	$(BUILD)/$(APP)-gui -ui

## Launch browser-based web UI
web: build
	$(BUILD)/$(APP) -web $(ARGS)

## ── Misc ─────────────────────────────────────────────────

## Tidy and download dependencies
tidy:
	go mod tidy

## Remove build artifacts
clean:
	rm -rf $(BUILD)/$(APP) $(BUILD)/$(APP)-gui $(BUILD)/$(APP)-* $(BUILD)/PrintSquare.app
