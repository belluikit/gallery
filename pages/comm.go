package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// CommunicationPage shows Badges, Snackbar, Tooltips, Progress, Loading, Toast.
func CommunicationPage(w bellui.Window) core.CanvasObject {
	// --- Badges ---
	dot := widget.NewBadge(widget.NewIconButton(theme.MailComposeIcon(), func() {}))
	count := widget.NewBadgeWithCount(widget.NewIconButton(theme.InfoIcon(), func() {}), 3)
	large := widget.NewBadgeWithCount(widget.NewIconButton(theme.WarningIcon(), func() {}), 999)

	badgeCard := widget.NewCard("Badges", "Dot and count indicators", container.NewHBox(
		container.NewVBox(widget.NewLabel("Dot"), dot),
		container.NewVBox(widget.NewLabel("Count"), count),
		container.NewVBox(widget.NewLabel("Large"), large),
	))
	badgeCard.Style = widget.CardStyleOutlined

	// --- Snackbar ---
	snackLabel := widget.NewLabel("—")
	snackCard := widget.NewCard("Snackbar", "Brief message at bottom", container.NewVBox(
		widget.NewButton("Show Snackbar", func() {
			widget.ShowSnackbar("Changes saved", "", nil, w.Canvas())
		}),
		widget.NewButton("Snackbar with Action", func() {
			widget.ShowSnackbar("Item archived", "Undo", func() {
				snackLabel.SetText("Undo clicked!")
			}, w.Canvas())
		}),
		snackLabel,
	))
	snackCard.Style = widget.CardStyleOutlined

	// --- Progress ---
	progress := widget.NewProgressBar()
	progress.SetValue(0.65)
	infinite := widget.NewProgressBarInfinite()
	activity := widget.NewActivity()
	activity.Start()

	progressCard := widget.NewCard("Progress Indicators", "Linear and circular", container.NewVBox(
		widget.NewLabel("Determinate (65%)"),
		progress,
		widget.NewLabel("Indeterminate"),
		infinite,
		container.NewHBox(widget.NewLabel("Activity"), activity),
	))
	progressCard.Style = widget.CardStyleOutlined

	// --- Loading ---
	loadRect := widget.NewLoadingIndicator(widget.LoadingShapeRect)
	loadCircle := widget.NewLoadingIndicator(widget.LoadingShapeCircle)
	loadText := widget.NewLoadingIndicator(widget.LoadingShapeText)

	loadCard := widget.NewCard("Loading Indicators", "Skeleton placeholders", container.NewHBox(
		container.NewVBox(widget.NewLabel("Card"), loadRect),
		container.NewVBox(widget.NewLabel("Avatar"), loadCircle),
		container.NewVBox(widget.NewLabel("Text"), loadText),
	))
	loadCard.Style = widget.CardStyleOutlined

	// --- Tooltip ---
	tooltipTarget := widget.NewTooltip(widget.NewButton("Hover me", func() {}), "This is a tooltip")
	tooltipCard := widget.NewCard("Tooltip", "Hover text after delay", tooltipTarget)
	tooltipCard.Style = widget.CardStyleOutlined

	// --- Dialogs ---
	dialogCard := widget.NewCard("Dialogs", "Confirm, information, entry, file open/save", container.NewVBox(
		widget.NewButton("Confirm Dialog", func() {
			dialog.ShowConfirm("Confirm", "Are you sure you want to proceed?", func(ok bool) {}, w)
		}),
		widget.NewButton("Information Dialog", func() {
			dialog.ShowInformation("Information", "This is an informational message.", w)
		}),
		widget.NewButton("Entry Dialog", func() {
			entry := widget.NewEntry()
			dialog.ShowForm("Input", "OK", "Cancel",
				[]*widget.FormItem{widget.NewFormItem("Enter your name:", entry)},
				func(bool) {}, w)
		}),
		widget.NewButton("File Open Dialog", func() {
			dialog.ShowFileOpen(func(reader bellui.URIReadCloser, err error) {
				if reader != nil {
					reader.Close()
				}
			}, w)
		}),
		widget.NewButton("File Save Dialog", func() {
			dialog.ShowFileSave(func(writer bellui.URIWriteCloser, err error) {
				if writer != nil {
					writer.Close()
				}
			}, w)
		}),
		widget.NewButton("Folder Open Dialog", func() {
			dialog.ShowFolderOpen(func(uri bellui.ListableURI, err error) {}, w)
		}),
	))
	dialogCard.Style = widget.CardStyleOutlined

	// --- Toast ---
	toastCard := widget.NewCard("Toast", "Transient notification", container.NewVBox(
		widget.NewButton("Show Toast", func() {
			widget.ShowToast("Action completed", w.Canvas())
		}),
	))
	toastCard.Style = widget.CardStyleOutlined

	return container.NewVBox(badgeCard, snackCard, progressCard, loadCard, tooltipCard, dialogCard, toastCard)
}
