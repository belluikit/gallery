//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Pickers page's wheels report the alarm they spell, and its box disables
// every picker on it (gallery#38).
func TestPickersPage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.PickersPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(900, 3000))

	var wheels []*widget.WheelPicker
	var disable *widget.Check
	var alarm *widget.Label
	var controls []core.Disableable
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *widget.WheelPicker:
			wheels = append(wheels, v)
		case *widget.Check:
			if v.Text == "Disable everything on this page" {
				disable = v
			}
		case *widget.Label:
			if len(v.Text) > 6 && v.Text[:6] == "Alarm:" {
				alarm = v
			}
		}
		switch o.(type) {
		case *widget.WheelPicker, *widget.DatePicker, *widget.TimePicker, *widget.DateTimePicker, *widget.DateEntry, *widget.DateTimeEntry:
			controls = append(controls, o.(core.Disableable))
		}
	})
	if len(wheels) != 4 || disable == nil || alarm == nil {
		t.Fatalf("found %d wheels, disable box %v, alarm line %v", len(wheels), disable != nil, alarm != nil)
	}
	if alarm.Text != "Alarm: 7:30 AM" {
		t.Errorf("first alarm %q, want 7:30 AM", alarm.Text)
	}
	wheels[0].SetSelected(10) // 11
	wheels[1].SetSelected(5)  // 05
	wheels[2].SetSelected(1)  // PM
	if alarm.Text != "Alarm: 11:05 PM" {
		t.Errorf("alarm %q after spinning, want 11:05 PM", alarm.Text)
	}

	// 4 wheels, the date picker, 3 time pickers, the date-time picker's own
	// (its inner date and time pickers count too), and the two entries.
	if len(controls) < 11 {
		t.Fatalf("found %d pickers, want at least 11", len(controls))
	}
	disable.SetChecked(true)
	for _, c := range controls {
		if !c.Disabled() {
			t.Errorf("%T not disabled by the page's box", c)
		}
	}
	disable.SetChecked(false)
	for _, c := range controls {
		if c.Disabled() {
			t.Errorf("%T still disabled after unticking", c)
		}
	}
}
