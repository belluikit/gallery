package pages

import (
	"os"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/storage"
	"github.com/belluikit/bellui/widget"
)

// FileTreePage demos widget.FileTree (fyne #445): a filesystem browser with
// asynchronous listings, hidden-file filtering and live file-system watching
// (fyne #454) — create or delete a file in the shown directory from another
// app and the branch re-lists itself.
func FileTreePage(_ bellui.Window) core.CanvasObject {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.TempDir()
	}
	root := storage.NewFileURI(home)

	treeSlot := container.NewStack()
	selected := widget.NewLabel("Select a file…")
	selected.Wrapping = bellui.TextWrapWord

	showHidden := false
	buildTree := func() {
		tree := widget.NewFileTree(root)
		tree.ShowHidden = showHidden
		tree.ShowRootPath = true
		tree.Watch = true // live: external changes re-list the open branches
		tree.OnSelected = func(id widget.TreeNodeID) {
			selected.SetText(id)
		}
		treeSlot.Objects = []core.CanvasObject{tree}
		treeSlot.Refresh()
	}

	hiddenSwitch := widget.NewSwitch(func(on bool) {
		showHidden = on
		buildTree()
	})
	buildTree()

	controls := container.NewHBox(
		widget.NewLabel("Show hidden files"), hiddenSwitch,
		widget.NewLabel("   Watching: try creating a file in a visible folder"),
	)

	return container.NewBorder(controls, selected, nil, nil, treeSlot)
}
