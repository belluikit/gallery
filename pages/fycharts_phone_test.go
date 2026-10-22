package pages

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/theme"

	"github.com/belluikit/bellui/x/charts/chart"
)

// phoneCharts shows the Charts page as main.go hosts it on a Pixel 9a, and
// returns its tabs.
func phoneCharts(t *testing.T) *container.AppTabs {
	t.Helper()
	test.NewTempApp(t)
	if !bellui.CurrentDevice().IsMobile() {
		t.Skip("phone layout; run with -tags mobile")
	}
	test.ApplyTheme(t, theme.DefaultTheme())
	w := test.NewWindow(nil)
	t.Cleanup(w.Close)
	w.Resize(geo.NewSize(411, 923)) // a Pixel 9a in dp
	page := FychartsPage(w)
	w.SetContent(container.NewScroll(container.NewPadded(page)))
	var tabs *container.AppTabs
	walkCharts(page, func(o core.CanvasObject) {
		if v, ok := o.(*container.AppTabs); ok {
			tabs = v
		}
	})
	if tabs == nil {
		t.Fatal("no tabs on the Charts page")
	}
	return tabs
}

// walkCharts visits o and the containers and scrolls under it.
func walkCharts(o core.CanvasObject, visit func(core.CanvasObject)) {
	visit(o)
	switch v := o.(type) {
	case *bellui.Container:
		for _, c := range v.Objects {
			walkCharts(c, visit)
		}
	case *container.Scroll:
		walkCharts(v.Content, visit)
	}
}

// On a phone every chart in every tab is phoneChartHeight tall: the
// drawing-tools row, which needs more, must not stretch the others, here
// or in another tab. The drawing chart itself keeps a plot of at least
// 300pt below its tool buttons (it showed only its title at 150pt).
// Run with -tags mobile.
func TestChartHeightsOnPhone(t *testing.T) {
	tabs := phoneCharts(t)
	for _, it := range tabs.Items {
		tabs.Select(it)
		var n int
		walkCharts(it.Content, func(o core.CanvasObject) {
			var h float32
			switch v := o.(type) {
			case *chart.Chart:
				h = v.Size().Height
			case *chart.ChartGroup:
				h = v.Size().Height
			default:
				return
			}
			n++
			if it.Text == "Trading" && n == 4 { // the drawing chart, below its tools
				if h < 300 {
					t.Errorf("Trading: drawing chart is %v tall, want at least 300", h)
				}
				return
			}
			if h != phoneChartHeight {
				t.Errorf("%s: chart %d is %v tall, want %v", it.Text, n, h, phoneChartHeight)
			}
		})
		if n == 0 {
			t.Errorf("%s: no charts found", it.Text)
		}
	}
}
