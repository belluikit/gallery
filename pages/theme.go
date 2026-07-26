package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"
)

// ThemePage demonstrates widget.ThemeSelector — a callback-driven appearance
// picker. The chosen variant and primary color are applied to the running app
// and stored in the gallery's own Preferences (restored at startup in
// main.go), never in the system-wide Bellui settings.
func ThemePage(_ bellui.Window) core.CanvasObject {
	a := bellui.CurrentApp()

	status := widget.NewLabel("")
	showSelection := func(s widget.ThemeSelection) {
		status.SetText(fmt.Sprintf("Variant: %v  •  Primary: %v", s.Variant, s.Primary))
	}

	sel := widget.NewThemeSelector(nil)
	if stored, ok := widget.LoadThemeSelection(a.Preferences()); ok {
		// Show the stored choice; assigning the field (rather than calling
		// SetSelection) avoids re-firing Apply on every page open.
		sel.Selection = stored
	}
	showSelection(sel.Selection)

	sel.OnChanged = func(s widget.ThemeSelection) {
		s.Apply()
		widget.SaveThemeSelection(s, a.Preferences())
		showSelection(s)
	}

	selectorCard := widget.NewCard("Theme Selector",
		"System / Light / Dark variant and primary color, applied live",
		container.NewVBox(sel, status))

	persistCard := widget.NewCard("Per-app persistence",
		"Stored in the gallery's own Preferences",
		widget.NewLabel("Pick a variant and color, then relaunch — the gallery restores your\n"+
			"appearance at startup without touching the system-wide Bellui settings."))

	return container.NewVBox(selectorCard, persistCard)
}
