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

// release lets go of content a slot no longer shows. A variable so tests can
// see what Show passes it: bellui's test driver has no Release.
var release = bellui.Release

// Show makes obj the only child of slot and refreshes it. If slot showed
// something else, Show then hands that to bellui.Release, so the caches let go
// of it at the next frame. Without the release, content dropped by assigning
// Objects stays cached until a frame at least a minute later, which an idle
// app never draws (belluikit/bellui#1083). The release comes after the
// assignment because it frees what the canvas no longer shows.
func Show(slot *bellui.Container, obj core.CanvasObject) {
	var old core.CanvasObject
	if len(slot.Objects) > 0 {
		old = slot.Objects[0]
	}
	slot.Objects = []core.CanvasObject{obj}
	slot.Refresh()
	if old != nil && old != obj {
		release(old)
	}
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
