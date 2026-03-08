APP      := print-square
BUILD    := ./build
SRC      := main.go
HEADER   := TS Printables
FOOTER   := I can do this
PAGES    := 2
export GOROOT ?= $(shell go env GOROOT 2>/dev/null || echo /usr/local/Cellar/go/1.26.1/libexec)
ARGS     := -header "$(HEADER)" -footer "$(FOOTER)" -pages "$(PAGES)"

.DEFAULT_GOAL := help
.PHONY: help build run print tidy clean

## Show help
help:
	@echo ""
	@echo "  ┌──────────────────────────────────────────┐"
	@echo "  │         🟥  print-square  🟥              │"
	@echo "  │   Generate grid-box PDFs for A4 paper    │"
	@echo "  └──────────────────────────────────────────┘"
	@echo ""
	@echo "  TARGETS:"
	@echo "    make build          Build the binary"
	@echo "    make run            Run with defaults (1.2cm boxes, no header/footer)"
	@echo "    make print          Run with HEADER & FOOTER from Makefile vars"
	@echo "    make tidy           Download Go dependencies"
	@echo "    make clean          Remove build artifacts"
	@echo ""
	@echo "  EXAMPLES:"
	@echo "    # Default grid with your name & motto"
	@echo "    make print"
	@echo ""
	@echo "    # 5 pages of 1.5cm boxes"
	@echo '    make print ARGS='"'"'-bw 1.5 -bh 1.5 -pages 5 -header "Tashini Sehansa" -footer "I can do this"'"'"''
	@echo ""
	@echo "    # Custom output path, 2cm boxes, 0.3cm row gap"
	@echo '    make print ARGS='"'"'-bw 2 -vgap 0.3 -o ./my-grid.pdf'"'"''
	@echo ""
	@echo "    # Bigger horizontal gap between boxes"
	@echo '    make print ARGS='"'"'-hgap 0.2 -vgap 0.5 -pages 3 -header "Practice" -footer "Keep going!"'"'"''
	@echo ""
	@echo "  FLAGS:"
	@echo "    -bw      Box width in cm         (default: 1.2)"
	@echo "    -bh      Box height in cm        (default: 1.2)"
	@echo "    -hgap    Horizontal gap in cm    (default: 0)"
	@echo "    -vgap    Vertical gap in cm      (default: 0.5)"
	@echo "    -pages   Number of pages         (default: 1)"
	@echo "    -header  Header text             (default: none)"
	@echo "    -footer  Footer text             (default: none)"
	@echo "    -margin  Left/right margin in cm (default: 0.5)"
	@echo "    -o       Output PDF path         (default: ~/Desktop/printme-<size>.pdf)"
	@echo ""

## Build the binary
build: tidy
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/$(APP) $(SRC)
	@echo "Built → $(BUILD)/$(APP)"

## Run with default settings (no header/footer)
run: build
	$(BUILD)/$(APP)

## Run with configured args
print: build
	$(BUILD)/$(APP) $(ARGS)

## Tidy and download dependencies
tidy:
	go mod tidy

## Remove build artifacts
clean:
	rm -rf $(BUILD)
