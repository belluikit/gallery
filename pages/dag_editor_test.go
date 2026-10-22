package pages

import (
	"encoding/json"
	"testing"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/dag/dagwidget"
)

// dagPage is the DAG editor page in a test window, with handles on the parts
// the tests drive.
type dagPage struct {
	w      bellui.Window
	canvas *dagwidget.DAGCanvas
}

func newDAGPage(t *testing.T) *dagPage {
	t.Helper()
	test.NewTempApp(t)
	w := test.NewWindow(nil)
	t.Cleanup(w.Close)
	w.Resize(geo.NewSize(1200, 800))
	w.SetContent(DAGEditorPage(w))
	p := &dagPage{w: w}
	walkObjects(w.Content(), func(o core.CanvasObject) {
		if v, ok := o.(*dagwidget.DAGCanvas); ok {
			p.canvas = v
		}
	})
	if p.canvas == nil {
		t.Fatal("page has no DAG canvas")
	}
	return p
}

// panel is the properties panel on screen now; the page makes a new one for
// each selection.
func (p *dagPage) panel(t *testing.T) *dagwidget.PropertiesPanel {
	t.Helper()
	var panel *dagwidget.PropertiesPanel
	walkObjects(p.w.Content(), func(o core.CanvasObject) {
		if v, ok := o.(*dagwidget.PropertiesPanel); ok {
			panel = v
		}
	})
	if panel == nil {
		t.Fatal("page shows no properties panel")
	}
	return panel
}

// shows reports whether o draws the text s.
func shows(o core.CanvasObject, s string) bool {
	var found bool
	walkObjects(o, func(x core.CanvasObject) {
		switch v := x.(type) {
		case *canvas.Text:
			found = found || v.Text == s
		case *widget.Label:
			found = found || v.Text == s
		case *widget.RichText:
			found = found || v.String() == s
		}
	})
	return found
}

// walkObjects visits o and everything it renders, overlays aside. It does
// not descend into the DAG canvas, whose node widgets the tests reach through
// its own API.
func walkObjects(o core.CanvasObject, visit func(core.CanvasObject)) {
	if o == nil {
		return
	}
	visit(o)
	switch v := o.(type) {
	case *dagwidget.DAGCanvas:
		return
	case *bellui.Container:
		for _, child := range v.Objects {
			walkObjects(child, visit)
		}
	case *container.AppTabs:
		for _, item := range v.Items {
			walkObjects(item.Content, visit)
		}
	case core.Widget:
		for _, child := range test.WidgetRenderer(v).Objects() {
			walkObjects(child, visit)
		}
	}
}

func (p *dagPage) button(t *testing.T, text string) *widget.Button {
	t.Helper()
	var found *widget.Button
	walkObjects(p.w.Content(), func(o core.CanvasObject) {
		if b, ok := o.(*widget.Button); ok && b.Text == text {
			found = b
		}
	})
	if found == nil {
		t.Fatalf("no %q button", text)
	}
	return found
}

// field returns the input the panel draws for the form row labelled label.
// It reads what is on screen, not the form's item list: in x/dag v1.2.0 the
// two drift apart (the panel replaces the items without redrawing the form).
func (p *dagPage) field(t *testing.T, label string) core.CanvasObject {
	t.Helper()
	var found core.CanvasObject
	walkObjects(p.panel(t), func(o core.CanvasObject) {
		f, ok := o.(*widget.Form)
		if !ok || found != nil {
			return
		}
		row := -1
		for i, item := range f.Items {
			if item.Text == label {
				row = i
			}
		}
		if row < 0 {
			return
		}
		var inputs []core.CanvasObject
		for _, ro := range test.WidgetRenderer(f).Objects() {
			walkObjects(ro, func(x core.CanvasObject) {
				switch x.(type) {
				case *widget.Entry, *widget.Check, *widget.Select:
					inputs = append(inputs, x)
				}
			})
		}
		if row < len(inputs) {
			found = inputs[row]
		}
	})
	if found == nil {
		t.Fatalf("panel draws no %q field", label)
	}
	return found
}

// typeInto focuses e and types s the way the app's driver does: each key
// goes to whatever has the focus at that moment. It fails if a keystroke
// finds nothing focused, or a different widget from the one typed into.
func (p *dagPage) typeInto(t *testing.T, e *widget.Entry, s string) {
	t.Helper()
	c := p.w.Canvas()
	c.Focus(e)
	c.Focused().(event.Focusable).TypedKey(&event.KeyEvent{Name: event.KeyEnd}) // focusing selects all; type after it instead
	for i, r := range s {
		f := c.Focused()
		if f == nil {
			t.Fatalf("typing %q: nothing has the focus before key %d", s, i)
		}
		f.TypedRune(r)
	}
}

// panelHasFocus reports whether the focused widget is one the panel shows now,
// rather than one it has thrown away.
func (p *dagPage) panelHasFocus(t *testing.T) bool {
	f := p.w.Canvas().Focused()
	var in bool
	walkObjects(p.panel(t), func(o core.CanvasObject) {
		if fo, ok := f.(core.CanvasObject); ok && o == fo {
			in = true
		}
	})
	return in
}

func (p *dagPage) prop(t *testing.T, nodeID, key string) any {
	t.Helper()
	meta := map[string]any{}
	if err := json.Unmarshal(p.canvas.Graph.Nodes[nodeID].Metadata, &meta); err != nil {
		t.Fatal(err)
	}
	return meta[key]
}

func (p *dagPage) label(nodeID string) string {
	if n := p.canvas.Graph.Nodes[nodeID]; n != nil {
		return n.Label
	}
	return ""
}

// Every node in the sample pipeline has its component in the registry, or
// the panel shows it with no attributes and the palette cannot make one.
func TestDAGRegistryCoversPipeline(t *testing.T) {
	reg := dagRegistry()
	for id, n := range dagPipeline(reg).Nodes {
		comp := reg.Get(n.TypeTag)
		if comp == nil {
			t.Errorf("node %s: type %q is not in the registry", id, n.TypeTag)
			continue
		}
		if len(comp.Properties) == 0 {
			t.Errorf("component %q declares no attributes", comp.ID)
		}
	}
}

func TestDAGDeleteButton(t *testing.T) {
	p := newDAGPage(t)
	del := p.button(t, "Delete")
	if !del.Disabled() {
		t.Fatal("Delete is enabled with nothing selected")
	}

	edges := len(p.canvas.Graph.Edges)
	p.canvas.SelectNode("join")
	if del.Disabled() {
		t.Fatal("Delete is disabled with a node selected")
	}
	test.Tap(del)
	if p.canvas.Graph.Nodes["join"] != nil {
		t.Fatal("Delete left the selected node")
	}
	if got := len(p.canvas.Graph.Edges); got != edges-5 {
		t.Fatalf("after deleting join: %d edges, want %d", got, edges-5)
	}
	if !del.Disabled() {
		t.Fatal("Delete is enabled after the selection was deleted")
	}

	test.Tap(p.button(t, "Undo"))
	if p.canvas.Graph.Nodes["join"] == nil || len(p.canvas.Graph.Edges) != edges {
		t.Fatalf("Undo: join back %v, %d edges, want %d", p.canvas.Graph.Nodes["join"] != nil, len(p.canvas.Graph.Edges), edges)
	}
}

func TestDAGDeleteButtonReadOnly(t *testing.T) {
	p := newDAGPage(t)
	var editable *widget.Check
	walkObjects(p.w.Content(), func(o core.CanvasObject) {
		if c, ok := o.(*widget.Check); ok && c.Text == "Editable" {
			editable = c
		}
	})
	p.canvas.SelectNode("join")
	editable.SetChecked(false)
	del := p.button(t, "Delete")
	if !del.Disabled() {
		t.Fatal("Delete is enabled on a read-only canvas")
	}
	p.canvas.SelectNode("train")
	if !del.Disabled() {
		t.Fatal("selecting a node re-enabled Delete on a read-only canvas")
	}
	editable.SetChecked(true)
	if del.Disabled() {
		t.Fatal("Delete stayed disabled when editing came back on with a node selected")
	}
}

// The Delete key on a read-only canvas deletes nothing, but x/dag v1.2.0
// reports the selection as cleared anyway, so the page empties its panel
// while the node stays highlighted (belluikit/bellui#638). This pins
// that: when the fix lands, the second check fails and should be inverted.
func TestDAGDeleteKeyReadOnly(t *testing.T) {
	p := newDAGPage(t)
	p.canvas.Capabilities = dagwidget.Capabilities{}
	p.canvas.SelectNode("train")
	p.canvas.TypedKey(&event.KeyEvent{Name: event.KeyDelete})
	if p.canvas.Graph.Nodes["train"] == nil {
		t.Fatal("a read-only canvas deleted a node")
	}
	if shows(p.panel(t), "Train model") {
		t.Fatal("the panel still shows the node: x/dag no longer reports a clear it did not do; invert this check")
	}
}

// openNodeMenu does what a right-click or a long press on a node does: the
// canvas selects the node, then calls the page's menu hook.
func (p *dagPage) openNodeMenu(t *testing.T, id string) *widget.PopUpMenu {
	t.Helper()
	p.canvas.SelectNode(id)
	p.canvas.OnNodeSecondaryTapped(id, &event.PointEvent{AbsolutePosition: geo.NewPos(300, 300)})
	var menu *widget.PopUpMenu
	walkObjects(p.w.Canvas().Overlays().Top(), func(o core.CanvasObject) {
		if m, ok := o.(*widget.PopUpMenu); ok {
			menu = m
		}
	})
	if menu == nil {
		t.Fatal("no node menu in the top overlay")
	}
	return menu
}

// tapMenuItem taps the menu row showing label, as a finger would.
func tapMenuItem(t *testing.T, m *widget.PopUpMenu, label string) {
	t.Helper()
	for _, it := range m.Items {
		if tp, ok := it.(event.Tappable); ok && shows(it, label) {
			tp.Tapped(&event.PointEvent{})
			return
		}
	}
	t.Fatalf("menu has no %q item", label)
}

func TestDAGNodeMenuDelete(t *testing.T) {
	p := newDAGPage(t)
	m := p.openNodeMenu(t, "kpis")
	tapMenuItem(t, m, "Delete")
	if p.canvas.Graph.Nodes["kpis"] != nil {
		t.Fatal("menu Delete left the node")
	}
	test.Tap(p.button(t, "Undo"))
	if p.canvas.Graph.Nodes["kpis"] == nil {
		t.Fatal("Undo did not bring the node back")
	}
}

func TestDAGNodeMenuRename(t *testing.T) {
	p := newDAGPage(t)
	m := p.openNodeMenu(t, "kpis")
	tapMenuItem(t, m, "Rename…")

	// The dialog's entry has the old name and the focus.
	entry, ok := p.w.Canvas().Focused().(*widget.Entry)
	if !ok || entry.Text != "Aggregate" {
		t.Fatalf("focused %T %q, want the name entry holding %q", p.w.Canvas().Focused(), entry.Text, "Aggregate")
	}
	entry.SetText("")
	p.typeInto(t, entry, "KPIs")
	var rename *widget.Button
	walkObjects(p.w.Canvas().Overlays().Top(), func(o core.CanvasObject) {
		if b, ok := o.(*widget.Button); ok && b.Text == "Rename" {
			rename = b
		}
	})
	test.Tap(rename)

	if got := p.label("kpis"); got != "KPIs" {
		t.Fatalf("label %q after rename, want KPIs", got)
	}
	if got := p.field(t, "Name").(*widget.Entry).Text; got != "KPIs" {
		t.Fatalf("panel Name %q after rename, want KPIs", got)
	}
	test.Tap(p.button(t, "Undo"))
	if got := p.label("kpis"); got != "Aggregate" {
		t.Fatalf("label %q after Undo, want Aggregate", got)
	}
	if got := p.field(t, "Name").(*widget.Entry).Text; got != "Aggregate" {
		t.Fatalf("panel Name %q after Undo, want Aggregate", got)
	}
}

func TestDAGMenuGatedByEditable(t *testing.T) {
	p := newDAGPage(t)
	p.canvas.Capabilities = dagwidget.Capabilities{}
	tapMenuItem(t, p.openNodeMenu(t, "kpis"), "Delete")
	if p.canvas.Graph.Nodes["kpis"] == nil {
		t.Fatal("menu Delete removed a node on a read-only canvas")
	}
	tapMenuItem(t, p.openNodeMenu(t, "kpis"), "Rename…")
	if _, ok := p.w.Canvas().Overlays().Top().(*widget.PopUpMenu); !ok {
		if _, isEntry := p.w.Canvas().Focused().(*widget.Entry); isEntry {
			t.Fatal("menu Rename opened the dialog on a read-only canvas")
		}
	}
}

func TestDAGPanelRename(t *testing.T) {
	p := newDAGPage(t)
	p.canvas.SelectNode("dash")
	name := p.field(t, "Name").(*widget.Entry)
	if name.Text != "Dashboard" {
		t.Fatalf("Name field %q, want Dashboard", name.Text)
	}
	p.typeInto(t, name, " v2")
	if !p.panelHasFocus(t) {
		t.Errorf("after typing in Name the focus is on %T, not on the panel's field", p.w.Canvas().Focused())
	}
	if got := p.label("dash"); got != "Dashboard v2" {
		t.Fatalf("label %q after typing, want %q", got, "Dashboard v2")
	}
	for range len(" v2") {
		test.Tap(p.button(t, "Undo"))
	}
	if got := p.label("dash"); got != "Dashboard" {
		t.Fatalf("label %q after undoing each keystroke, want Dashboard", got)
	}
	if got := p.field(t, "Name").(*widget.Entry).Text; got != "Dashboard" {
		t.Fatalf("panel Name %q after Undo, want Dashboard", got)
	}
}

// Typing an attribute must keep every keystroke, and what the panel shows
// must be what the node holds.
func TestDAGPanelAttributes(t *testing.T) {
	p := newDAGPage(t)
	p.canvas.SelectNode("dash")
	sched := p.field(t, "Schedule").(*widget.Entry)
	if sched.Text != "daily at 06:00" {
		t.Fatalf("Schedule field %q, want the default", sched.Text)
	}
	sched.SetText("")
	p.typeInto(t, sched, "hourly")
	if !p.panelHasFocus(t) {
		t.Errorf("after typing in Schedule the focus is on %T, not on the panel's field", p.w.Canvas().Focused())
	}
	if got := p.prop(t, "dash", "schedule"); got != "hourly" {
		t.Fatalf("schedule %q after typing %q", got, "hourly")
	}
	if got := p.field(t, "Schedule").(*widget.Entry).Text; got != "hourly" {
		t.Fatalf("panel shows %q after typing %q", got, "hourly")
	}

	enabled := p.field(t, "Enabled").(*widget.Check)
	test.Tap(enabled)
	if got := p.prop(t, "dash", "enabled"); got != false {
		t.Fatalf("enabled %v after unticking", got)
	}
	test.Tap(p.button(t, "Undo"))
	if got := p.prop(t, "dash", "enabled"); got != true {
		t.Fatalf("enabled %v after Undo", got)
	}
}

// The library side of issue #16: tapping one palette item several times
// fans the nodes out instead of stacking them, and Undo removes each one.
func TestDAGPaletteTapFansOut(t *testing.T) {
	p := newDAGPage(t)
	var palette *dagwidget.Palette
	walkObjects(p.w.Content(), func(o core.CanvasObject) {
		if v, ok := o.(*dagwidget.Palette); ok {
			palette = v
		}
	})
	if palette == nil {
		t.Fatal("no palette")
	}
	var sinkItem event.Tappable
	walkObjects(palette, func(o core.CanvasObject) {
		if tp, ok := o.(event.Tappable); ok && sinkItem == nil && shows(o, "Sink") {
			sinkItem = tp
		}
	})
	if sinkItem == nil {
		t.Fatal("no Sink item in the palette")
	}

	before := len(p.canvas.Graph.Nodes)
	for range 3 {
		sinkItem.Tapped(&event.PointEvent{})
	}
	if got := len(p.canvas.Graph.Nodes); got != before+3 {
		t.Fatalf("%d nodes after three taps, want %d", got, before+3)
	}
	seen := map[[2]float32]string{}
	for id, n := range p.canvas.Graph.Nodes {
		if other, dup := seen[n.Position]; dup {
			t.Fatalf("nodes %s and %s share position %v", id, other, n.Position)
		}
		seen[n.Position] = id
	}
	for range 3 {
		test.Tap(p.button(t, "Undo"))
	}
	if got := len(p.canvas.Graph.Nodes); got != before {
		t.Fatalf("%d nodes after three undos, want %d", got, before)
	}
}

// The panel follows the selection: after one node, another node's fields are
// drawn, not the first one's.
func TestDAGPanelFollowsSelection(t *testing.T) {
	p := newDAGPage(t)
	p.canvas.SelectNode("dash")
	if got := p.field(t, "Schedule").(*widget.Entry).Text; got != "daily at 06:00" {
		t.Fatalf("Schedule %q", got)
	}
	p.canvas.SelectNode("sessions")
	if got := p.field(t, "Name").(*widget.Entry).Text; got != "Sessionize" {
		t.Fatalf("after selecting Sessionize the panel's Name field holds %q", got)
	}
	if got := p.field(t, "Batch size").(*widget.Entry).Text; got != "5000" {
		t.Fatalf("after selecting Sessionize the panel's Batch size holds %q, want 5000", got)
	}
	p.canvas.SelectGroup("cleaning")
	if got := p.field(t, "Label").(*widget.Entry).Text; got != "Cleaning" {
		t.Fatalf("after selecting the group the panel's Label holds %q", got)
	}
	p.canvas.ClearSelection()
	if !shows(p.panel(t), "Select a node to view properties") {
		t.Fatal("the panel is not empty after the selection was cleared")
	}
}

// On a phone the palette and the panel go below the canvas in tabs, the
// canvas keeps most of the screen, and selecting a node brings up its
// properties. Run with -tags mobile.
func TestDAGPhoneLayout(t *testing.T) {
	p := newDAGPage(t)
	if !bellui.CurrentDevice().IsMobile() {
		t.Skip("desktop layout; run with -tags mobile")
	}
	p.w.Resize(geo.NewSize(412, 860)) // a Pixel 9a in dp
	var tabs *container.AppTabs
	walkObjects(p.w.Content(), func(o core.CanvasObject) {
		if v, ok := o.(*container.AppTabs); ok {
			tabs = v
		}
	})
	if tabs == nil {
		t.Fatal("no palette/properties tabs on a phone")
	}
	cs := p.canvas.Size()
	if cs.Width < 380 || cs.Height < 300 {
		t.Fatalf("canvas is %v on a 412×860 screen", cs)
	}
	if min := p.w.Content().MinSize(); min.Width > 412 {
		t.Fatalf("page needs %v wide, more than the screen", min.Width)
	}
	tabs.SelectIndex(1) // Add node
	p.canvas.SelectNode("dash")
	if tabs.Selected().Text != "Properties" {
		t.Fatalf("selecting a node left the %q tab up", tabs.Selected().Text)
	}
}
