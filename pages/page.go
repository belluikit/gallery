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

// leaveHooks run when the user navigates away from the current page. Pages
// that own something running in the background (a media player decoding
// frames) register a hook with OnLeave to stop it; the gallery drops the old
// page's widgets but cannot know what they started.
var leaveHooks []func()

// OnLeave registers fn to run when the current page is left.
func OnLeave(fn func()) { leaveHooks = append(leaveHooks, fn) }

// LeavePage runs and clears the hooks registered by the current page. main.go
// calls it before building the next page.
func LeavePage() {
	hooks := leaveHooks
	leaveHooks = nil
	for _, fn := range hooks {
		fn()
	}
}
