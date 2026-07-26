package pages

import (
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
)

// pageMinSize wraps obj so it reports at least the given minimum size.
// Component pages host widgets that manage their own scrolling (grids,
// charts, editors) and would otherwise collapse inside the gallery's
// outer vertical scroll container.
//
// On mobile the desktop-class minimum width (these pages ask for 900–1100px)
// would force horizontal overflow on a ~360–410dp phone screen. There we drop
// the forced width and let each widget use its own minimum, relying on the
// gallery's bidirectional scroll (see main.go) so wide content pans rather
// than clips. The height floor is kept so pages don't collapse.
func pageMinSize(w, h float32, obj core.CanvasObject) core.CanvasObject {
	if bellui.CurrentDevice().IsMobile() {
		w = 0
	}
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(geo.NewSize(w, h))
	return container.NewStack(spacer, obj)
}
