//go:build !android

package pages_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Communication page's More Dialogs card opens the colour picker and the
// three custom dialogs, and reports what each one returned (gallery#47).
func TestMoreDialogs(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.CommunicationPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(900, page.MinSize().Height))

	button := func(root core.CanvasObject, text string) *widget.Button {
		var b *widget.Button
		walk(root, func(o core.CanvasObject) {
			if v, ok := o.(*widget.Button); ok && v.Text == text {
				b = v
			}
		})
		if b == nil {
			t.Fatalf("no %q button", text)
		}
		return b
	}
	var status *widget.Label
	walk(page, func(o core.CanvasObject) {
		if l, ok := o.(*widget.Label); ok && l.Text == "Nothing chosen yet" {
			status = l
		}
	})
	if status == nil {
		t.Fatal("the More Dialogs card has no status line")
	}
	top := func() core.CanvasObject { return w.Canvas().Overlays().Top() }
	open := func() bool { return len(w.Canvas().Overlays().List()) > 0 }

	// The colour picker: a swatch tapped in it is the colour chosen.
	test.Tap(button(page, "Colour Picker"))
	if !open() {
		t.Fatal("Colour Picker opens no dialog")
	}
	// The swatches are bellui's unexported colorButton; there is no other
	// way to find one from outside.
	var swatch event.Tappable
	walk(top(), func(o core.CanvasObject) {
		if v, ok := o.(event.Tappable); ok && swatch == nil && fmt.Sprintf("%T", o) == "*dialog.colorButton" {
			swatch = v
		}
	})
	if swatch == nil {
		t.Fatal("the colour picker shows no swatch to tap")
	}
	swatch.Tapped(&event.PointEvent{})
	if !strings.HasPrefix(status.Text, "Colour: #") {
		t.Errorf("after a swatch was tapped the status says %q", status.Text)
	}
	if open() {
		t.Error("the colour picker stays open after a colour was chosen")
	}

	// The advanced picker opens with the Advanced section.
	test.Tap(button(page, "Advanced Colour Picker"))
	var acc *widget.Accordion
	walk(top(), func(o core.CanvasObject) {
		if v, ok := o.(*widget.Accordion); ok {
			acc = v
		}
	})
	if acc == nil {
		t.Error("the advanced colour picker has no Advanced section")
	}
	test.Tap(button(top(), "Cancel"))

	// Custom: the content, and one button that closes it.
	test.Tap(button(page, "Custom"))
	shown := false
	walk(top(), func(o core.CanvasObject) {
		if l, ok := o.(*widget.Label); ok && l.Text == "Lake at dawn.jpg" {
			shown = true
		}
	})
	if !shown {
		t.Error("the custom dialog doesn't show its content")
	}
	test.Tap(button(top(), "Close"))
	if open() {
		t.Error("Close leaves the custom dialog open")
	}

	// Custom confirm: what was typed comes back with Rename.
	test.Tap(button(page, "Custom Confirm"))
	var name *widget.Entry
	walk(top(), func(o core.CanvasObject) {
		if e, ok := o.(*widget.Entry); ok {
			name = e
		}
	})
	if name == nil {
		t.Fatal("the custom confirm dialog has no entry")
	}
	name.SetText("Misty lake")
	test.Tap(button(top(), "Rename"))
	if status.Text != "Renamed to Misty lake" {
		t.Errorf("after Rename the status says %q", status.Text)
	}

	// Custom without buttons: a choice in the content closes it.
	test.Tap(button(page, "Custom, No Buttons"))
	test.Tap(button(top(), "Link"))
	if status.Text != "Shared by link" {
		t.Errorf("after Link the status says %q", status.Text)
	}
	if open() {
		t.Error("choosing Link leaves the dialog open")
	}
}
