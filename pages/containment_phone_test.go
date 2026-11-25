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

// findObject returns the first object under o that match accepts, looking
// through containers and widget renderers.
func findObject(o core.CanvasObject, match func(core.CanvasObject) bool) core.CanvasObject {
	if match(o) {
		return o
	}
	var children []core.CanvasObject
	switch v := o.(type) {
	case *bellui.Container:
		children = v.Objects
	case core.Widget:
		children = test.WidgetRenderer(v).Objects()
	}
	for _, c := range children {
		if found := findObject(c, match); found != nil {
			return found
		}
	}
	return nil
}

func cardTitled(title string) func(core.CanvasObject) bool {
	return func(o core.CanvasObject) bool {
		c, ok := o.(*widget.Card)
		return ok && c.Title == title
	}
}

func isTreeView(o core.CanvasObject) bool {
	_, ok := o.(*widget.TreeView)
	return ok
}

// On a desktop the Containment page keeps the TreeView in a rail beside the
// cards, and the three card styles side by side.
func TestContainmentDesktopLayout(t *testing.T) {
	a := test.NewTempApp(t)
	if bellui.CurrentDevice().IsMobile() {
		t.Skip("phone layout; run without -tags mobile")
	}
	a.Settings().SetTheme(theme.DefaultTheme())
	w := test.NewTempWindow(t, nil)
	page := pages.ContainmentPage(w)
	if _, ok := page.(*container.Split); !ok {
		t.Fatalf("desktop Containment page is a %T, want the TreeView rail in a *container.Split", page)
	}
	styles := findObject(page, cardTitled("Cards")).(*widget.Card).Content
	elevated := findObject(styles, cardTitled("Elevated Card"))
	if got, want := styles.MinSize().Width, 3*elevated.MinSize().Width; got < want {
		t.Errorf("card styles need %.0f dp, less than three side by side (%.0f)", got, want)
	}
}

// On a phone the Containment page fits 360 dp: the card styles stack, and
// the TreeView sits above the cards with room for several rows (gallery#51).
// Run with -tags mobile.
func TestContainmentPhoneLayout(t *testing.T) {
	a := test.NewTempApp(t)
	if !bellui.CurrentDevice().IsMobile() {
		t.Skip("phone layout; run with -tags mobile")
	}
	a.Settings().SetTheme(theme.DefaultTheme())
	w := test.NewTempWindow(t, nil)
	page := pages.ContainmentPage(w)
	padded := container.NewPadded(page)
	if got := padded.MinSize().Width; got > phoneWidth {
		t.Fatalf("Containment page needs %.0f dp with its padding, more than a phone's %d", got, phoneWidth)
	}
	w.SetContent(container.NewScroll(padded))
	w.Resize(geo.NewSize(phoneWidth, 800))

	tree := findObject(page, isTreeView)
	cards := findObject(page, cardTitled("Cards"))
	if tree == nil || cards == nil {
		t.Fatal("Containment page has no TreeView or no Cards card")
	}
	treePos := a.Driver().AbsolutePositionForObject(tree)
	cardsPos := a.Driver().AbsolutePositionForObject(cards)
	if treePos.Y+tree.Size().Height > cardsPos.Y {
		t.Errorf("TreeView spans y %.0f–%.0f, overlapping the Cards card at y %.0f; want it above", treePos.Y, treePos.Y+tree.Size().Height, cardsPos.Y)
	}
	if h, one := tree.Size().Height, tree.MinSize().Height; h < 3*one {
		t.Errorf("TreeView is %.0f dp tall, its one-row minimum is %.0f; want room for several rows", h, one)
	}
	elevated := findObject(cards, cardTitled("Elevated Card"))
	outlined := findObject(cards, cardTitled("Outlined Card"))
	if e, o := a.Driver().AbsolutePositionForObject(elevated), a.Driver().AbsolutePositionForObject(outlined); o.Y <= e.Y {
		t.Errorf("Outlined Card at y %.0f, Elevated Card at y %.0f; want the styles stacked", o.Y, e.Y)
	}
}
