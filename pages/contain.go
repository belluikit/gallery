package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// ContainmentPage shows Cards, Carousel, Bottom/Side Sheets, Dividers.
func ContainmentPage(w bellui.Window) core.CanvasObject {
	// --- Cards ---
	elevated := widget.NewCard("Elevated Card", "Shadow + surface fill", widget.NewLabel("CardStyleElevated (default)"))
	filled := widget.NewCard("Filled Card", "No shadow, filled surface", widget.NewLabel("CardStyleFilled"))
	filled.Style = widget.CardStyleFilled
	outlined := widget.NewCard("Outlined Card", "Border, no shadow", widget.NewLabel("CardStyleOutlined"))
	outlined.Style = widget.CardStyleOutlined

	cardsCard := widget.NewCard("Cards", "3 MD3 card styles", container.NewGridWithColumns(
		3,
		elevated, filled, outlined,
	))
	cardsCard.Style = widget.CardStyleOutlined

	// --- Carousel ---
	items := []core.CanvasObject{
		widget.NewCard("Slide 1", "", widget.NewLabel("First slide content")),
		widget.NewCard("Slide 2", "", widget.NewLabel("Second slide content")),
		widget.NewCard("Slide 3", "", widget.NewLabel("Third slide content")),
	}
	carousel := widget.NewCarousel(items...)
	carouselCard := widget.NewCard("Carousel", "Swipeable showcase", carousel)
	carouselCard.Style = widget.CardStyleOutlined

	// --- Bottom Sheet ---
	bottomSheetCard := widget.NewCard("Bottom Sheet", "Content from bottom", container.NewVBox(
		widget.NewButton("Show Bottom Sheet", func() {
			sheet := widget.NewBottomSheet(nil)
			content := container.NewVBox(
				widget.NewLabel("Bottom Sheet Content"),
				widget.NewLabel("Add your content here"),
				widget.NewButton("Close", func() { sheet.Hide() }),
			)
			sheet.Content = content
			sheet.Show(w.Canvas())
		}),
	))
	bottomSheetCard.Style = widget.CardStyleOutlined

	// --- Side Sheet ---
	sideSheetCard := widget.NewCard("Side Sheet", "Panel from right", container.NewVBox(
		widget.NewButton("Show Side Sheet", func() {
			content := container.NewVBox(
				widget.NewLabel("Side Sheet Content"),
				widget.NewLabel("Supplementary information"),
			)
			widget.ShowSideSheet(content, w.Canvas())
		}),
	))
	sideSheetCard.Style = widget.CardStyleOutlined

	// --- Divider ---
	dividerCard := widget.NewCard("Divider", "Visual separator", container.NewVBox(
		widget.NewLabel("Content above"),
		widget.NewSeparator(),
		widget.NewLabel("Content below"),
	))
	dividerCard.Style = widget.CardStyleOutlined

	// --- Lists ---
	data := []string{"Inbox", "Starred", "Sent", "Drafts", "Trash"}
	list := widget.NewList(
		func() int { return len(data) },
		func() core.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj core.CanvasObject) {
			obj.(*widget.Label).SetText(data[id])
		},
	)
	list.Resize(geo.NewSize(300, 200))
	listCard := widget.NewCard("Lists", "Scrollable item list", list)
	listCard.Style = widget.CardStyleOutlined

	// --- TreeView (left rail) ---
	tv := newDemoTreeView()
	treeTitle := widget.NewLabelWithStyle("TreeView", bellui.TextAlignLeading, bellui.TextStyle{Bold: true})
	treeSub := widget.NewLabel("Drag reorder · rename · delete")
	treeSub.Wrapping = bellui.TextWrapWord
	treeHeader := container.NewVBox(treeTitle, treeSub, widget.NewSeparator())
	treeRail := container.NewBorder(treeHeader, nil, nil, nil, tv)

	cards := container.NewVBox(cardsCard, carouselCard, bottomSheetCard, sideSheetCard, dividerCard, listCard)
	split := container.NewHSplit(treeRail, cards)
	split.SetOffset(0.25)
	return split
}

func newDemoTreeView() *widget.TreeView {
	tv := widget.NewTreeView()
	tv.AllowReorder = true
	tv.AllowReparent = true
	tv.AllowRename = true
	tv.AllowDelete = true
	tv.ShowGuideLines = true

	_ = tv.AddNode("", widget.TreeViewNode{ID: "src", Text: "src", Icon: theme.FolderIcon()})
	_ = tv.AddNode("src", widget.TreeViewNode{ID: "main.go", Text: "main.go", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("src", widget.TreeViewNode{ID: "utils.go", Text: "utils.go", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("src", widget.TreeViewNode{ID: "widgets", Text: "widgets", Icon: theme.FolderIcon()})
	_ = tv.AddNode("widgets", widget.TreeViewNode{ID: "button.go", Text: "button.go", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("widgets", widget.TreeViewNode{ID: "label.go", Text: "label.go", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("", widget.TreeViewNode{ID: "docs", Text: "docs", Icon: theme.FolderIcon()})
	_ = tv.AddNode("docs", widget.TreeViewNode{ID: "readme", Text: "README.md", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("docs", widget.TreeViewNode{ID: "api", Text: "api", Icon: theme.FolderIcon()})
	_ = tv.AddNode("api", widget.TreeViewNode{ID: "spec", Text: "spec.yaml", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("", widget.TreeViewNode{ID: "tests", Text: "tests", Icon: theme.FolderIcon()})
	_ = tv.AddNode("tests", widget.TreeViewNode{ID: "main_test.go", Text: "main_test.go", Icon: theme.DocumentIcon()})
	_ = tv.AddNode("", widget.TreeViewNode{ID: "go.mod", Text: "go.mod", Icon: theme.FileIcon()})
	_ = tv.AddNode("", widget.TreeViewNode{ID: "go.sum", Text: "go.sum", Icon: theme.FileIcon()})
	_ = tv.AddNode("", widget.TreeViewNode{ID: "Makefile", Text: "Makefile", Icon: theme.FileIcon()})

	tv.OpenBranch("src")
	tv.OpenBranch("docs")

	return tv
}
