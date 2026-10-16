package pages

import (
	"strings"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"
)

// The large-data tab uses 5,000 rows on mobile devices and 10,000 elsewhere,
// and its info line reports the real count. Run with -tags mobile for the
// mobile case.
func TestLargeDataTabRowCount(t *testing.T) {
	test.NewTempApp(t)
	want := "10000 rows generated"
	if bellui.CurrentDevice().IsMobile() {
		want = "5000 rows generated"
	}

	var found bool
	var walk func(core.CanvasObject)
	walk = func(o core.CanvasObject) {
		if l, ok := o.(*widget.Label); ok && strings.Contains(l.Text, "rows generated") {
			found = true
			if !strings.Contains(l.Text, want) {
				t.Errorf("info line %q, want it to contain %q", l.Text, want)
			}
		}
		if c, ok := o.(*bellui.Container); ok {
			for _, child := range c.Objects {
				walk(child)
			}
		}
	}
	walk(makeLargeDataTab())
	if !found {
		t.Fatal("no info line found in the large-data tab")
	}
}
