package pages

import (
	"fmt"
	"strings"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"
)

// ScrollingPage demos the fyne #449 scroll bar grab targets interactively:
// the minimum thumb length on very long content, the always-wide invisible
// hit area (the bar expands before the pointer reaches the thin strip), and
// the constant reserved gutter next to word-wrapped content.
func ScrollingPage(_ bellui.Window) core.CanvasObject {
	intro := widget.NewLabel("Things to try: move the pointer toward a right edge — the bar expands " +
		"while you are still ~9dp away (the hit area is always full width, only the visual is thin). " +
		"Grab and drag the thumb: even at 50,000 rows it stays a 36dp handle instead of a 12dp nub.")
	intro.Wrapping = bellui.TextWrapWord

	// --- 50,000-row list: the original fycoder complaint scenario ---
	list := widget.NewList(
		func() int { return 50000 },
		func() core.CanvasObject { return widget.NewLabel("Row template") },
		func(i widget.ListItemID, o core.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("Row %d of a very long list", i+1))
		},
	)
	listCard := widget.NewCard("50,000 rows", "The thumb keeps its minimum length; drag maps the full range", list)
	listCard.Style = widget.CardStyleOutlined

	// --- 5,000-line editor, wrap off: both bars active ---
	editor := widget.NewMultiLineEntry()
	editor.Wrapping = bellui.TextWrapOff
	var b strings.Builder
	for i := 1; i <= 5000; i++ {
		fmt.Fprintf(&b, "%5d  A long unwrapped editor line to give the horizontal bar something to do as well.\n", i)
	}
	editor.SetText(b.String())
	editorCard := widget.NewCard("5,000-line editor", "Vertical and horizontal bars; hover either edge", editor)
	editorCard.Style = widget.CardStyleOutlined

	// --- Word-wrapped text: the reserved gutter, toggleable ---
	wrapped := widget.NewRichTextFromMarkdown(strings.Repeat(
		"Word-wrapped paragraphs keep a constant gutter for the scroll bar, so the text never sits "+
			"under the expanded bar and never re-flows when the bar appears. ", 60,
	))
	wrapped.Wrapping = bellui.TextWrapWord
	wrapScroll := container.NewScroll(wrapped)

	reserve := widget.NewSwitch(func(on bool) {
		wrapScroll.ReserveBarSpace = on
		wrapScroll.Refresh()
	})
	reserve.On = true
	reserveRow := container.NewHBox(widget.NewLabel("Reserve bar space"), reserve)
	wrapCard := widget.NewCard("Wrapped text + reserved gutter",
		"Toggle the reservation to see the bar overlap the text's edge",
		container.NewBorder(reserveRow, nil, nil, nil, wrapScroll))
	wrapCard.Style = widget.CardStyleOutlined

	grid := container.NewGridWithColumns(3, listCard, editorCard, wrapCard)
	return container.NewBorder(intro, nil, nil, nil, grid)
}
