package pages

import (
	"fmt"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"
)

// InputsPage shows Text Fields, Checkbox, Radio, Slider, Switch, Chips, NumericEntry, SearchField.
func InputsPage(_ bellui.Window) core.CanvasObject {
	// --- Text Fields ---
	filledEntry := widget.NewEntry()
	filledEntry.PlaceHolder = "Filled text field"
	outlinedEntry := &widget.Entry{PlaceHolder: "Outlined text field", Style: widget.EntryStyleOutlined}
	outlinedEntry.ExtendBaseWidget(outlinedEntry)

	textCard := widget.NewCard("Text Fields", "Filled and Outlined", container.NewVBox(
		widget.NewLabel("Filled"),
		filledEntry,
		widget.NewLabel("Outlined"),
		outlinedEntry,
	))
	textCard.Style = widget.CardStyleOutlined

	// --- Checkbox & Radio ---
	check1 := widget.NewCheck("Option A", func(bool) {})
	check2 := widget.NewCheck("Option B", func(bool) {})
	check2.SetChecked(true)
	radio := widget.NewRadioGroup([]string{"Choice 1", "Choice 2", "Choice 3"}, func(string) {})
	radio.SetSelected("Choice 1")

	selCard := widget.NewCard("Selection Controls", "Checkbox and Radio", container.NewHBox(
		container.NewVBox(widget.NewLabel("Checkbox"), check1, check2),
		container.NewVBox(widget.NewLabel("Radio"), radio),
	))
	selCard.Style = widget.CardStyleOutlined

	// --- Switch ---
	sw := widget.NewSwitch(func(bool) {})
	swOn := widget.NewSwitch(func(bool) {})
	swOn.On = true
	swDisabled := widget.NewSwitch(nil)
	swDisabled.Disable()

	switchCard := widget.NewCard("Switch", "Toggle on/off", container.NewHBox(
		container.NewVBox(widget.NewLabel("Off"), sw),
		container.NewVBox(widget.NewLabel("On"), swOn),
		container.NewVBox(widget.NewLabel("Disabled"), swDisabled),
	))
	switchCard.Style = widget.CardStyleOutlined

	// --- Slider ---
	sliderLabel := widget.NewLabel("50")
	slider := widget.NewSlider(0, 100)
	slider.Value = 50
	slider.Step = 1
	slider.ShowValueIndicator = true
	slider.OnChanged = func(v float64) { sliderLabel.SetText(fmt.Sprintf("%.0f", v)) }

	tickSlider := widget.NewSlider(0, 5)
	tickSlider.Step = 1
	tickSlider.ShowTickMarks = true

	rangeSlider := widget.NewRangeSlider(0, 100)
	rangeSlider.LowValue = 20
	rangeSlider.HighValue = 80

	sliderCard := widget.NewCard("Sliders", "Continuous, discrete, range", container.NewVBox(
		container.NewHBox(widget.NewLabel("Continuous:"), sliderLabel),
		slider,
		widget.NewLabel("Discrete (with ticks):"),
		tickSlider,
		widget.NewLabel("Range:"),
		rangeSlider,
	))
	sliderCard.Style = widget.CardStyleOutlined

	// --- Chips ---
	filter1 := widget.NewFilterChip("Running", func(bool) {})
	filter2 := widget.NewFilterChip("Cycling", func(bool) {})
	filter2.Selected = true
	assist := widget.NewChip("Get directions", func() {})

	chipCard := widget.NewCard("Chips", "Filter and Assist", container.NewHBox(filter1, filter2, assist))
	chipCard.Style = widget.CardStyleOutlined

	// --- NumericEntry ---
	numEntry := widget.NewNumericEntry()
	numEntry.PlaceHolder = "0.00"
	intEntry := widget.NewIntegerEntry()
	intEntry.PlaceHolder = "0"

	numCard := widget.NewCard("Numeric Entry", "Right-aligned, filtered input", container.NewVBox(
		widget.NewLabel("Float:"),
		numEntry,
		widget.NewLabel("Integer:"),
		intEntry,
	))
	numCard.Style = widget.CardStyleOutlined

	// --- Search ---
	search := widget.NewSearchField(func(string) {})
	search.PlaceHolder = "Search..."
	searchCard := widget.NewCard("Search Field", "Search with clear button", search)
	searchCard.Style = widget.CardStyleOutlined

	// --- MD3 Date Picker ---
	datePickerLabel := widget.NewLabel("No date selected")
	datePickerBtn := &widget.Button{Text: "Pick Date", Style: widget.ButtonStyleOutlined, OnTapped: func() {}}
	datePickerBtn.OnTapped = func() {
		widget.ShowDatePicker(time.Now(), func(d time.Time) {
			datePickerLabel.SetText(d.Format("Jan 2, 2006"))
		}, bellui.CurrentApp().Driver().AllWindows()[0].Canvas())
	}
	dateCard := widget.NewCard("Date Picker (MD3)", "Calendar grid with month navigation", container.NewVBox(datePickerLabel, datePickerBtn))
	dateCard.Style = widget.CardStyleOutlined

	// --- MD3 Time Picker ---
	timePickerLabel := widget.NewLabel("No time selected")
	timePickerBtn := &widget.Button{Text: "Pick Time", Style: widget.ButtonStyleOutlined, OnTapped: func() {}}
	timePickerBtn.OnTapped = func() {
		now := time.Now()
		widget.ShowTimePicker(now.Hour(), now.Minute(), false, func(h, m int) {
			timePickerLabel.SetText(fmt.Sprintf("%02d:%02d", h, m))
		}, bellui.CurrentApp().Driver().AllWindows()[0].Canvas())
	}
	timeCard := widget.NewCard("Time Picker (MD3)", "Clock dial with hour/minute selection", container.NewVBox(timePickerLabel, timePickerBtn))
	timeCard.Style = widget.CardStyleOutlined

	// --- MD3 DateTime Picker ---
	dtPickerLabel := widget.NewLabel("No date/time selected")
	dtPickerBtn := &widget.Button{Text: "Pick Date & Time", Style: widget.ButtonStyleOutlined, OnTapped: func() {}}
	dtPickerBtn.OnTapped = func() {
		widget.ShowDateTimePicker(time.Now(), false, func(dt time.Time) {
			dtPickerLabel.SetText(dt.Format("Jan 2, 2006 03:04 PM"))
		}, bellui.CurrentApp().Driver().AllWindows()[0].Canvas())
	}
	dtCard := widget.NewCard("DateTime Picker (MD3)", "Combined date and time selection", container.NewVBox(dtPickerLabel, dtPickerBtn))
	dtCard.Style = widget.CardStyleOutlined

	// --- Select / Dropdown ---
	sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, func(string) {})
	sel.PlaceHolder = "Choose an option"
	selectCard := widget.NewCard("Dropdown Menu", "Select from options", sel)
	selectCard.Style = widget.CardStyleOutlined

	return container.NewVBox(textCard, selCard, switchCard, sliderCard, chipCard, numCard, searchCard, dateCard, timeCard, dtCard, selectCard)
}
