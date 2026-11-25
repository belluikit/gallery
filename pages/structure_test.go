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

// The Structure page's status bar names what is selected in the table, tree
// and grid, and its box lets several accordion sections be open (gallery#38).
func TestStructurePage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.StructurePage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(phoneWidth, page.MinSize().Height))

	var table *widget.Table
	var tree *widget.Tree
	var grid *widget.GridWrap
	var acc *widget.Accordion
	var multi *widget.Check
	var status *widget.StatusBar
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *widget.Table:
			table = v
		case *widget.Tree:
			tree = v
		case *widget.GridWrap:
			grid = v
		case *widget.Accordion:
			acc = v
		case *widget.Check:
			multi = v
		case *widget.StatusBar:
			status = v
		}
	})
	if table == nil || tree == nil || grid == nil || acc == nil || multi == nil || status == nil {
		t.Fatal("a widget of the page is missing")
	}
	selection := func() string {
		// The bar's items: message, separator, selection, spacer, encoding.
		return status.Items[2].(*widget.StatusBarLabel).Text
	}

	for _, c := range []struct {
		do   func()
		want string
	}{
		{func() { table.Select(widget.TableCellID{Row: 3, Col: 3}) }, "Mars: Moons"},
		{func() { tree.Select("Notes.txt") }, "Tree: Notes.txt"},
		{func() { grid.Select(4) }, "Grid: tile 5"},
	} {
		c.do()
		if got := selection(); got != c.want {
			t.Errorf("status bar says %q, want %q", got, c.want)
		}
	}

	open := func() int {
		n := 0
		for _, it := range acc.Items {
			if it.Open {
				n++
			}
		}
		return n
	}
	acc.Open(1)
	if open() != 1 {
		t.Errorf("%d sections open with one allowed", open())
	}
	multi.SetChecked(true)
	acc.Open(2)
	if open() != 2 {
		t.Errorf("%d sections open with several allowed, want 2", open())
	}
}
