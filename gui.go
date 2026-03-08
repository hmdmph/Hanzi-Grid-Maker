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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ── Custom light theme with red accent ───────────────────

var accentRed = color.NRGBA{R: 220, G: 53, B: 34, A: 255}

type appTheme struct{}

var _ fyne.Theme = (*appTheme)(nil)

func (t *appTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return accentRed
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (t *appTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// ── GUI ──────────────────────────────────────────────────

func launchGUI() {
	a := app.NewWithID(AppID)
	a.Settings().SetTheme(&appTheme{})
	appIcon := fyne.NewStaticResource("icon.png", AppIconBytes)
	a.SetIcon(appIcon)
	w := a.NewWindow(AppName)
	w.SetIcon(appIcon)

	// ── Form fields ──────────────────────────────────────

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

	// ── Save location ────────────────────────────────────

	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "Documents", "berries", "print-format")
	saveEntry := widget.NewEntry()
	saveEntry.SetText(defaultDir)

	browseBtn := widget.NewButtonWithIcon("Browse", theme.FolderOpenIcon(), func() {
		dlg := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			saveEntry.SetText(uri.Path())
		}, w)
		dlg.Show()
	})

	// ── Status ───────────────────────────────────────────

	statusLabel := widget.NewLabel("")

	// ── Generate handler ─────────────────────────────────

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

		outDir := saveEntry.Text
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			dialog.ShowError(err, w)
			return
		}
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

	// ── Layout ───────────────────────────────────────────

	// Title with app logo
	logoImg := canvas.NewImageFromResource(fyne.NewStaticResource("icon.png", AppIconBytes))
	logoImg.SetMinSize(fyne.NewSize(72, 72))
	logoImg.FillMode = canvas.ImageFillContain

	titleText := canvas.NewText("Hanzi Grid Maker - Xizi", accentRed)
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	subtitleText := canvas.NewText(
		"Generate printable grid-box PDFs for Chinese character practice",
		color.NRGBA{R: 130, G: 130, B: 130, A: 255},
	)
	subtitleText.TextSize = 12

	titleBlock := container.NewVBox(
		container.NewHBox(logoImg, titleText),
		subtitleText,
	)

	// BOX DIMENSIONS card - 2-col then 3-col rows like web UI
	dimRow1 := container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Width (cm)"), bwEntry),
		container.NewVBox(widget.NewLabel("Height (cm)"), bhEntry),
	)
	dimRow2 := container.NewGridWithColumns(3,
		container.NewVBox(widget.NewLabel("H-Gap (cm)"), hgapEntry),
		container.NewVBox(widget.NewLabel("V-Gap (cm)"), vgapEntry),
		container.NewVBox(widget.NewLabel("Margin (cm)"), marginEntry),
	)
	dimCard := widget.NewCard("BOX DIMENSIONS", "", container.NewVBox(dimRow1, dimRow2))

	// GUIDE LINE STYLE card
	styleCard := widget.NewCard("GUIDE LINE STYLE", "", container.NewPadded(styleGroup))

	// TEXT & PAGES card
	textContent := container.NewVBox(
		widget.NewLabel("Header Text"),
		headerEntry,
		widget.NewLabel("Footer Text"),
		footerEntry,
		widget.NewLabel("Number of Pages"),
		pagesEntry,
	)
	textCard := widget.NewCard("TEXT & PAGES", "", textContent)

	// SAVE LOCATION card
	saveLine := container.NewBorder(nil, nil, nil, browseBtn, saveEntry)
	saveCard := widget.NewCard("SAVE LOCATION", "", saveLine)

	// Assemble
	content := container.NewVBox(
		titleBlock,
		dimCard,
		styleCard,
		textCard,
		saveCard,
		generateBtn,
		statusLabel,
	)

	scroll := container.NewVScroll(container.NewPadded(content))
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(520, 760))
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
