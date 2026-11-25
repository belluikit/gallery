//go:build !android

package pages_test

import (
	"math"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// The Gestures page's speed dial opens over its FAB in both directions: its
// actions line up with the FAB's trailing edge, the right one left to right
// and the left one right to left. Since bellui v1.4.0 FABMenu.Show takes the
// FAB's top-left corner in a right-to-left window; given the top-right one,
// the menu opened a FAB's width away from it.
func TestFABMenuOpensAtTheFABsTrailingEdge(t *testing.T) {
	for _, dir := range []bellui.Direction{bellui.DirectionLTR, bellui.DirectionRTL} {
		t.Run(map[bellui.Direction]string{bellui.DirectionLTR: "LTR", bellui.DirectionRTL: "RTL"}[dir], func(t *testing.T) {
			test.NewTempApp(t)
			bellui.SetLayoutDirection(dir)
			defer bellui.SetLayoutDirection(bellui.DirectionLTR)
			w := test.NewTempWindow(t, nil)
			page := pages.GesturesPage(w)
			w.SetContent(page)
			w.Resize(geo.NewSize(412, page.MinSize().Height))

			var fab *widget.FAB
			walk(page, func(o core.CanvasObject) {
				if f, ok := o.(*widget.FAB); ok {
					fab = f
				}
			})
			if fab == nil {
				t.Fatal("the page has no FAB")
			}
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

			abs := bellui.CurrentApp().Driver().AbsolutePositionForObject
			fabEdge, noteEdge := abs(fab).X+fab.Size().Width, abs(note).X+note.Size().Width
			if dir == bellui.DirectionRTL {
				fabEdge, noteEdge = abs(fab).X, abs(note).X
			}
			// The menu's padding insets its actions by a few dp.
			if d := math.Abs(float64(noteEdge - fabEdge)); d > 4 {
				t.Errorf("the Note action's trailing edge is at %.0f, the FAB's at %.0f: %.0f dp apart", noteEdge, fabEdge, d)
			}
		})
	}
}
