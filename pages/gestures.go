package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/layout"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// GesturesPage shows SwipeArea (a card swiped away, and a plain swipe
// detector), ContextMenu and FABMenu (gallery#38).
func GesturesPage(w bellui.Window) core.CanvasObject {
	// --- Swipe card ---
	cardNo := 1
	makeCard := func() core.CanvasObject {
		c := widget.NewCard(fmt.Sprintf("Card %d", cardNo), "Swipe me", widget.NewLabel("Left or right throws me off."))
		return c
	}
	cardStatus := widget.NewLabel("Not swiped yet")
	swipeCard := widget.NewSwipeArea(makeCard(), nil)
	swipeCard.Follow = true
	swipeCard.Directions = map[widget.SwipeDirection]widget.SwipeCardMode{
		widget.SwipeUp:   widget.SwipeStay,   // springs back, still reported
		widget.SwipeDown: widget.SwipeIgnore, // does not move at all
	}
	swipeCard.OnSwiped = func(d widget.SwipeDirection) {
		if d == widget.SwipeUp {
			cardStatus.SetText(fmt.Sprintf("Card %d: swiped up, and it sprang back", cardNo))
			return
		}
		cardStatus.SetText(fmt.Sprintf("Card %d: thrown %s", cardNo, d))
		cardNo++
		swipeCard.Content = makeCard()
		swipeCard.Refresh()
	}

	// --- Plain swipe detector ---
	areaStatus := widget.NewLabel("Swipe across the box")
	plainBox := swipeHint("Swipe here, any way")
	plain := widget.NewSwipeArea(plainBox, func(d widget.SwipeDirection) { areaStatus.SetText("Swiped " + string(d)) })

	swipeCardBox := container.NewGridWrap(geo.NewSize(300, 150), swipeCard)
	plainBoxRow := container.NewGridWrap(geo.NewSize(300, 80), plain)
	swipeSection := widget.NewCard("Swipe Area", "Swipes, with a finger or the mouse", container.NewVBox(
		note("A card that follows the drag. Up springs it back and is only reported; down does nothing."),
		swipeCardBox, cardStatus,
		widget.NewSeparator(),
		note("Without Follow the area stays put and only reports the swipe."),
		plainBoxRow, areaStatus,
	))
	swipeSection.Style = widget.CardStyleOutlined

	// --- Context menu ---
	menuStatus := widget.NewLabel("Nothing chosen")
	choose := func(what string) func() {
		return func() { menuStatus.SetText("Chose " + what + " on report.pdf") }
	}
	file := widget.NewCard("", "", container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("report.pdf")))
	withMenu := widget.NewContextMenu(file,
		bellui.NewMenuItem("Copy", choose("Copy")),
		bellui.NewMenuItem("Rename", choose("Rename")),
		bellui.NewMenuItemSeparator(),
		bellui.NewMenuItem("Delete", choose("Delete")),
	)
	menuCard := widget.NewCard("Context Menu", "A menu on right-click", container.NewVBox(
		note("Right-click the file (a secondary tap). Other taps go to the content as usual."),
		container.NewHBox(withMenu), menuStatus,
	))
	menuCard.Style = widget.CardStyleOutlined

	// --- FAB menu ---
	fabStatus := widget.NewLabel("Nothing created")
	create := func(what string) func() { return func() { fabStatus.SetText("Created a " + what) } }
	speedDial := widget.NewFABMenu(theme.ContentAddIcon(), []*widget.FABMenuItem{
		{Label: "Note", Icon: theme.DocumentCreateIcon(), OnTapped: create("note")},
		{Label: "List", Icon: theme.ListIcon(), OnTapped: create("list")},
		{Label: "Upload", Icon: theme.UploadIcon(), OnTapped: create("upload")},
	})
	var fab *widget.FAB
	fab = widget.NewFAB(theme.ContentAddIcon(), func() {
		// The menu opens above the FAB's top corner on its trailing side:
		// the top-right one, opening up and to the left, or in a
		// right-to-left window the top-left one, opening up and to the
		// right. The card keeps the FAB at the trailing edge, where MD3
		// puts it.
		pos := bellui.CurrentApp().Driver().AbsolutePositionForObject(fab)
		if bellui.CurrentLayoutDirection() != bellui.DirectionRTL {
			pos = pos.AddXY(fab.Size().Width, 0)
		}
		speedDial.Toggle(w.Canvas(), pos)
	})
	fabCard := widget.NewCard("FAB Menu", "A speed dial", container.NewVBox(
		note("The button opens its actions beside it; pick one, or tap the button again to close them."),
		container.NewHBox(fabStatus, layout.NewSpacer(), fab),
	))
	fabCard.Style = widget.CardStyleOutlined

	disable := widget.NewCheck("Disable the swipe areas and the FAB", func(on bool) {
		for _, d := range []core.Disableable{swipeCard, plain, fab} {
			if on {
				d.Disable()
			} else {
				d.Enable()
			}
		}
	})

	return container.NewVBox(disable, swipeSection, menuCard, fabCard)
}

// swipeHint is a framed area with a hint in the middle, for a swipe to
// start on.
func swipeHint(hint string) core.CanvasObject {
	return widget.NewCard("", "", container.NewCenter(widget.NewLabel(hint)))
}
