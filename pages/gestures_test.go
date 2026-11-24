//go:build !android

package pages_test

import (
	"strings"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Gestures page reports swipes, its FAB opens a speed dial whose actions
// run, its context menu's items run, and its box disables the swipe areas and
// the FAB (gallery#38).
func TestGesturesPage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.GesturesPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(412, page.MinSize().Height))

	var areas []*widget.SwipeArea
	var fab *widget.FAB
	var menu *widget.ContextMenu
	var disable *widget.Check
	label := func(prefix string) string {
		got := ""
		walk(page, func(o core.CanvasObject) {
			if l, ok := o.(*widget.Label); ok && strings.HasPrefix(l.Text, prefix) {
				got = l.Text
			}
		})
		return got
	}
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *widget.SwipeArea:
			areas = append(areas, v)
		case *widget.FAB:
			fab = v
		case *widget.ContextMenu:
			menu = v
		case *widget.Check:
			disable = v
		}
	})
	if len(areas) != 2 || fab == nil || menu == nil || disable == nil {
		t.Fatalf("found %d swipe areas, FAB %v, context menu %v, box %v", len(areas), fab != nil, menu != nil, disable != nil)
	}
	abs := func(o core.CanvasObject) geo.Position {
		return bellui.CurrentApp().Driver().AbsolutePositionForObject(o)
	}
	middle := func(o core.CanvasObject) geo.Position {
		return abs(o).Add(geo.NewPos(o.Size().Width/2, o.Size().Height/2))
	}

	// The plain area reports a sideways swipe.
	test.Drag(w.Canvas(), middle(areas[1]), 150, 0)
	if got := label("Swiped"); got != "Swiped right" {
		t.Errorf("plain area: %q, want Swiped right", got)
	}
	// The card is thrown off to the left and the next one comes.
	test.Drag(w.Canvas(), middle(areas[0]), -200, 0)
	if got := label("Card "); got != "Card 1: thrown left" {
		t.Errorf("card: %q, want Card 1: thrown left", got)
	}

	// The FAB opens its actions; Note runs and closes them.
	test.Tap(fab)
	var note *widget.FAB
	for _, o := range w.Canvas().Overlays().List() {
		walk(o, func(c core.CanvasObject) {
			if e, ok := c.(*widget.FAB); ok && e.Text == "Note" {
				note = e
			}
		})
	}
	if note == nil {
		t.Fatal("the FAB opens no Note action")
	}
	if p := abs(note); p.X < 0 || p.Y < 0 || p.X+note.Size().Width > w.Canvas().Size().Width {
		t.Errorf("the Note action is at %v, %v wide: off the %v canvas", p, note.Size().Width, w.Canvas().Size())
	}
	test.Tap(note)
	if got := label("Created"); got != "Created a note" {
		t.Errorf("FAB menu: %q, want Created a note", got)
	}

	// The context menu's Rename runs.
	test.TapSecondary(menu)
	var rename *widget.PopUpMenu
	for _, o := range w.Canvas().Overlays().List() {
		walk(o, func(c core.CanvasObject) {
			if m, ok := c.(*widget.PopUpMenu); ok {
				rename = m
			}
		})
	}
	if rename == nil {
		t.Fatal("a right-click opens no menu")
	}
	rename.ActivateNext() // Copy
	rename.ActivateNext() // Rename
	rename.TriggerLast()
	if got := label("Chose"); got != "Chose Rename on report.pdf" {
		t.Errorf("context menu: %q, want Chose Rename on report.pdf", got)
	}

	disable.SetChecked(true)
	for _, d := range []core.Disableable{areas[0], areas[1], fab} {
		if !d.Disabled() {
			t.Errorf("%T not disabled by the page's box", d)
		}
	}
	test.Drag(w.Canvas(), middle(areas[1]), 0, 150)
	if got := label("Swiped"); got != "Swiped right" {
		t.Errorf("a disabled area reported a swipe: %q", got)
	}
	disable.SetChecked(false)
	for _, d := range []core.Disableable{areas[0], areas[1], fab} {
		if d.Disabled() {
			t.Errorf("%T still disabled after unticking", d)
		}
	}
}
