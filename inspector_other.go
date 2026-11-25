//go:build js || android || ios

package main

import "github.com/belluikit/bellui"

// addDeveloperMenu does nothing here: the widget inspector opens a window of
// its own, and mobile and the web have only one. See inspector_desktop.go.
func addDeveloperMenu(bellui.Window) {}
