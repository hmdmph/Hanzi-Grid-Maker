//go:build gui

package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func launchGUI() {
	a := app.NewWithID(AppID)
	w := a.NewWindow(AppName)

	// ── Form fields ──────────────────────────────────────────

	bwEntry := newNumEntry("1.2")
	bhEntry := newNumEntry("1.2")
	hgapEntry := newNumEntry("0")
	vgapEntry := newNumEntry("0.5")
	marginEntry := newNumEntry("0.5")
	pagesEntry := newNumEntry("2")

	headerEntry := widget.NewEntry()
	headerEntry.SetPlaceHolder("e.g. TS Printables")

	footerEntry := widget.NewEntry()
	footerEntry.SetPlaceHolder("e.g. I can do this")

	styleGroup := widget.NewRadioGroup(
		[]string{"Plain", "Tián Zì Gé (田字格)", "Mǐ Zì Gé (米字格)"},
		nil,
	)
	styleGroup.SetSelected("Plain")
	styleGroup.Horizontal = true

	statusLabel := widget.NewLabel("")

	// ── Generate handler ─────────────────────────────────────

	generateBtn := widget.NewButton("Generate PDF", func() {
		style := StyleNone
		switch styleGroup.Selected {
		case "Tián Zì Gé (田字格)":
			style = StyleTianZiGe
		case "Mǐ Zì Gé (米字格)":
			style = StyleMiZiGe
		}

		cfg := GridConfig{
			BoxW:   safeParseFl(bwEntry.Text, 1.2),
			BoxH:   safeParseFl(bhEntry.Text, 1.2),
			HGap:   safeParseFl(hgapEntry.Text, 0),
			VGap:   safeParseFl(vgapEntry.Text, 0.5),
			Margin: safeParseFl(marginEntry.Text, 0.5),
			Pages:  safeParseIn(pagesEntry.Text, 2),
			Header: headerEntry.Text,
			Footer: footerEntry.Text,
			Style:  style,
		}

		home, _ := os.UserHomeDir()
		outDir := filepath.Join(home, "Documents", "berries", "print-format")
		os.MkdirAll(outDir, 0o755)
		outPath := filepath.Join(outDir, fmt.Sprintf("printme-%.1fcm.pdf", cfg.BoxW))

		f, err := os.Create(outPath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer f.Close()

		info, err := GeneratePDF(cfg, f)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		statusLabel.SetText(fmt.Sprintf("✅  %d × %d grid, %d page(s)", info.Cols, info.Rows, info.Pages))

		dialog.ShowConfirm("PDF Generated",
			fmt.Sprintf("Saved to:\n%s\n\nOpen the file?", outPath),
			func(open bool) {
				if open {
					openBrowser(outPath)
				}
			}, w)
	})
	generateBtn.Importance = widget.HighImportance

	// ── Layout ───────────────────────────────────────────────

	title := widget.NewRichTextFromMarkdown("# 🟥 " + AppName)
	subtitle := canvas.NewText("Grid PDF for Chinese character practice", color.NRGBA{R: 140, G: 140, B: 140, A: 255})
	subtitle.TextSize = 13

	dimForm := widget.NewForm(
		widget.NewFormItem("Width (cm)", bwEntry),
		widget.NewFormItem("Height (cm)", bhEntry),
		widget.NewFormItem("H-Gap (cm)", hgapEntry),
		widget.NewFormItem("V-Gap (cm)", vgapEntry),
		widget.NewFormItem("Margin (cm)", marginEntry),
	)

	textForm := widget.NewForm(
		widget.NewFormItem("Header", headerEntry),
		widget.NewFormItem("Footer", footerEntry),
		widget.NewFormItem("Pages", pagesEntry),
	)

	sectionStyle := widget.NewLabelWithStyle("Guide Line Style", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sectionText := widget.NewLabelWithStyle("Header & Footer", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sectionDim := widget.NewLabelWithStyle("Box Dimensions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		sectionDim,
		dimForm,
		widget.NewSeparator(),
		sectionStyle,
		styleGroup,
		widget.NewSeparator(),
		sectionText,
		textForm,
		layout.NewSpacer(),
		generateBtn,
		statusLabel,
	)

	scroll := container.NewVScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(480, 680))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func newNumEntry(val string) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(val)
	return e
}

func safeParseFl(s string, def float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}

func safeParseIn(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
