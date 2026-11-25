package main

import (
	"slices"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/layout"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// showAppBar shows the gallery's top bar at the top of a window of the given
// width, as main lays it out, and returns it with the window and a count of
// the page rebuilds its actions asked for.
func showAppBar(t *testing.T, width float32) (*widget.TopAppBar, bellui.Window, *int) {
	t.Helper()
	a := test.NewTempApp(t)
	t.Cleanup(func() { bellui.SetLayoutDirection(bellui.DirectionAuto) })
	rebuilds := new(int)
	w := test.NewTempWindow(t, nil)
	bar := newAppBar(a, w, "Communication", func() {}, func() { *rebuilds++ })
	w.SetContent(&bellui.Container{Layout: layout.NewVBoxLayout(), Objects: []core.CanvasObject{bar}})
	w.SetPadded(false)
	w.Resize(geo.NewSize(width, 640))
	return bar, w, rebuilds
}

// barButtons returns the buttons the bar draws after its navigation button,
// in object order, named by their accessible labels.
func barButtons(t *testing.T, bar *widget.TopAppBar) (buttons []*widget.IconButton, names []string) {
	t.Helper()
	for _, o := range test.TempWidgetRenderer(t, bar).Objects() {
		b, ok := o.(*widget.IconButton)
		if !ok || b.Icon == bar.NavigationIcon {
			continue
		}
		node, _ := core.SnapshotAccessible(b)
		buttons, names = append(buttons, b), append(names, node.Label)
	}
	return buttons, names
}

// overflowItems taps the bar's ⋮ button and returns the items of the menu it
// opens, named by their accessible labels.
func overflowItems(t *testing.T, w bellui.Window, overflow *widget.IconButton) (items []core.CanvasObject, names []string) {
	t.Helper()
	test.Tap(overflow)
	// The menu is shown inside an overlay container of bellui's own.
	var menu *widget.PopUpMenu
	for _, o := range test.LaidOutObjects(w.Canvas().Overlays().Top()) {
		if m, ok := o.(*widget.PopUpMenu); ok {
			menu = m
			break
		}
	}
	if menu == nil {
		t.Fatalf("tapping ⋮ opened no *widget.PopUpMenu: the top overlay is %T", w.Canvas().Overlays().Top())
	}
	for _, o := range menu.Items {
		node, _ := core.SnapshotAccessible(o)
		items, names = append(items, o), append(names, node.Label)
	}
	return items, names
}

// assertDistinctIcons checks that no two buttons on the bar draw the same
// icon: the direction toggle used ⋮ in RTL, the overflow button's icon (#66).
func assertDistinctIcons(t *testing.T, buttons []*widget.IconButton) {
	t.Helper()
	seen := map[string]int{}
	for i, b := range buttons {
		if j, dup := seen[b.Icon.Name()]; dup {
			t.Errorf("buttons %d and %d both draw %s", j, i, b.Icon.Name())
		}
		seen[b.Icon.Name()] = i
	}
}

// On a 360 dp phone (the width the gallery targets, #48) the bar keeps the
// variant and direction toggles and a ⋮, whose menu lists the other three
// actions by label and runs them. It stays that way in RTL, where no button
// shares the ⋮ icon either (#66).
func TestAppBarFoldsOnPhone(t *testing.T) {
	bar, w, rebuilds := showAppBar(t, 360)
	variant := "Dark theme" // the test app starts light
	if bellui.CurrentApp().Settings().ThemeVariant() == theme.VariantDark {
		variant = "Light theme"
	}

	buttons, names := barButtons(t, bar)
	if want := []string{variant, "Right-to-left layout", "More options"}; !slices.Equal(names, want) {
		t.Fatalf("at 360 dp the bar shows %q, want %q", names, want)
	}
	assertDistinctIcons(t, buttons)
	if buttons[1].Icon != textDirectionRTLIcon {
		t.Errorf("in LTR the direction toggle draws %s, want the right-to-left icon", buttons[1].Icon.Name())
	}

	_, items := overflowItems(t, w, buttons[2])
	if want := []string{"Dense theme", "Account", "About"}; !slices.Equal(items, want) {
		t.Errorf("the ⋮ menu lists %q, want %q", items, want)
	}
	w.Canvas().Overlays().Remove(w.Canvas().Overlays().Top())

	// Into RTL from the bar: the toggle now offers the way back, with an icon
	// that is still not the ⋮ one.
	test.Tap(buttons[1])
	if *rebuilds != 1 {
		t.Errorf("the direction toggle rebuilt the page %d times, want 1", *rebuilds)
	}
	buttons, names = barButtons(t, bar)
	if want := []string{variant, "Left-to-right layout", "More options"}; !slices.Equal(names, want) {
		t.Fatalf("in RTL the bar shows %q, want %q", names, want)
	}
	assertDistinctIcons(t, buttons)
	if buttons[1].Icon != textDirectionLTRIcon {
		t.Errorf("in RTL the direction toggle draws %s, want the left-to-right icon", buttons[1].Icon.Name())
	}

	// A folded toggle works from the menu. The dense theme's buttons are small
	// enough for all five to fit at 360 dp, so the bar unfolds, and the
	// toggle on it offers the way back, which folds the bar again.
	menu, _ := overflowItems(t, w, buttons[2])
	test.Tap(menu[0].(event.Tappable))
	if *rebuilds != 2 {
		t.Errorf("choosing Dense theme rebuilt the page %d times in all, want 2", *rebuilds)
	}
	buttons, names = barButtons(t, bar)
	if want := []string{variant, "Left-to-right layout", "Default theme", "Account", "About"}; !slices.Equal(names, want) {
		t.Fatalf("under the dense theme the bar shows %q, want %q", names, want)
	}
	assertDistinctIcons(t, buttons)
	test.Tap(buttons[2])
	buttons, names = barButtons(t, bar)
	if want := []string{variant, "Left-to-right layout", "More options"}; !slices.Equal(names, want) {
		t.Fatalf("back on the default theme the bar shows %q, want %q", names, want)
	}

	// Account and About open their dialogs from the menu: choosing one closes
	// the menu, so what is on top afterwards is the dialog.
	for i, name := range []string{"Account", "About"} {
		menu, items := overflowItems(t, w, buttons[2])
		if items[i+1] != name {
			t.Fatalf("menu item %d is %q, want %q", i+1, items[i+1], name)
		}
		test.Tap(menu[i+1].(event.Tappable))
		dialog := w.Canvas().Overlays().Top()
		if dialog == nil {
			t.Errorf("choosing %s from the menu opened nothing", name)
			continue
		}
		for _, o := range test.LaidOutObjects(dialog) {
			if _, ok := o.(*widget.PopUpMenu); ok {
				t.Errorf("after choosing %s the menu is still on top", name)
			}
		}
		w.Canvas().Overlays().Remove(dialog)
	}
}

// On a desktop-wide bar all five actions show as icons, with no ⋮.
func TestAppBarUnfoldedOnDesktop(t *testing.T) {
	bar, _, _ := showAppBar(t, 900)
	buttons, names := barButtons(t, bar)
	if len(names) != 5 || names[4] != "About" {
		t.Errorf("at 900 dp the bar shows %q, want the five actions ending in About", names)
	}
	assertDistinctIcons(t, buttons)
}
