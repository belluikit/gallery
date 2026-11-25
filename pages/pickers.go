package pages

import (
	"fmt"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// PickersPage shows the pickers drawn in place on the page, not in a dialog:
// WheelPicker, DatePicker, TimePicker (12-hour, 24-hour and following the
// locale), DateTimePicker, and the DateEntry and DateTimeEntry fields
// (gallery#38). The Inputs page shows the dialog forms.
func PickersPage(_ bellui.Window) core.CanvasObject {
	var disableables []core.Disableable
	now := time.Now().Truncate(time.Minute)

	// --- Wheel pickers ---
	labels := func(texts ...string) []core.CanvasObject {
		items := make([]core.CanvasObject, len(texts))
		for i, s := range texts {
			items[i] = widget.NewLabelWithStyle(s, bellui.TextAlignCenter, bellui.TextStyle{})
		}
		return items
	}
	var hours, minutes []string
	for h := 1; h <= 12; h++ {
		hours = append(hours, fmt.Sprintf("%d", h))
	}
	for m := 0; m < 60; m++ {
		minutes = append(minutes, fmt.Sprintf("%02d", m))
	}
	alarm := widget.NewLabel("")
	hour := widget.NewWheelPicker(labels(hours...), nil)
	hour.Wrap = true
	hour.SetSelectedSilently(6) // 7
	minute := widget.NewWheelPicker(labels(minutes...), nil)
	minute.Wrap = true
	minute.SetSelectedSilently(30)
	ampm := widget.NewWheelPicker(labels("AM", "PM"), nil)
	showAlarm := func(int) {
		alarm.SetText(fmt.Sprintf("Alarm: %s:%s %s", hours[hour.Selected], minutes[minute.Selected], []string{"AM", "PM"}[ampm.Selected]))
	}
	hour.OnChanged, minute.OnChanged, ampm.OnChanged = showAlarm, showAlarm, showAlarm
	showAlarm(0)
	hour.SetAccessibilityLabel("Hour")
	minute.SetAccessibilityLabel("Minute")
	ampm.SetAccessibilityLabel("AM or PM")

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	day := widget.NewWheelPicker(labels(days...), nil)
	day.Orientation = widget.Horizontal
	day.VisibleItems = 3
	dayStatus := widget.NewLabel("Day: Monday")
	day.OnChanged = func(i int) { dayStatus.SetText("Day: " + days[i]) }
	day.SetAccessibilityLabel("Day")
	disableables = append(disableables, hour, minute, ampm, day)

	wheelBox := func(w float32, o core.CanvasObject) core.CanvasObject {
		return container.NewGridWrap(geo.NewSize(w, o.MinSize().Height), o)
	}
	wheelCard := widget.NewCard("Wheel Picker", "Spin it; the middle item is the choice", container.NewVBox(
		note("Hour and minute wrap round; AM/PM stops at its ends."),
		container.NewHBox(wheelBox(64, hour), wheelBox(64, minute), wheelBox(64, ampm)), alarm,
		widget.NewSeparator(),
		widget.NewLabel("Sideways, three at a time:"), wheelBox(300, day), dayStatus,
	))
	wheelCard.Style = widget.CardStyleOutlined

	// --- Date picker ---
	dateStatus := widget.NewLabel("Date: " + now.Format("Mon, Jan 2, 2006"))
	date := widget.NewDatePicker(now, func(t time.Time) { dateStatus.SetText("Date: " + t.Format("Mon, Jan 2, 2006")) })
	disableables = append(disableables, date)
	dateSection := uncarded("Date Picker", "A month's calendar, in place.", container.NewHBox(date), dateStatus)

	// --- Time pickers ---
	timeStatus := widget.NewLabel("")
	showTime := func(which string) func(h, m int) {
		return func(h, m int) { timeStatus.SetText(fmt.Sprintf("%s: %02d:%02d", which, h, m)) }
	}
	t12 := widget.NewTimePicker(now.Hour(), now.Minute(), false, showTime("12-hour"))
	t24 := widget.NewTimePicker(now.Hour(), now.Minute(), true, showTime("24-hour"))
	tAuto := widget.NewAutoTimePicker(now.Hour(), now.Minute(), showTime("Locale's clock"))
	timeStatus.SetText(fmt.Sprintf("Time: %02d:%02d", now.Hour(), now.Minute()))
	disableables = append(disableables, t12, t24, tAuto)
	cells := []core.CanvasObject{
		container.NewVBox(widget.NewLabel("12-hour"), t12),
		container.NewVBox(widget.NewLabel("24-hour"), t24),
		container.NewVBox(widget.NewLabel("Locale (Auto)"), tAuto),
	}
	cell := geo.Size{}
	for _, c := range cells {
		cell = cell.Max(c.MinSize())
	}
	// A grid that wraps, so the three dials sit side by side on a wide
	// window and stack on a phone.
	timeCard := widget.NewCard("Time Picker", "A clock dial, in place", container.NewVBox(
		note("12-hour, 24-hour, and the locale's own clock (Auto). Side by side on a wide window, stacked on a phone."),
		container.NewGridWrap(cell, cells...),
		timeStatus,
	))
	timeCard.Style = widget.CardStyleOutlined

	// --- Date and time picker ---
	dtStatus := widget.NewLabel("Date and time: " + now.Format("Jan 2, 2006 15:04"))
	dt := widget.NewAutoDateTimePicker(now, func(t time.Time) { dtStatus.SetText("Date and time: " + t.Format("Jan 2, 2006 15:04")) })
	disableables = append(disableables, dt)
	dtSection := uncarded("Date & Time Picker", "Two tabs, in place.", container.NewHBox(dt), dtStatus)

	// --- Entries ---
	dateEntryStatus := widget.NewLabel("Type a date, or open the calendar")
	dateEntry := widget.NewDateEntry()
	dateEntry.OnChanged = func(t time.Time) { dateEntryStatus.SetText("Date entry: " + t.Format("Mon, Jan 2, 2006")) }
	dtEntryStatus := widget.NewLabel("Type a date and time, or open the picker")
	dtEntry := widget.NewDateTimeEntry()
	dtEntry.OnChanged = func(t time.Time) { dtEntryStatus.SetText("Date-time entry: " + t.Format("Jan 2, 2006 15:04")) }
	disableables = append(disableables, dateEntry, dtEntry)
	entryCard := widget.NewCard("Date Entries", "Type it, or use the button", container.NewVBox(
		note("Text fields that read dates in the locale's format, with a picker button."),
		dateEntry, dateEntryStatus, dtEntry, dtEntryStatus,
	))
	entryCard.Style = widget.CardStyleOutlined

	disable := widget.NewCheck("Disable everything on this page", func(on bool) {
		for _, d := range disableables {
			if on {
				d.Disable()
			} else {
				d.Enable()
			}
		}
	})

	return container.NewVBox(disable, wheelCard, dateSection, timeCard, dtSection, entryCard)
}

// uncarded is a section headed like a card but without one. The date picker
// (328 dp) and the date-and-time picker (344 dp) are too wide for a card's
// 30 dp on a 360 dp phone, where a page has 348 dp (gallery#48).
func uncarded(title, explanation string, content ...core.CanvasObject) core.CanvasObject {
	head := []core.CanvasObject{
		widget.NewLabelWithStyle(title, bellui.TextAlignLeading, theme.TypographyTitleMedium),
		note(explanation),
	}
	return container.NewVBox(append(head, content...)...)
}

// note is an explanation that wraps to the page's width, so a page fits a
// phone: a card's subtitle is a single line.
func note(text string) *widget.Label {
	l := widget.NewLabel(text)
	l.Wrapping = bellui.TextWrapWord
	return l
}
