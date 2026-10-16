//go:build !android

package pages_test

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"

	"github.com/belluikit/gallery/pages"
)

// TestComponentPagesConstruct builds and renders each component page under
// the test driver, catching construction-time panics and renderer errors.
func TestComponentPagesConstruct(t *testing.T) {
	componentPages := map[string]func(bellui.Window) core.CanvasObject{
		"Animated GIF":  pages.AnimatedGifPage,
		"Charts":        pages.FychartsPage,
		"File Tree":     pages.FileTreePage,
		"Responsive":    pages.ResponsivePage,
		"Scrolling":     pages.ScrollingPage,
		"Code Editor":   pages.FycoderPage,
		"Data Grid":     pages.FytablesPage,
		"File Transfer": pages.UploaderPage,
		"Theme":         pages.ThemePage,
		"DAG Editor":    pages.DAGEditorPage,
	}

	for name, page := range componentPages {
		t.Run(name, func(t *testing.T) {
			a := test.NewApp()
			defer a.Quit()

			w := test.NewWindow(nil)
			defer w.Close()
			w.Resize(geo.NewSize(1200, 800))

			obj := page(w)
			if obj == nil {
				t.Fatalf("%s page returned nil content", name)
			}
			w.SetContent(obj)
		})
	}
}
