//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/driver/desktop"
	"github.com/belluikit/bellui/driver/mobile"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// pressable is what the Gestures page's tile is to a test: a widget of the
// gallery's own, so its type is not exported, that takes taps, mouse presses
// and touches through the PressTarget it embeds.
type pressable interface {
	core.Widget
	core.Accessible
	core.Disableable
	event.Tappable
	event.Hoverable
	desktop.Mouseable
	mobile.Touchable
}

// The Gestures page's Play tile, built on bellui's PressTarget, shows the
// MD3 pressed state: a press shows the ripple, leaving the tile cancels it, a
// tap runs the tile's action, and the page's box disables it, which ends a
// press under way and then shows no press and runs nothing (gallery#49).
func TestGesturesPressTile(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.GesturesPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(phoneWidth, page.MinSize().Height))

	var tile pressable
	var disable *widget.Check
	var status *widget.Label
	walk(page, func(o core.CanvasObject) {
		// Check and Label first: a Check is pressable too.
		switch v := o.(type) {
		case *widget.Check:
			disable = v
		case *widget.Label:
			if v.Text == "Not tapped yet" {
				status = v
			}
		case pressable:
			if v.AccessibleRole() == core.RoleButton && v.AccessibleLabel() == "Play" {
				tile = v
			}
		}
	})
	if tile == nil || disable == nil || status == nil {
		t.Fatalf("tile %v, box %v, status %v", tile != nil, disable != nil, status != nil)
	}
	var ripple *canvas.Ripple
	walk(tile, func(o core.CanvasObject) {
		if r, ok := o.(*canvas.Ripple); ok {
			ripple = r
		}
	})
	if ripple == nil {
		t.Fatal("the tile draws no ripple")
	}
	if ripple.Visible() {
		t.Error("the ripple shows before any press")
	}

	// A mouse press shows the ripple from where it went down; the pointer
	// leaving cancels it.
	down := &desktop.MouseEvent{PointEvent: event.PointEvent{Position: geo.NewPos(20, 10)}, Button: desktop.MouseButtonPrimary}
	tile.MouseDown(down)
	if !ripple.Visible() || ripple.Center != geo.NewPos(20, 10) {
		t.Errorf("after a press at (20, 10): ripple shown %v, at %v", ripple.Visible(), ripple.Center)
	}
	tile.HoverOut()
	if ripple.Visible() {
		t.Error("leaving the tile mid-press left the ripple showing")
	}

	// A finger does the same.
	tile.TouchDown(&mobile.TouchEvent{PointEvent: event.PointEvent{Position: geo.NewPos(5, 5)}})
	if !ripple.Visible() {
		t.Error("a touch shows no ripple")
	}
	tile.TouchCancel(nil)
	if ripple.Visible() {
		t.Error("a cancelled touch left the ripple showing")
	}

	// A tap runs the action and, with no press before it, still shows one.
	// The test app runs the release to its end at once, so the ripple's
	// radius is what shows it spread.
	ripple.Radius = 0
	test.Tap(tile)
	if status.Text != "Tapped once" {
		t.Errorf("after a tap the status says %q", status.Text)
	}
	if ripple.Radius == 0 {
		t.Error("a tap showed no ripple")
	}
	test.Tap(tile)
	if status.Text != "Tapped 2 times" {
		t.Errorf("after two taps the status says %q", status.Text)
	}

	// Disabled mid-press, the press ends with nothing left showing; then
	// it takes no press and runs nothing.
	tile.MouseDown(down)
	disable.SetChecked(true)
	if !tile.Disabled() {
		t.Fatal("the page's box doesn't disable the tile")
	}
	if ripple.Visible() {
		t.Error("disabling the tile mid-press left the ripple showing")
	}
	ripple.Radius = 0
	tile.MouseDown(down)
	if ripple.Visible() {
		t.Error("a disabled tile showed a press")
	}
	test.Tap(tile)
	if status.Text != "Tapped 2 times" || ripple.Radius != 0 {
		t.Errorf("a disabled tile ran its action (%q) or showed a tap (radius %v)", status.Text, ripple.Radius)
	}
	disable.SetChecked(false)
	if tile.Disabled() {
		t.Error("the tile is still disabled after unticking")
	}
}
