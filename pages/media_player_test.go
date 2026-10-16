//go:build !android && !js

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"

	"github.com/belluikit/gallery/pages"
)

// The Media Player page builds on the bundled sample, and leaving the page
// closes the player (it registers an OnLeave hook).
func TestMediaPlayerPageBuildsAndLeaves(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(900, 600))

	obj := pages.MediaPlayerPage(w)
	if obj == nil {
		t.Fatal("Media Player page returned nil content")
	}
	w.SetContent(obj)
	pages.LeavePage() // must close the player without panicking
	pages.LeavePage() // and running it again is a no-op
}
