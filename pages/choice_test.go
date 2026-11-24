//go:build !android

package pages_test

import (
	"strings"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/gallery/pages"
)

// walk calls f for every object under o, through containers and widget renderers.
func walk(o core.CanvasObject, f func(core.CanvasObject)) {
	f(o)
	switch v := o.(type) {
	case *bellui.Container:
		for _, c := range v.Objects {
			walk(c, f)
		}
	case core.Widget:
		for _, c := range test.WidgetRenderer(v).Objects() {
			walk(c, f)
		}
	}
}

// The Choice page's three-state group reports what it includes and excludes,
// and the page's box disables every control on it (gallery#38).
func TestChoicePage(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := pages.ChoicePage(w)
	w.SetContent(page)
	w.Resize(geo.NewSize(900, 1400))

	var groups []*widget.ChipGroup
	var disable *widget.Check
	var labels []*widget.Label
	var disableables []core.Disableable
	walk(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *widget.ChipGroup:
			groups = append(groups, v)
		case *widget.Check:
			if v.Text == "Disable everything on this page" {
				disable = v
			}
		case *widget.Label:
			labels = append(labels, v)
		}
		switch o.(type) {
		case *widget.ChipGroup, *widget.CheckGroup, *widget.SpinBox:
			disableables = append(disableables, o.(core.Disableable))
		}
	})
	if len(groups) != 3 || disable == nil {
		t.Fatalf("found %d chip groups and disable box %v", len(groups), disable != nil)
	}
	three := groups[2]
	if !three.Excludable || !three.Multiple {
		t.Fatal("the third group is not the three-state one")
	}

	three.SetSelected([]string{"Jazz"})
	three.SetExcluded([]string{"Pop"})
	status := ""
	for _, l := range labels {
		if strings.HasPrefix(l.Text, "Include:") {
			status = l.Text
		}
	}
	want := "Include: Jazz\nExclude: Pop\nLast change: Pop is now excluded"
	if status != want {
		t.Errorf("status %q, want %q", status, want)
	}

	if len(disableables) != 8 {
		t.Fatalf("found %d controls, want 8 (3 chip groups, 3 check groups, 2 spin boxes)", len(disableables))
	}
	disable.SetChecked(true)
	for _, d := range disableables {
		if !d.Disabled() {
			t.Errorf("%T not disabled by the page's box", d)
		}
	}
	disable.SetChecked(false)
	for _, d := range disableables {
		if d.Disabled() {
			t.Errorf("%T still disabled after unticking", d)
		}
	}
}
