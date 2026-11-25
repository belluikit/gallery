//go:build !js && !android && !ios

package main

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/dev/inspector"
)

// developerMenu is the desktop's Developer menu. Its Widget Inspector item
// opens bellui's read-only widget-tree inspector on w, in a window of its
// own, with pick mode (gallery#47). That is a different tool from the F12
// overlay dev.Enable adds. The inspector needs a second window, which
// mobile and the web don't have, so other builds have no such menu.
func developerMenu(w bellui.Window) *bellui.Menu {
	return bellui.NewMenu("Developer",
		bellui.NewMenuItem("Widget Inspector", func() {
			// One inspector at a time: bring an open one forward. The
			// inspector doesn't say whether its window is still open, so look
			// for the window by its title.
			for _, o := range bellui.CurrentApp().Driver().AllWindows() {
				if o.Title() == inspectorTitle {
					o.RequestFocus()
					return
				}
			}
			inspector.Show(w)
		}),
	)
}

// inspectorTitle is the title inspector.Show gives its window.
const inspectorTitle = "Bellui Inspector"

// addDeveloperMenu gives w a main menu with the Developer menu.
func addDeveloperMenu(w bellui.Window) {
	w.SetMainMenu(bellui.NewMainMenu(developerMenu(w)))
}
