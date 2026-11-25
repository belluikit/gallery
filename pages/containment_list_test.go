//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// firstList returns the first widget.List under o, looking through
// containers and widget renderers.
func firstList(o core.CanvasObject) *widget.List {
	if l, ok := o.(*widget.List); ok {
		return l
	}
	var children []core.CanvasObject
	switch v := o.(type) {
	case *bellui.Container:
		children = v.Objects
	case core.Widget:
		children = test.WidgetRenderer(v).Objects()
	}
	for _, c := range children {
		if l := firstList(c); l != nil {
			return l
		}
	}
	return nil
}

// The Containment page's Lists card shows all five rows: scrolled to the
// bottom, the list has nowhere to go. It used to show one row, because the
// page's VBox gives the list its minimum height (gallery#61). Runs on a
// desktop, and on a phone with -tags mobile.
func TestContainmentListShowsEveryRow(t *testing.T) {
	a := test.NewTempApp(t)
	a.Settings().SetTheme(theme.DefaultTheme())
	w := test.NewTempWindow(t, nil)
	page := pages.ContainmentPage(w)
	w.SetContent(container.NewScroll(container.NewPadded(page)))
	w.Resize(geo.NewSize(900, 800))
	if bellui.CurrentDevice().IsMobile() {
		w.Resize(geo.NewSize(360, 800))
	}

	list := firstList(page)
	if list == nil {
		t.Fatal("Containment page has no widget.List")
	}
	list.ScrollToBottom()
	if off := list.GetScrollOffset(); off != 0 {
		t.Errorf("the list scrolls %.0f dp to reach its last row in a %.0f dp box; want every row showing", off, list.Size().Height)
	}
}
