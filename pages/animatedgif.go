package pages

import (
	_ "embed"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// earthGif is the rotating-earth test image from the bellui repository
// (CC BY-SA 3.0, Wikipedia user Marvel, derived from NASA imagery).
//
//go:embed earth.gif
var earthGif []byte

// AnimatedGifPage demos widget.AnimatedGif (fyne #443): frame playback driven
// by the canvas animation ticker, honoring the file's own frame delays, and
// widget.TwoStateToolbarAction (fyne #446) as its play/pause control.
func AnimatedGifPage(_ bellui.Window) core.CanvasObject {
	gif, err := widget.NewAnimatedGifFromResource(bellui.NewStaticResource("earth.gif", earthGif))
	if err != nil {
		return widget.NewLabel("Failed to load earth.gif: " + err.Error())
	}
	gif.SetMinSize(geo.NewSize(280, 280))

	status := widget.NewLabelWithStyle("playing", bellui.TextAlignCenter, bellui.TextStyle{Italic: true})

	// TwoStateToolbarAction rides shotgun: a play/pause ToolbarItem.
	var playPause *widget.TwoStateToolbarAction
	playPause = widget.NewTwoStateToolbarAction(theme.MediaPauseIcon(), theme.MediaPlayIcon(),
		func(paused bool) {
			if paused {
				gif.Stop()
				status.SetText("stopped — the last frame stays visible")
			} else {
				gif.Start()
				status.SetText("playing")
			}
		})
	toolbar := widget.NewToolbar(
		playPause,
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			gif.Stop()
			gif.Start()
			playPause.SetOn(false)
			status.SetText("playing")
		}),
	)

	note := widget.NewLabel("The gif plays at its own cadence (44 frames × 90 ms here) on the canvas " +
		"animation ticker — no goroutine, and playback pauses with the window. The toolbar's " +
		"play/pause button is a TwoStateToolbarAction, also new in this release.")
	note.Wrapping = bellui.TextWrapWord

	gif.Start()
	card := widget.NewCard("Animated GIF", "widget.AnimatedGif + TwoStateToolbarAction",
		container.NewVBox(toolbar, container.NewCenter(gif), status))
	card.Style = widget.CardStyleOutlined

	return container.NewVBox(note, card)
}
