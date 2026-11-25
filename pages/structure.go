package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/storage"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// StructurePage shows the widgets that arrange other content: Accordion, the
// classic Table, Tree and GridWrap, ScrollIndicator, FileIcon and StatusBar
// (gallery#38). The Data Grid page shows x/tables' DataGrid.
func StructurePage(_ bellui.Window) core.CanvasObject {
	// --- Status bar, updated by the cards above it ---
	statusMsg := widget.NewStatusBarLabelWithIcon("Ready", theme.InfoIcon())
	statusSel := widget.NewStatusBarLabel("Nothing selected")
	status := widget.NewStatusBar(statusMsg, widget.NewStatusBarSeparator(), statusSel, widget.NewStatusBarSpacer(),
		widget.NewStatusBarLabel("UTF-8"))

	// --- Accordion ---
	accordion := widget.NewAccordion(
		widget.NewAccordionItem("Shipping", note("Orders over $50 ship free. Others ship for $4.99.")),
		widget.NewAccordionItem("Returns", note("Return anything within 30 days of delivery.")),
		widget.NewAccordionItem("Warranty", note("Two years on all parts and labour.")),
	)
	accordion.Open(0)
	multi := widget.NewCheck("Let several sections be open", func(on bool) {
		accordion.MultiOpen = on
		if !on {
			accordion.CloseAll()
			accordion.Open(0)
		}
	})
	accordionCard := widget.NewCard("Accordion", "Sections that expand", container.NewVBox(multi, accordion))
	accordionCard.Style = widget.CardStyleOutlined

	// --- Table ---
	planets := [][]string{
		{"Mercury", "0.39", "88", "0"},
		{"Venus", "0.72", "225", "0"},
		{"Earth", "1.00", "365", "1"},
		{"Mars", "1.52", "687", "2"},
		{"Jupiter", "5.20", "4,333", "95"},
		{"Saturn", "9.54", "10,759", "146"},
	}
	headers := []string{"Planet", "Distance (AU)", "Year (days)", "Moons"}
	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(planets), len(headers) },
		func() core.CanvasObject { return widget.NewLabel("Distance (AU)") },
		func(id widget.TableCellID, o core.CanvasObject) { o.(*widget.Label).SetText(planets[id.Row][id.Col]) },
	)
	table.ShowHeaderColumn = false
	table.StickyColumnCount = 1
	table.UpdateHeader = func(id widget.TableCellID, o core.CanvasObject) {
		if id.Row < 0 && id.Col >= 0 {
			o.(*widget.Label).SetText(headers[id.Col])
		}
	}
	table.OnSelected = func(id widget.TableCellID) {
		statusSel.SetText(fmt.Sprintf("%s: %s", planets[id.Row][0], headers[id.Col]))
	}
	tableCard := widget.NewCard("Table", "Rows and columns that scroll", container.NewVBox(
		note("A header row, and the Planet column stays put while the rest scrolls sideways. Select a cell: the status bar at the bottom names it."),
		container.NewGridWrap(geo.NewSize(300, 260), table),
	))
	tableCard.Style = widget.CardStyleOutlined

	// --- Tree ---
	tree := widget.NewTreeWithStrings(map[string][]string{
		"":          {"Documents", "Pictures", "Music"},
		"Documents": {"Report.pdf", "Notes.txt"},
		"Pictures":  {"Holiday", "Profile.png"},
		"Holiday":   {"Beach.jpg", "Hills.jpg"},
		"Music":     {"Song.mp3"},
	})
	tree.OpenBranch("Documents")
	tree.OnSelected = func(id widget.TreeNodeID) { statusSel.SetText("Tree: " + id) }
	treeCard := widget.NewCard("Tree", "Branches that open and close", container.NewVBox(
		container.NewGridWrap(geo.NewSize(300, 240), tree),
	))
	treeCard.Style = widget.CardStyleOutlined

	// --- Grid wrap ---
	icons := []bellui.Resource{
		theme.HomeIcon(), theme.MailComposeIcon(), theme.SearchIcon(), theme.SettingsIcon(), theme.AccountIcon(),
		theme.DocumentIcon(), theme.FolderIcon(), theme.MediaPlayIcon(), theme.MediaMusicIcon(), theme.MediaPhotoIcon(),
		theme.StorageIcon(), theme.ComputerIcon(), theme.ColorPaletteIcon(), theme.HistoryIcon(), theme.InfoIcon(),
		theme.HelpIcon(), theme.UploadIcon(), theme.DownloadIcon(), theme.ContentCutIcon(), theme.ContentCopyIcon(),
	}
	grid := widget.NewGridWrap(
		func() int { return len(icons) },
		func() core.CanvasObject {
			return container.NewVBox(widget.NewIcon(theme.HomeIcon()), widget.NewLabelWithStyle("00", bellui.TextAlignCenter, bellui.TextStyle{}))
		},
		func(id widget.GridWrapItemID, o core.CanvasObject) {
			box := o.(*bellui.Container)
			box.Objects[0].(*widget.Icon).SetResource(icons[id])
			box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id+1))
		},
	)
	grid.OnSelected = func(id widget.GridWrapItemID) { statusSel.SetText(fmt.Sprintf("Grid: tile %d", id+1)) }
	gridCard := widget.NewCard("Grid Wrap", "Tiles that wrap and scroll", container.NewVBox(
		note("A list laid out in rows: as many tiles fit across as the width allows, and only the visible ones are built."),
		container.NewGridWrap(geo.NewSize(300, 200), grid),
	))
	gridCard.Style = widget.CardStyleOutlined

	// --- Scroll indicator ---
	var rows []core.CanvasObject
	for i := 1; i <= 60; i++ {
		rows = append(rows, widget.NewLabel(fmt.Sprintf("Line %d of a long document", i)))
	}
	scroll := container.NewVScroll(container.NewVBox(rows...))
	indicator := widget.NewScrollIndicator(scroll)
	indicator.Placement = widget.ScrollIndicatorTopRight
	indicatorCard := widget.NewCard("Scroll Indicator", "Where you are in a long scroll", container.NewVBox(
		note("A small map laid over the scroll area; its box follows the visible part as you scroll."),
		container.NewGridWrap(geo.NewSize(300, 220), container.NewStack(scroll, indicator)),
	))
	indicatorCard.Style = widget.CardStyleOutlined

	// --- File icons ---
	var files []core.CanvasObject
	for _, name := range []string{"report.pdf", "photo.png", "song.mp3", "notes.txt", "archive.zip"} {
		icon := widget.NewFileIcon(storage.NewFileURI("/tmp/" + name))
		files = append(files, container.NewVBox(
			container.NewGridWrap(geo.NewSize(48, 48), icon),
			widget.NewLabelWithStyle(name, bellui.TextAlignCenter, bellui.TextStyle{}),
		))
	}
	fileCard := widget.NewCard("File Icon", "An icon for a file's type", container.NewVBox(
		container.NewGridWrap(geo.NewSize(96, 84), files...),
	))
	fileCard.Style = widget.CardStyleOutlined

	statusCard := widget.NewCard("Status Bar", "A strip of short messages", container.NewVBox(
		note("Labels, separators and a spacer; the selection label follows the table, tree and grid above."),
		status,
	))
	statusCard.Style = widget.CardStyleOutlined

	return container.NewVBox(accordionCard, tableCard, treeCard, gridCard, indicatorCard, fileCard, statusCard)
}
