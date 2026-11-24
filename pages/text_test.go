//go:build !android

package pages_test

import (
	"strings"
	"testing"

	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Text page's completion field shows suggestions for what is typed and
// takes the one picked, and its box disables every editable widget
// (gallery#38).
func TestTextPage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.TextPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(900, 2400))

	var comp *widget.Entry
	var disable *widget.Check
	var picked *widget.Label
	var controls []core.Disableable
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *widget.Entry:
			if v.PlaceHolder == "Type b or c" {
				comp = v
			}
		case *widget.Check:
			if v.Text == "Disable everything on this page" {
				disable = v
			}
		case *widget.Label:
			if strings.HasPrefix(v.Text, "Picked:") {
				picked = v
			}
		}
		switch o.(type) {
		case *widget.SelectEntry, *widget.SearchSelect, *widget.RichEntry:
			controls = append(controls, o.(core.Disableable))
		}
	})
	if comp == nil || disable == nil || picked == nil {
		t.Fatalf("completion field %v, disable box %v, picked line %v", comp != nil, disable != nil, picked != nil)
	}
	controls = append(controls, comp)

	popup := func() *widget.CompletionPopup {
		var found *widget.CompletionPopup
		for _, o := range w.Canvas().Overlays().List() {
			if !o.Visible() {
				continue
			}
			walk(o, func(c core.CanvasObject) {
				if p, ok := c.(*widget.CompletionPopup); ok && p.Visible() {
					found = p
				}
			})
		}
		return found
	}
	comp.SetText("ca")
	p := popup()
	if p == nil {
		t.Fatal("typing \"ca\" shows no suggestions")
	}
	p.HandleKey(&event.KeyEvent{Name: event.KeyReturn})
	if comp.Text != "canvas" || picked.Text != "Picked: canvas" {
		t.Errorf("Enter picked %q (%q), want canvas", comp.Text, picked.Text)
	}
	if popup() != nil {
		t.Error("the suggestions stay open after a pick")
	}
	comp.SetText("zz")
	if popup() != nil {
		t.Error("suggestions for text nothing starts with")
	}

	if len(controls) != 4 {
		t.Fatalf("found %d editable widgets, want 4", len(controls))
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
