package pages

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"

	"github.com/belluikit/bellui/x/charts/chart"
)

// On a phone the drawing-tools chart sits below four rows of tool buttons.
// It must still get a plot-sized height, not just its 150pt minimum, which
// left only the title visible. Run with -tags mobile.
func TestDrawingChartHeightOnPhone(t *testing.T) {
	test.NewTempApp(t)
	if !bellui.CurrentDevice().IsMobile() {
		t.Skip("phone layout; run with -tags mobile")
	}
	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(412, 860)) // a Pixel 9a in dp
	page := FychartsPage(w)
	w.SetContent(container.NewScroll(container.NewPadded(page))) // as main.go hosts it on a phone

	var tabs *container.AppTabs
	var findTabs func(core.CanvasObject)
	findTabs = func(o core.CanvasObject) {
		switch v := o.(type) {
		case *container.AppTabs:
			tabs = v
		case *bellui.Container:
			for _, c := range v.Objects {
				findTabs(c)
			}
		}
	}
	findTabs(page)
	for _, it := range tabs.Items {
		if it.Text != "Trading" {
			continue
		}
		tabs.Select(it)
		cells := it.Content.(*bellui.Container).Objects
		var draw *chart.Chart
		var findChart func(core.CanvasObject)
		findChart = func(o core.CanvasObject) {
			switch v := o.(type) {
			case *chart.Chart:
				draw = v
			case *bellui.Container:
				for _, c := range v.Objects {
					findChart(c)
				}
			}
		}
		findChart(cells[len(cells)-1])
		if draw == nil {
			t.Fatal("no chart in the drawing-tools cell")
		}
		if h := draw.Size().Height; h < 300 {
			t.Fatalf("drawing chart is %v tall on a phone, want at least 300", h)
		}
		return
	}
	t.Fatal("no Trading tab")
}
