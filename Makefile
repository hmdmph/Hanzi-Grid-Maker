APP      := print-square
BUILD    := ./build
HEADER   := TS Printables
FOOTER   := I can do this
PAGES    := 2
export GOROOT ?= $(shell go env GOROOT 2>/dev/null || echo /usr/local/Cellar/go/1.26.1/libexec)
ARGS     := -header "$(HEADER)" -footer "$(FOOTER)" -pages "$(PAGES)"

PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64

.DEFAULT_GOAL := help
.PHONY: help build build-all run print ui tidy clean

## Show help
help:
	@echo ""
	@echo "  ┌──────────────────────────────────────────────┐"
	@echo "  │           🟥  print-square  🟥                │"
	@echo "  │  Generate grid-box PDFs for A4 paper         │"
	@echo "  │  Chinese character practice (田字格 / 米字格)  │"
	@echo "  └──────────────────────────────────────────────┘"
	@echo ""
	@echo "  TARGETS:"
	@echo "    make build          Build binary for current OS"
	@echo "    make build-all      Cross-compile for macOS, Linux, Windows"
	@echo "    make run            Run with defaults (1.2cm plain boxes)"
	@echo "    make print          Run with HEADER/FOOTER/PAGES from Makefile vars"
	@echo "    make ui             Launch web UI (opens browser)"
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
	@echo "    # Launch the web UI on port 9090"
	@echo '    make ui ARGS='"'"'-port 9090'"'"''
	@echo ""
	@echo "  FLAGS:"
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
	@echo "    -ui      Launch web UI           (flag, no value)"
	@echo "    -port    Web UI port             (default: 8080)"
	@echo ""

## Build binary for current OS/arch
build: tidy
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/$(APP) .
	@echo "Built → $(BUILD)/$(APP)"

## Cross-compile for all platforms
build-all: tidy
	@mkdir -p $(BUILD)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d/ -f1); \
		arch=$$(echo $$platform | cut -d/ -f2); \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		out="$(BUILD)/$(APP)-$$os-$$arch$$ext"; \
		echo "Building $$out ..."; \
		GOOS=$$os GOARCH=$$arch go build -o $$out . || exit 1; \
	done
	@echo ""
	@echo "All binaries → $(BUILD)/"
	@ls -lh $(BUILD)/

## Run with default settings
run: build
	$(BUILD)/$(APP)

## Run with configured args
print: build
	$(BUILD)/$(APP) $(ARGS)

## Launch web UI
ui: build
	$(BUILD)/$(APP) -ui $(ARGS)

## Tidy and download dependencies
tidy:
	go mod tidy

## Remove build artifacts
clean:
	rm -rf $(BUILD)
