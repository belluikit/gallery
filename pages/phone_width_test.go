//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/theme"

	"github.com/belluikit/gallery/pages"
)

// phoneWidth is the width a page must fit without scrolling sideways: a
// 412 dp phone (Pixel). Since bellui v1.4.0 its pickers fit a 360 dp phone
// (the DateTimePicker needs 344 dp, belluikit/bellui#1421), but the gallery's
// own cards don't yet: its fixed 340 dp demo boxes and the DateTimePicker's
// card need 370–374 dp (gallery#48).
const phoneWidth = 412

// The pages added for gallery#38 fit a phone's width: a card's subtitle is
// a single line, and a long one made the Choice and Pickers pages over 630 dp
// wide, cut off on a phone.
func TestPagesFitAPhone(t *testing.T) {
	a := test.NewTempApp(t)
	a.Settings().SetTheme(theme.DefaultTheme())
	w := test.NewTempWindow(t, nil)
	for name, page := range map[string]func(bellui.Window) core.CanvasObject{
		"Choice":     pages.ChoicePage,
		"Pickers":    pages.PickersPage,
		"Text":       pages.TextPage,
		"Gestures":   pages.GesturesPage,
		"Structure":  pages.StructurePage,
		"Containers": pages.ContainersPage,
	} {
		if got := page(w).MinSize().Width; got > phoneWidth {
			t.Errorf("%s page needs %.0f dp, more than a phone's %d", name, got, phoneWidth)
		}
	}
}
