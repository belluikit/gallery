//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/theme"

	"github.com/belluikit/gallery/pages"
)

// The Theme page's override panel draws its controls with the other variant
// from the app's: dark in a light app, light in a dark one (gallery#47).
func TestThemeOverridePanel(t *testing.T) {
	for _, appVariant := range []bellui.ThemeVariant{theme.VariantLight, theme.VariantDark} {
		a := test.NewTempApp(t)
		// The test app's own theme has one set of colours for both
		// variants; the gallery's has two.
		a.Settings().SetTheme(theme.DefaultTheme())
		a.Settings().SetThemeVariant(appVariant)
		w := test.NewTempWindow(t, nil)
		page := pages.ThemePage(w)
		var override *container.ThemeOverride
		walk(page, func(o core.CanvasObject) {
			if v, ok := o.(*container.ThemeOverride); ok {
				override = v
			}
		})
		if override == nil {
			t.Fatal("the Theme page has no ThemeOverride")
		}
		other := theme.VariantDark
		if appVariant == theme.VariantDark {
			other = theme.VariantLight
		}
		th := theme.Current()
		got := override.Theme.Color(theme.ColorNameBackground, appVariant)
		if want := th.Color(theme.ColorNameBackground, other); got != want {
			t.Errorf("app variant %v: the panel's background is %v, want the other variant's %v", appVariant, got, want)
		}
		if got == th.Color(theme.ColorNameBackground, appVariant) {
			t.Errorf("app variant %v: the panel draws the app's own background", appVariant)
		}
	}
}
