package pages

import (
	"encoding/json"
	"fmt"
	"image/color"
	"maps"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/dag/dag"
	"github.com/belluikit/bellui/x/dag/dagwidget"
)

// DAGEditorPage demos x/dag: a data pipeline on the interactive canvas, with
// styled nodes, a group and an anchored note, a component palette to add
// nodes, a properties panel to rename a node and edit its attributes, a node
// menu (right-click, or long press on a touch screen), undo/redo, auto-layout
// and a read-only switch. Everything works by touch alone.
func DAGEditorPage(w bellui.Window) core.CanvasObject {
	reg := dagRegistry()
	c := dagwidget.NewDAGCanvas(dagPipeline(reg))
	c.NodeStyler = dagStyler
	c.RefreshNodeStyles()
	c.Zoom = 0.5 // the whole pipeline fits beside the side tabs at the default window size
	c.Offset = geo.NewPos(5, 0)

	// The properties panel sits in a slot and is made anew for each thing it
	// shows: x/dag v1.2.0's panel does not redraw its fields when it is
	// pointed at another node, so a reused panel kept showing the first
	// node's fields (belluikit/bellui#638).
	panelSlot := container.NewStack()
	showInPanel := func(set func(*dagwidget.PropertiesPanel)) {
		p := dagwidget.NewPropertiesPanel(reg, c.Commands)
		p.CanvasRef = c // gives the panel its Name field, and its edits the canvas's undo stack
		set(p)
		panelSlot.Objects = []core.CanvasObject{p}
		panelSlot.Refresh()
	}
	showInPanel(func(p *dagwidget.PropertiesPanel) { p.SetNode(nil) })

	status := widget.NewLabel("Tap a node to edit it. Drag to pan, pinch or scroll to zoom.")
	deleteBtn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), c.DeleteSelected)
	deleteBtn.Disable() // nothing is selected yet

	// showPanel brings the properties tab to the front; it shares a tab strip
	// with the palette. Set once the tabs exist, below.
	var showPanel func()

	// selected redraws the panel for the current selection, after an undo, a
	// redo or a rename from the menu; nil when nothing is selected.
	var selected func()
	setSelection := func(show func(*dagwidget.PropertiesPanel)) {
		selected = func() { showInPanel(show) }
		selected()
		if c.Capabilities != (dagwidget.Capabilities{}) {
			deleteBtn.Enable()
		}
		showPanel()
	}
	c.OnNodeSelected = func(id string) {
		if n := c.Graph.Nodes[id]; n != nil {
			status.SetText(fmt.Sprintf("%s: %d in, %d out", n.Label, len(n.Inputs), len(n.Outputs)))
		}
		setSelection(func(p *dagwidget.PropertiesPanel) { p.SetNode(c.Graph.Nodes[id]) })
	}
	c.OnGroupSelected = func(id string) {
		setSelection(func(p *dagwidget.PropertiesPanel) { p.SetGroup(c.Graph.Groups[id]) })
	}
	c.OnAnnotationSelected = func(id string) {
		setSelection(func(p *dagwidget.PropertiesPanel) { p.SetAnnotation(c.Graph.Annotations[id]) })
	}
	c.OnEdgeSelected = func(string) {
		setSelection(func(p *dagwidget.PropertiesPanel) { p.SetNode(nil) })
	}
	c.OnSelectionCleared = func() {
		selected = nil
		showInPanel(func(p *dagwidget.PropertiesPanel) { p.SetNode(nil) })
		deleteBtn.Disable()
		status.SetText("Nothing selected.")
	}

	refreshPanel := func() {
		if selected != nil {
			selected()
		}
	}
	undo := func() { c.Commands.Undo(); refreshPanel() }
	redo := func() { c.Commands.Redo(); refreshPanel() }

	c.OnNodeSecondaryTapped = func(id string, ev *event.PointEvent) {
		rename := bellui.NewMenuItem("Rename…", func() { dagRenameDialog(w, c, id, refreshPanel) })
		rename.Disabled = !c.Capabilities.EditProperties
		del := bellui.NewMenuItemWithIcon("Delete", theme.DeleteIcon(), c.DeleteSelected) // the tap selected the node
		del.Disabled = !c.Capabilities.RemoveNodes
		widget.ShowPopUpMenuAtPosition(bellui.NewMenu("", rename, del), w.Canvas(), ev.AbsolutePosition)
	}

	editable := widget.NewCheck("Editable", func(on bool) {
		if on {
			c.Capabilities = dagwidget.AllCapabilities()
			if selected != nil {
				deleteBtn.Enable()
			}
		} else {
			c.Capabilities = dagwidget.Capabilities{}
			deleteBtn.Disable()
		}
		c.Refresh()
	})
	editable.SetChecked(true)

	tools := []core.CanvasObject{
		widget.NewButtonWithIcon("Undo", theme.ContentUndoIcon(), undo),
		widget.NewButtonWithIcon("Redo", theme.ContentRedoIcon(), redo),
		deleteBtn,
		widget.NewButton("Arrange", func() { _ = c.AutoLayout() }),
		widget.NewButton("Group", c.AddGroupAtCenter),
		widget.NewButton("Note", c.AddAnnotationAtCenter),
		editable,
	}
	palette := dagwidget.NewPalette(reg, c) // tap or drag a component onto the canvas

	// The panel and the palette share tabs, so neither squeezes the canvas:
	// beside it on desktop, below it on a phone.
	propsTab := container.NewTabItem("Properties", panelSlot)
	tabs := container.NewAppTabs(propsTab, container.NewTabItem("Add node", palette))
	showPanel = func() { tabs.Select(propsTab) }

	if !bellui.CurrentDevice().IsMobile() {
		// Wide enough for the panel's fields and a three-input port list.
		side := canvas.NewRectangle(color.Transparent)
		side.SetMinSize(geo.NewSize(320, 0))
		return container.NewBorder(container.NewHBox(tools...), status, nil, container.NewStack(side, tabs), c)
	}

	// A toolbar of seven controls is wider than a phone, so it wraps to rows
	// of four, and the tabs go below the canvas, which keeps most of the
	// height.
	split := container.NewVSplit(c, tabs)
	split.Offset = 0.6
	return container.NewBorder(container.NewGridWithColumns(4, tools...), status, nil, nil, split)
}

// dagRenameDialog asks for a node's new name and renames it as one undoable
// step.
func dagRenameDialog(w bellui.Window, c *dagwidget.DAGCanvas, id string, done func()) {
	n := c.Graph.Nodes[id]
	if n == nil {
		return
	}
	entry := widget.NewEntry()
	entry.SetText(n.Label)
	dialog.ShowForm("Rename node", "Rename", "Cancel",
		[]*widget.FormItem{widget.NewFormItem("Name", entry)},
		func(ok bool) {
			if !ok || entry.Text == "" || entry.Text == n.Label {
				return
			}
			c.Commands.Do(&dagwidget.RenameNodeCommand{Canvas: c, NodeID: id, Old: n.Label, New: entry.Text})
			done()
		}, w)
	w.Canvas().Focus(entry)
}

// dagRegistry lists the components the palette offers, with the attributes
// each one declares; the properties panel edits those.
func dagRegistry() *dag.Registry {
	reg := dag.NewRegistry()
	// One concrete stream type, so the palette's port summaries read
	// "rows → rows" rather than the wildcard's "* → *", while every component
	// still connects to every other.
	port := func(name string, dir dag.PortDirection) dag.PortDef {
		return dag.PortDef{Name: name, Direction: dir, Stream: dag.StreamType{Name: "rows"}}
	}
	in := func(names ...string) []dag.PortDef {
		ps := make([]dag.PortDef, len(names))
		for i, n := range names {
			ps[i] = port(n, dag.PortInput)
		}
		return ps
	}
	out := func(name string) []dag.PortDef { return []dag.PortDef{port(name, dag.PortOutput)} }
	for _, comp := range []*dag.Component{
		{
			ID: "source", Name: "Source", Category: "Sources", Outputs: out("out"),
			Properties: []dag.PropertyDef{
				{Key: "system", Label: "System", Kind: dag.PropertyEnum, EnumValues: []string{"postgres", "kafka", "s3"}, Default: "postgres"},
				{Key: "connection", Label: "Connection", Kind: dag.PropertyString, Required: true},
			},
		},
		{
			ID: "transform", Name: "Transform", Category: "Transforms", Inputs: in("in"), Outputs: out("out"),
			Properties: []dag.PropertyDef{
				{Key: "batch_size", Label: "Batch size", Kind: dag.PropertyInt, Default: 1000},
				{Key: "parallel", Label: "Parallel", Kind: dag.PropertyBool, Default: true},
			},
		},
		{
			ID: "join", Name: "Join", Category: "Transforms", Inputs: in("orders", "sessions", "customers"), Outputs: out("table"),
			Properties: []dag.PropertyDef{
				{Key: "key", Label: "Join key", Kind: dag.PropertyString, Default: "customer_id"},
				{Key: "how", Label: "How", Kind: dag.PropertyEnum, EnumValues: []string{"inner", "left", "outer"}, Default: "inner"},
			},
		},
		{
			ID: "ml", Name: "Train model", Category: "Models", Inputs: in("in"), Outputs: out("model"),
			Properties: []dag.PropertyDef{
				{Key: "algorithm", Label: "Algorithm", Kind: dag.PropertyEnum, EnumValues: []string{"boosted trees", "random forest", "logistic regression"}, Default: "boosted trees"},
				{Key: "learning_rate", Label: "Learning rate", Kind: dag.PropertyFloat, Default: 0.1},
			},
		},
		{
			ID: "sink", Name: "Sink", Category: "Sinks", Inputs: in("in"),
			Properties: []dag.PropertyDef{
				{Key: "schedule", Label: "Schedule", Kind: dag.PropertyString, Default: "daily at 06:00"},
				{Key: "enabled", Label: "Enabled", Kind: dag.PropertyBool, Default: true},
			},
		},
	} {
		reg.MustRegister(comp)
	}
	return reg
}

// dagPipeline is a nightly analytics pipeline laid out to fit the page. Its
// nodes are made from the registry's components, so their ports and
// attributes are the ones the properties panel describes.
func dagPipeline(reg *dag.Registry) *dag.Graph {
	g := dag.NewGraph("nightly", "Nightly analytics")
	node := func(id, label, comp string, x, y float32, props map[string]any) {
		n := reg.Get(comp).NewNode(id, [2]float32{x, y})
		n.Label = label
		if len(props) > 0 {
			meta := map[string]any{}
			_ = json.Unmarshal(n.Metadata, &meta)
			maps.Copy(meta, props)
			n.Metadata, _ = json.Marshal(meta)
		}
		_ = g.AddNode(n)
	}
	node("orders", "Orders DB", "source", 30, 80, map[string]any{"connection": "postgres://orders/prod"})
	node("clicks", "Clickstream", "source", 30, 260, map[string]any{"system": "kafka", "connection": "kafka://events/clicks"})
	node("crm", "CRM export", "source", 30, 440, map[string]any{"system": "s3", "connection": "s3://crm-exports/daily/"})
	node("clean", "Clean orders", "transform", 260, 80, nil)
	node("sessions", "Sessionize", "transform", 260, 260, map[string]any{"batch_size": 5000})
	node("dedupe", "Dedupe", "transform", 260, 440, map[string]any{"parallel": false})
	node("join", "Join", "join", 490, 245, nil)
	node("train", "Train model", "ml", 720, 110, nil)
	node("kpis", "Aggregate", "transform", 720, 390, nil)
	node("deploy", "Deploy model", "sink", 950, 110, map[string]any{"schedule": "after training"})
	node("dash", "Dashboard", "sink", 950, 390, nil)

	edge := func(from, fromPort, to, toPort string) {
		_ = g.AddEdge(&dag.Edge{
			ID:     from + "-" + to,
			Source: dag.PortID{NodeID: from, Name: fromPort},
			Target: dag.PortID{NodeID: to, Name: toPort},
		})
	}
	edge("orders", "out", "clean", "in")
	edge("clicks", "out", "sessions", "in")
	edge("crm", "out", "dedupe", "in")
	edge("clean", "out", "join", "orders")
	edge("sessions", "out", "join", "sessions")
	edge("dedupe", "out", "join", "customers")
	edge("join", "table", "train", "in")
	edge("join", "table", "kpis", "in")
	edge("train", "model", "deploy", "in")
	edge("kpis", "out", "dash", "in")

	_ = g.AddGroup(&dag.Group{
		ID: "cleaning", Label: "Cleaning", Members: []string{"clean", "sessions", "dedupe"},
		Position: [2]float32{240, 20}, Size: [2]float32{200, 520}, Color: "#4f8ef7",
	})
	_ = g.AddAnnotation(&dag.Annotation{
		ID:       "note",
		Text:     "Join waits for all three inputs. A late CRM file holds the whole run.",
		Position: [2]float32{490, 560}, Size: [2]float32{300, 50},
		AnchorKind: dag.AnchorNode, AnchorID: "join",
	})
	return g
}

var dagSubtitles = map[string]string{
	"orders": "postgres · 4.1 M rows", "clicks": "kafka · last 24 h", "crm": "s3 · daily CSV",
	"clean": "drop test orders", "sessions": "30 min idle gap", "dedupe": "customers by email",
	"join": "by customer ID", "train": "churn, boosted trees", "kpis": "revenue, churn, LTV",
	"deploy": "model registry", "dash": "daily at 06:00",
}

func dagStyler(n *dag.Node) dagwidget.NodeStyle {
	s := dagwidget.NodeStyle{Subtitle: dagSubtitles[n.ID]}
	switch n.TypeTag {
	case "source":
		s.Icon, s.Accent = theme.StorageIcon(), color.NRGBA{R: 79, G: 142, B: 247, A: 255}
	case "transform":
		s.Icon, s.Accent = theme.SettingsIcon(), color.NRGBA{R: 245, G: 166, B: 35, A: 255}
	case "join":
		s.Icon, s.Accent = theme.ListIcon(), color.NRGBA{R: 149, G: 117, B: 205, A: 255}
	case "ml":
		s.Icon, s.Accent = theme.ComputerIcon(), color.NRGBA{R: 236, G: 64, B: 122, A: 255}
	case "sink":
		s.Icon, s.Accent = theme.UploadIcon(), color.NRGBA{R: 38, G: 166, B: 91, A: 255}
	}
	return s
}
