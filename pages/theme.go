package pages

import (
	"fmt"
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// ThemePage demonstrates widget.ThemeSelector — a callback-driven appearance
// picker. The chosen variant and primary color are applied to the running app
// and stored in the gallery's own Preferences (restored at startup in
// main.go), never in the system-wide Bellui settings. container.ThemeOverride
// draws one part of the page in the other variant (gallery#47).
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

	overrideCard := widget.NewCard("Theme Override", "One part of the window, another theme",
		themeOverrideDemo(a.Settings().ThemeVariant()))
	overrideCard.Style = widget.CardStyleOutlined

	return container.NewVBox(selectorCard, persistCard, overrideCard)
}

// otherVariant is the app's theme with its colours always taken from one
// variant, whatever the app uses.
type otherVariant struct {
	bellui.Theme
	variant bellui.ThemeVariant
}

func (t otherVariant) Color(n bellui.ThemeColorName, _ bellui.ThemeVariant) color.Color {
	return t.Theme.Color(n, t.variant)
}

// themeOverrideDemo shows the same controls twice: with the app's theme, and
// inside a container.ThemeOverride whose theme takes the other variant, so a
// light app shows a dark panel and a dark app a light one.
func themeOverrideDemo(appVariant bellui.ThemeVariant) core.CanvasObject {
	other, name := theme.VariantDark, "dark"
	if appVariant == theme.VariantDark {
		other, name = theme.VariantLight, "light"
	}
	controls := func() core.CanvasObject {
		check := widget.NewCheck("Remember me", nil)
		check.SetChecked(true)
		slider := widget.NewSlider(0, 10)
		slider.SetValue(6)
		return container.NewVBox(
			&widget.Button{Text: "Save", Style: widget.ButtonStyleFilled, OnTapped: func() {}},
			check, slider,
		)
	}
	bg := canvas.NewRectangle(nil)
	bg.CornerRadius = 12
	panel := container.NewThemeOverride(
		container.NewStack(bg, container.NewPadded(controls())),
		otherVariant{Theme: theme.Current(), variant: other},
	)
	// The panel's background is the overridden theme's surface, so the
	// other variant shows as a block, not only in the controls.
	bg.FillColor = otherVariant{Theme: theme.Current(), variant: other}.Color(theme.ColorNameBackground, other)
	return container.NewVBox(
		note("The same controls, then inside a ThemeOverride whose theme takes the "+name+" variant. The override reaches every widget under it."),
		container.NewGridWithColumns(2, container.NewPadded(controls()), panel),
	)
}
