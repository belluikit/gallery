//go:build !js && !android && !ios

package main

import (
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/test"
)

// The Developer menu's Widget Inspector item opens bellui's inspector in a
// window of its own, and a second use brings that window forward rather than
// opening another (gallery#47).
func TestDeveloperMenuOpensOneInspector(t *testing.T) {
	a := test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	menu := developerMenu(w)
	if len(menu.Items) != 1 || menu.Items[0].Label != "Widget Inspector" {
		t.Fatalf("Developer menu items: %v", menu.Items)
	}
	inspectors := func() int {
		n := 0
		for _, o := range a.Driver().AllWindows() {
			if o.Title() == inspectorTitle {
				n++
			}
		}
		return n
	}
	bellui.DoAndWait(menu.Items[0].Action)
	if got := inspectors(); got != 1 {
		t.Fatalf("%d inspector windows after the item was chosen once, want 1", got)
	}
	bellui.DoAndWait(menu.Items[0].Action)
	if got := inspectors(); got != 1 {
		t.Errorf("%d inspector windows after the item was chosen twice, want 1", got)
	}
}
