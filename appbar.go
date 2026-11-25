package main

import (
	_ "embed"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
	"github.com/belluikit/gallery/pages"
)

// The top bar's actions, in bar order. On a narrow bar, bellui folds the
// labelled actions into a ⋮ menu starting from the last, so the ones a viewer
// flips while browsing pages come first and the two dialogs last (#66).
const (
	actionVariant   = iota // light or dark theme
	actionDirection        // left-to-right or right-to-left layout
	actionDensity          // default or dense theme
	actionAccount          // the sign-in dialog
	actionAbout            // the about dialog
)

// The direction toggle's icons: Material's format_textdirection_l_to_r and
// format_textdirection_r_to_l (Apache 2.0; see THIRD_PARTY_NOTICES.md).
// bellui has no text-direction icon. The ⋯ and ⋮ the toggle used before read
// as "more", and ⋮ is the bar's own overflow button (#66). These icons are not
// mirrored in RTL: they show a direction.
var (
	//go:embed icons/format-textdirection-l-to-r.svg
	textDirectionLTRSVG []byte
	//go:embed icons/format-textdirection-r-to-l.svg
	textDirectionRTLSVG []byte

	textDirectionLTRIcon = theme.NewThemedResource(&bellui.StaticResource{
		StaticName: "format-textdirection-l-to-r.svg", StaticContent: textDirectionLTRSVG,
	})
	textDirectionRTLIcon = theme.NewThemedResource(&bellui.StaticResource{
		StaticName: "format-textdirection-r-to-l.svg", StaticContent: textDirectionRTLSVG,
	})
)

// newAppBar builds the gallery's top bar: a navigation button calling onNav,
// the page title, and five actions. rebuild re-creates the current page after
// a theme or direction change.
//
// Each action has a Label, which names its button for screen readers and lets
// it fold into the ⋮ menu on a phone. The three toggles show what a tap does:
// the variant toggle in dark mode reads "Light theme" with the chromatic icon.
func newAppBar(a bellui.App, w bellui.Window, title string, onNav, rebuild func()) *widget.TopAppBar {
	onUI := func(fn func()) { bellui.Do(fn) } // as in main: taps can arrive off the render goroutine

	bar := widget.NewTopAppBar(title, theme.MenuIcon(), onNav)
	actions := make([]widget.TopAppBarAction, actionAbout+1)

	// Seeded from the current variant, since a stored Theme-page selection
	// may already have applied light mode.
	setVariant := func(v bellui.ThemeVariant) {
		if v == theme.VariantDark {
			actions[actionVariant].Icon, actions[actionVariant].Label = theme.ColorChromaticIcon(), "Light theme"
		} else {
			actions[actionVariant].Icon, actions[actionVariant].Label = theme.ColorAchromaticIcon(), "Dark theme"
		}
	}
	setVariant(a.Settings().ThemeVariant())
	actions[actionVariant].OnTapped = func() {
		onUI(func() {
			v := theme.VariantDark
			if a.Settings().ThemeVariant() == theme.VariantDark {
				v = theme.VariantLight
			}
			a.Settings().SetThemeVariant(v)
			setVariant(v)
			bar.Refresh()
			rebuild()
		})
	}

	// Toggle the process-wide layout direction to preview the whole gallery
	// right-to-left. SetLayoutDirection re-lays-out and repaints the running
	// window (mirroring containers, alignment, scroll, focus order and
	// directional icons); rebuild() refreshes the current page to match.
	isRTL := false
	setDirection := func() {
		if isRTL {
			actions[actionDirection].Icon, actions[actionDirection].Label = textDirectionLTRIcon, "Left-to-right layout"
		} else {
			actions[actionDirection].Icon, actions[actionDirection].Label = textDirectionRTLIcon, "Right-to-left layout"
		}
	}
	setDirection()
	actions[actionDirection].OnTapped = func() {
		onUI(func() {
			isRTL = !isRTL
			if isRTL {
				bellui.SetLayoutDirection(bellui.DirectionRTL)
			} else {
				bellui.SetLayoutDirection(bellui.DirectionLTR)
			}
			setDirection()
			bar.Refresh()
			rebuild()
		})
	}

	isDense := false
	setDensity := func() {
		actions[actionDensity].Label = "Dense theme"
		if isDense {
			actions[actionDensity].Label = "Default theme"
		}
	}
	setDensity()
	actions[actionDensity].Icon = theme.GridIcon()
	actions[actionDensity].OnTapped = func() {
		onUI(func() {
			isDense = !isDense
			if isDense {
				a.Settings().SetTheme(theme.DenseTheme())
			} else {
				a.Settings().SetTheme(theme.DefaultTheme())
			}
			setDensity()
			bar.Refresh()
			rebuild()
		})
	}

	actions[actionAccount] = widget.TopAppBarAction{
		Icon:     theme.AccountIcon(),
		Label:    "Account",
		OnTapped: func() { onUI(func() { pages.ShowAuthDialog(w) }) },
	}
	// dialog.About (fyne #446/#453): MD3-native about box built from the
	// app metadata, with the parallax icon underlay behind the content.
	actions[actionAbout] = widget.TopAppBarAction{
		Icon:     theme.InfoIcon(),
		Label:    "About",
		OnTapped: func() { onUI(func() { showAbout(a, w) }) },
	}

	bar.Actions = actions
	return bar
}
