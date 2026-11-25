package pages

import (
	"fmt"
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// moreDialogs shows the dialogs #38 didn't list (gallery#47): the colour
// picker, simple and advanced, and the three custom dialogs, which take
// content of the app's own.
func moreDialogs(w bellui.Window) core.CanvasObject {
	status := widget.NewLabel("Nothing chosen yet")
	status.Wrapping = bellui.TextWrapWord

	// --- Colour picker ---
	swatch := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	swatch.CornerRadius = 8
	swatch.SetMinSize(geo.NewSize(40, 40))
	picked := func(c color.Color) {
		swatch.FillColor = c
		swatch.Refresh()
		r, g, b, _ := c.RGBA()
		status.SetText(fmt.Sprintf("Colour: #%02X%02X%02X", r>>8, g>>8, b>>8))
	}
	simple := widget.NewButton("Colour Picker", func() {
		dialog.ShowColorPicker("Pick a Colour", "Tap a swatch", picked, w)
	})
	advanced := widget.NewButton("Advanced Colour Picker", func() {
		p := dialog.NewColorPicker("Pick a Colour", "A swatch, or set it exactly", picked, w)
		p.Advanced = true
		p.Show()
		p.SetColor(swatch.FillColor)
	})

	// --- Custom dialog: content of your own, one button that closes it ---
	custom := widget.NewButton("Custom", func() {
		dialog.ShowCustom("Photo Details", "Close", container.NewVBox(
			widget.NewLabel("Lake at dawn.jpg"),
			widget.NewLabel("4032 × 3024, 3.1 MB"),
			widget.NewLabel("Taken 12 May, 05:42"),
		), w)
	})

	// --- Custom confirm: your content, with confirm and dismiss ---
	confirm := widget.NewButton("Custom Confirm", func() {
		name := widget.NewEntry()
		name.SetText("Lake at dawn")
		dialog.ShowCustomConfirm("Rename Photo", "Rename", "Cancel", container.NewVBox(
			widget.NewLabel("New name:"), name,
		), func(ok bool) {
			if ok {
				status.SetText("Renamed to " + name.Text)
			} else {
				status.SetText("Rename cancelled")
			}
		}, w)
	})

	// --- Custom without buttons: the content closes it ---
	noButtons := widget.NewButton("Custom, No Buttons", func() {
		var d *dialog.CustomDialog
		choose := func(what string) func() {
			return func() {
				status.SetText("Shared by " + what)
				d.Hide()
			}
		}
		d = dialog.NewCustomWithoutButtons("Share", container.NewVBox(
			note("No button row: each choice closes the dialog itself."),
			widget.NewButtonWithIcon("Email", theme.MailSendIcon(), choose("email")),
			widget.NewButtonWithIcon("Link", theme.ContentCopyIcon(), choose("link")),
		), w)
		d.Show()
	})

	return container.NewVBox(
		container.NewHBox(container.NewCenter(swatch), container.NewVBox(simple, advanced)),
		widget.NewSeparator(),
		custom, confirm, noButtons,
		status,
	)
}
