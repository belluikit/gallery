package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// ActionsPage shows Buttons, FAB, Icon Buttons, Segmented Buttons, Button Groups.
func ActionsPage(_ bellui.Window) core.CanvasObject {
	// --- Buttons ---
	buttonsCard := widget.NewCard("Buttons", "5 MD3 button styles", container.NewVBox(
		container.NewHBox(
			&widget.Button{Text: "Filled", Style: widget.ButtonStyleFilled, OnTapped: func() {}},
			&widget.Button{Text: "Filled Tonal", Style: widget.ButtonStyleFilledTonal, OnTapped: func() {}},
			&widget.Button{Text: "Elevated", Style: widget.ButtonStyleElevated, OnTapped: func() {}},
		),
		container.NewHBox(
			&widget.Button{Text: "Outlined", Style: widget.ButtonStyleOutlined, OnTapped: func() {}},
			&widget.Button{Text: "Text", Style: widget.ButtonStyleText, OnTapped: func() {}},
		),
	))
	buttonsCard.Style = widget.CardStyleOutlined

	// --- FAB ---
	fabCard := widget.NewCard("Floating Action Button", "Standard and Extended", container.NewHBox(
		widget.NewFAB(theme.ContentAddIcon(), func() {}),
		widget.NewExtendedFAB("Create", theme.ContentAddIcon(), func() {}),
	))
	fabCard.Style = widget.CardStyleOutlined

	// --- Icon Buttons ---
	standard := widget.NewIconButton(theme.SettingsIcon(), func() {})
	filled := widget.NewIconButton(theme.SettingsIcon(), func() {})
	filled.Style = widget.IconButtonStyleFilled
	tonal := widget.NewIconButton(theme.SettingsIcon(), func() {})
	tonal.Style = widget.IconButtonStyleFilledTonal
	outlined := widget.NewIconButton(theme.SettingsIcon(), func() {})
	outlined.Style = widget.IconButtonStyleOutlined

	iconBtnCard := widget.NewCard("Icon Buttons", "Standard, Filled, Tonal, Outlined", container.NewHBox(
		container.NewVBox(widget.NewLabel("Standard"), standard),
		container.NewVBox(widget.NewLabel("Filled"), filled),
		container.NewVBox(widget.NewLabel("Tonal"), tonal),
		container.NewVBox(widget.NewLabel("Outlined"), outlined),
	))
	iconBtnCard.Style = widget.CardStyleOutlined

	// --- Segmented Button ---
	seg := widget.NewSegmentedButton([]string{"Day", "Week", "Month", "Year"}, func(int) {})
	segCard := widget.NewCard("Segmented Button", "Single-select toggle group", seg)
	segCard.Style = widget.CardStyleOutlined

	// --- Button Group ---
	bg := widget.NewButtonGroup(
		&widget.Button{Text: "Cut", Style: widget.ButtonStyleOutlined, OnTapped: func() {}},
		&widget.Button{Text: "Copy", Style: widget.ButtonStyleOutlined, OnTapped: func() {}},
		&widget.Button{Text: "Paste", Style: widget.ButtonStyleOutlined, OnTapped: func() {}},
	)
	bgCard := widget.NewCard("Button Group", "Connected action buttons", bg)
	bgCard.Style = widget.CardStyleOutlined

	// --- Split Button ---
	split := widget.NewSplitButton("Save", func() {}, []*bellui.MenuItem{
		bellui.NewMenuItem("Save As...", func() {}),
		bellui.NewMenuItem("Export PDF", func() {}),
	})
	splitCard := widget.NewCard("Split Button", "Primary action + dropdown menu", split)
	splitCard.Style = widget.CardStyleOutlined

	return container.NewVBox(buttonsCard, fabCard, iconBtnCard, segCard, bgCard, splitCard)
}
