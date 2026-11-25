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

// phoneWidth is the width a page must fit without scrolling sideways: a
// 360 dp phone, the narrowest common Android width (gallery#48).
const phoneWidth = 360

// The pages added for gallery#38 fit a phone's width. A card's subtitle is a
// single line, and a long one made the Choice and Pickers pages over 630 dp
// wide, cut off on a phone. The page is measured inside the padding main.go's
// setPage puts round it, so it has 348 dp, not 360.
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
		// Older pages, fitted for gallery#51.
		"File Transfer": pages.UploaderPage,
	} {
		if got := container.NewPadded(page(w)).MinSize().Width; got > phoneWidth {
			t.Errorf("%s page needs %.0f dp with its padding, more than a phone's %d", name, got, phoneWidth)
		}
	}
}
