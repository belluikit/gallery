//go:build !js

package pages

import (
	"testing"

	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"
)

// Show hands the child it replaces to bellui.Release, and only once the slot
// shows the new one: Release frees what the canvas no longer shows, so a call
// made before the swap would keep the old child. A first Show, into an empty
// slot, and showing the same child again release nothing.
func TestShowReleasesTheReplacedChild(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	slot := container.NewStack()
	var released []core.CanvasObject
	var shownAtRelease []core.CanvasObject
	saved := release
	release = func(obj core.CanvasObject) {
		released = append(released, obj)
		shownAtRelease = append(shownAtRelease, slot.Objects[0])
	}
	defer func() { release = saved }()

	first := widget.NewLabel("first")
	second := widget.NewLabel("second")

	Show(slot, first)
	if len(released) != 0 {
		t.Fatalf("Show into an empty slot released %v, want nothing", released)
	}

	Show(slot, second)
	if len(slot.Objects) != 1 || slot.Objects[0] != second {
		t.Fatalf("slot shows %v, want only the second label", slot.Objects)
	}
	if len(released) != 1 || released[0] != first {
		t.Fatalf("released %v, want the first label", released)
	}
	if shownAtRelease[0] != second {
		t.Errorf("Release ran while the slot still showed %v; it must run after the swap", shownAtRelease[0])
	}

	Show(slot, second)
	if len(released) != 1 {
		t.Errorf("showing the shown child again released %v", released[1:])
	}
}
