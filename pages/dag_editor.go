package pages

import (
	"fmt"
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/dag/dag"
	"github.com/belluikit/bellui/x/dag/dagwidget"
)

// DAGEditorPage demos x/dag: a data pipeline on the interactive canvas, with
// styled nodes, a group and an anchored note, a component palette to add
// nodes, undo/redo, auto-layout and a read-only switch.
func DAGEditorPage(_ bellui.Window) core.CanvasObject {
	reg := dagRegistry()
	c := dagwidget.NewDAGCanvas(dagPipeline())
	c.NodeStyler = dagStyler
	c.RefreshNodeStyles()
	c.Zoom = 0.55 // the whole pipeline fits beside the palette at the default window size
	c.Offset = geo.NewPos(5, 0)

	status := widget.NewLabel("Click a node, drag to pan, scroll to zoom.")
	c.OnNodeSelected = func(id string) {
		if n := c.Graph.Nodes[id]; n != nil {
			status.SetText(fmt.Sprintf("%s: %d in, %d out", n.Label, len(n.Inputs), len(n.Outputs)))
		}
	}
	c.OnSelectionCleared = func() { status.SetText("Nothing selected.") }

	editable := widget.NewCheck("Editable", func(on bool) {
		if on {
			c.Capabilities = dagwidget.AllCapabilities()
		} else {
			c.Capabilities = dagwidget.Capabilities{}
		}
		c.Refresh()
	})
	editable.SetChecked(true)

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("Undo", theme.ContentUndoIcon(), func() { c.Commands.Undo() }),
		widget.NewButtonWithIcon("Redo", theme.ContentRedoIcon(), func() { c.Commands.Redo() }),
		widget.NewButton("Arrange", func() { _ = c.AutoLayout() }),
		widget.NewButton("Group", c.AddGroupAtCenter),
		widget.NewButton("Note", c.AddAnnotationAtCenter),
		editable,
	)
	palette := dagwidget.NewPalette(reg, c) // click or drag a component onto the canvas

	return container.NewBorder(toolbar, status, palette, nil, c)
}

// dagRegistry lists the components the palette offers.
func dagRegistry() *dag.Registry {
	reg := dag.NewRegistry()
	// One concrete stream type, so the palette's port summaries read
	// "rows → rows" rather than the wildcard's "* → *", while every component
	// still connects to every other.
	port := func(name string, dir dag.PortDirection) dag.PortDef {
		return dag.PortDef{Name: name, Direction: dir, Stream: dag.StreamType{Name: "rows"}}
	}
	for _, comp := range []*dag.Component{
		{ID: "source", Name: "Source", Category: "Sources", Outputs: []dag.PortDef{port("out", dag.PortOutput)}},
		{ID: "transform", Name: "Transform", Category: "Transforms",
			Inputs: []dag.PortDef{port("in", dag.PortInput)}, Outputs: []dag.PortDef{port("out", dag.PortOutput)}},
		{ID: "sink", Name: "Sink", Category: "Sinks", Inputs: []dag.PortDef{port("in", dag.PortInput)}},
	} {
		reg.MustRegister(comp)
	}
	return reg
}

func dagPorts(dir dag.PortDirection, names ...string) []dag.Port {
	ps := make([]dag.Port, len(names))
	for i, n := range names {
		ps[i] = dag.Port{Name: n, Direction: dir}
	}
	return ps
}

// dagPipeline is a nightly analytics pipeline laid out to fit the page.
func dagPipeline() *dag.Graph {
	in := func(names ...string) []dag.Port { return dagPorts(dag.PortInput, names...) }
	out := func(names ...string) []dag.Port { return dagPorts(dag.PortOutput, names...) }

	g := dag.NewGraph("nightly", "Nightly analytics")
	for _, n := range []*dag.Node{
		{ID: "orders", Label: "Orders DB", TypeTag: "source", Outputs: out("rows"), Position: [2]float32{30, 80}},
		{ID: "clicks", Label: "Clickstream", TypeTag: "source", Outputs: out("events"), Position: [2]float32{30, 260}},
		{ID: "crm", Label: "CRM export", TypeTag: "source", Outputs: out("file"), Position: [2]float32{30, 440}},
		{ID: "clean", Label: "Clean orders", TypeTag: "transform", Inputs: in("in"), Outputs: out("out"), Position: [2]float32{260, 80}},
		{ID: "sessions", Label: "Sessionize", TypeTag: "transform", Inputs: in("in"), Outputs: out("out"), Position: [2]float32{260, 260}},
		{ID: "dedupe", Label: "Dedupe", TypeTag: "transform", Inputs: in("in"), Outputs: out("out"), Position: [2]float32{260, 440}},
		{ID: "join", Label: "Join", TypeTag: "join", Inputs: in("orders", "sessions", "customers"), Outputs: out("table"), Position: [2]float32{490, 245}},
		{ID: "train", Label: "Train model", TypeTag: "ml", Inputs: in("in"), Outputs: out("model"), Position: [2]float32{720, 110}},
		{ID: "kpis", Label: "Aggregate", TypeTag: "transform", Inputs: in("in"), Outputs: out("out"), Position: [2]float32{720, 390}},
		{ID: "deploy", Label: "Deploy model", TypeTag: "sink", Inputs: in("in"), Position: [2]float32{950, 110}},
		{ID: "dash", Label: "Dashboard", TypeTag: "sink", Inputs: in("in"), Position: [2]float32{950, 390}},
	} {
		_ = g.AddNode(n)
	}
	edge := func(from, fromPort, to, toPort string) {
		_ = g.AddEdge(&dag.Edge{
			ID:     from + "-" + to,
			Source: dag.PortID{NodeID: from, Name: fromPort},
			Target: dag.PortID{NodeID: to, Name: toPort},
		})
	}
	edge("orders", "rows", "clean", "in")
	edge("clicks", "events", "sessions", "in")
	edge("crm", "file", "dedupe", "in")
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
