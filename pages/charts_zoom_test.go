//go:build !android

package pages

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"

	"github.com/belluikit/bellui/x/charts/chart"
	"github.com/belluikit/bellui/x/charts/component"
)

// walkPage calls f for every object under o: through containers, and through
// every tab of an AppTabs, shown or not. It does not enter widgets.
func walkPage(o core.CanvasObject, f func(core.CanvasObject)) {
	f(o)
	switch v := o.(type) {
	case *bellui.Container:
		for _, c := range v.Objects {
			walkPage(c, f)
		}
	case *container.AppTabs:
		for _, item := range v.Items {
			walkPage(item.Content, f)
		}
	}
}

// The Charts page's zoom slider drives the dashed-lines chart, and the chart
// moves it back when a pinch zooms it (gallery#47).
func TestChartsZoomSliderDrivesItsChart(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := FychartsPage(w)

	var slider *component.DataZoomSlider
	var driven *chart.Chart
	walkPage(page, func(o core.CanvasObject) {
		c, ok := o.(*bellui.Container)
		if !ok {
			return
		}
		var s *component.DataZoomSlider
		var ch *chart.Chart
		for _, child := range c.Objects {
			switch v := child.(type) {
			case *component.DataZoomSlider:
				s = v
			case *chart.Chart:
				ch = v
			}
		}
		if s != nil && ch != nil {
			slider, driven = s, ch
		}
	})
	if slider == nil || driven == nil {
		t.Fatalf("no zoom slider beside a chart: slider %v, chart %v", slider != nil, driven != nil)
	}
	driven.Resize(geo.NewSize(400, 300))
	slider.Resize(geo.NewSize(400, 30))

	// Drag the right handle from the end to the middle.
	drag := func(x float32) {
		slider.Dragged(&event.DragEvent{PointEvent: event.PointEvent{Position: geo.NewPos(x, 15)}})
	}
	drag(398)
	drag(200)
	slider.DragEnd()
	if got := driven.ZoomWindow(); got.Start != 0 || got.End < 0.49 || got.End > 0.51 {
		t.Errorf("after dragging the right handle to the middle the chart shows %+v, want 0 to 0.5", got)
	}

	// A pinch on the chart moves the slider to the chart's new window.
	driven.Pinched(&event.PinchEvent{PointEvent: event.PointEvent{Position: geo.NewPos(200, 150)}, Scale: 2})
	if got, want := slider.Window(), driven.ZoomWindow(); got != want {
		t.Errorf("after a pinch the slider shows %+v, the chart %+v", got, want)
	}
	if slider.Window().Span() > 0.3 {
		t.Errorf("the pinch zoomed nothing: the slider's window is %+v", slider.Window())
	}
}
