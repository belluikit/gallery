package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
)

// PageDef describes a single gallery destination: the label and icon shown in
// the navigation drawer, and a constructor that builds the page content for a
// given window.
type PageDef struct {
	Name string
	Icon bellui.Resource
	New  func(bellui.Window) core.CanvasObject
	// Fill hosts the page directly in the content area, without the gallery's
	// outer scroll. Use it for full-bleed interactive widgets that manage their
	// own scrolling and pointer handling (e.g. the code editor) — wrapping such a
	// widget in an outer container.Scroll lets the scroll intercept presses and
	// drags, breaking focus, gutter taps, and modifier-clicks.
	Fill bool
}
