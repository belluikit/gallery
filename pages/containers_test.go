//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Containers page's city buttons push a page onto the navigation and back
// returns, its slider sets the card's opacity, and the second form is laid out
// right-to-left (gallery#38).
func TestContainersPage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.ContainersPage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(phoneWidth, page.MinSize().Height))

	var nav *container.Navigation
	var kyoto *widget.Button
	var faded *container.Opacity
	var slider *widget.Slider
	var scope *container.DirectionScope
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *container.Navigation:
			nav = v
		case *widget.Button:
			if v.Text == "Kyoto" {
				kyoto = v
			}
		case *container.Opacity:
			faded = v
		case *widget.Slider:
			slider = v
		case *container.DirectionScope:
			scope = v
		}
	})
	if nav == nil || kyoto == nil || faded == nil || slider == nil || scope == nil {
		t.Fatal("a container of the page is missing")
	}

	shows := func(text string) bool {
		found := false
		walk(nav, func(o core.CanvasObject) {
			if l, ok := o.(*widget.Label); ok && l.Text == text && l.Visible() {
				found = true
			}
		})
		return found
	}
	test.Tap(kyoto)
	if !shows("All about Kyoto.") {
		t.Error("tapping Kyoto pushes no Kyoto page")
	}
	nav.Back()
	if shows("All about Kyoto.") {
		t.Error("back leaves the Kyoto page showing")
	}

	slider.SetValue(0.25)
	if faded.Opacity != 0.25 {
		t.Errorf("opacity %v after the slider moved to 0.25", faded.Opacity)
	}

	if scope.ScopedDirection() != bellui.DirectionRTL {
		t.Errorf("the second form's direction is %v, want right-to-left", scope.ScopedDirection())
	}
}
