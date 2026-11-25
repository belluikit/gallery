package pages

import (
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/event"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// pressTile is a widget made outside bellui, from its public API alone, that
// shows the same MD3 pressed state as bellui's own Hyperlink and Chip: a
// ripple that spreads from the press and fades on release (gallery#49,
// belluikit/bellui#1449). Embedding widget.PressTarget gives it the touch and
// mouse press events; the rest is the three calls PressTarget asks for:
//
//   - AttachPress in CreateRenderer, with the ripple placed above the
//     background and below the content;
//   - TapPress in Tapped, so a tap with no press before it still shows one;
//   - CancelPress in HoverOut, since a mouse released elsewhere sends no
//     MouseUp.
type pressTile struct {
	widget.DisableableWidget
	widget.PressTarget

	text     string
	icon     bellui.Resource
	onTapped func()
}

var (
	_ event.Tappable   = (*pressTile)(nil)
	_ event.Hoverable  = (*pressTile)(nil)
	_ core.Disableable = (*pressTile)(nil)
	_ core.Accessible  = (*pressTile)(nil)
)

func newPressTile(text string, icon bellui.Resource, onTapped func()) *pressTile {
	t := &pressTile{text: text, icon: icon, onTapped: onTapped}
	t.ExtendBaseWidget(t)
	return t
}

// pressTileRadius rounds the tile, and the state layer with it.
const pressTileRadius = 12

func (t *pressTile) CreateRenderer() core.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &pressTileRenderer{
		tile:  t,
		bg:    canvas.NewRectangle(nil),
		label: canvas.NewText(t.text, nil),
		icon:  canvas.NewImageFromResource(nil),
	}
	r.bg.CornerRadius = pressTileRadius
	r.label.TextStyle = theme.TypographyLabelLarge
	// The shape is asked at each press. It covers the whole tile, in the
	// colour of the tile's text, and refuses the press while disabled.
	r.ripple = widget.AttachPress(&t.PressTarget, func() (widget.PressShape, bool) {
		return widget.PressShape{
			Size:         t.Size(),
			CornerRadius: pressTileRadius,
			Content:      r.color(r.contentColor()),
		}, !t.Disabled()
	})
	r.Refresh()
	return r
}

// Tapped runs the tile's action. TapPress shows the press for a tap that had
// none before it, a tap from code or an accessibility action; after a
// pointer press it does nothing.
func (t *pressTile) Tapped(ev *event.PointEvent) {
	if t.Disabled() {
		return
	}
	widget.TapPress(&t.PressTarget, ev)
	if t.onTapped != nil {
		t.onTapped()
	}
}

// Disable ends a press under way first: the tile stops being a target
// mid-press, and the release would otherwise still show.
func (t *pressTile) Disable() {
	widget.CancelPress(&t.PressTarget)
	t.DisableableWidget.Disable()
}

func (t *pressTile) HoverIn(*event.HoverEvent)    {}
func (t *pressTile) HoverMoved(*event.HoverEvent) {}

// HoverOut ends a press the pointer took off the tile.
func (t *pressTile) HoverOut() { widget.CancelPress(&t.PressTarget) }

func (t *pressTile) AccessibleRole() core.AccessibleRole { return core.RoleButton }
func (t *pressTile) AccessibleLabel() string             { return t.text }

type pressTileRenderer struct {
	tile   *pressTile
	bg     *canvas.Rectangle
	ripple *canvas.Ripple
	icon   *canvas.Image
	label  *canvas.Text
}

// The tile is a tonal surface: the secondary container, with its text and
// icon in the colour that goes on it, or the disabled colours.
func (r *pressTileRenderer) contentColor() bellui.ThemeColorName {
	if r.tile.Disabled() {
		return theme.ColorNameDisabled
	}
	return theme.ColorNameOnSecondaryContainer
}

func (r *pressTileRenderer) color(n bellui.ThemeColorName) color.Color {
	return r.tile.Theme().Color(n, bellui.CurrentApp().Settings().ThemeVariant())
}

func (r *pressTileRenderer) iconSize() float32 { return r.tile.Theme().Size(theme.SizeNameInlineIcon) }

func (r *pressTileRenderer) Layout(size geo.Size) {
	r.bg.Resize(size)
	r.ripple.Resize(size)
	pad := r.tile.Theme().Size(theme.SizeNamePadding) * 2
	is := r.iconSize()
	ls := r.label.MinSize()
	// Icon then label, centred together. Under right-to-left the icon
	// leads from the right.
	x := (size.Width - (is + pad/2 + ls.Width)) / 2
	iconX, labelX := x, x+is+pad/2
	if bellui.CurrentLayoutDirection() == bellui.DirectionRTL {
		iconX, labelX = size.Width-x-is, x
	}
	r.icon.Resize(geo.NewSquareSize(is))
	r.icon.Move(geo.NewPos(iconX, (size.Height-is)/2))
	r.label.Resize(ls)
	r.label.Move(geo.NewPos(labelX, (size.Height-ls.Height)/2))
}

func (r *pressTileRenderer) MinSize() geo.Size {
	pad := r.tile.Theme().Size(theme.SizeNamePadding) * 2
	ls := r.label.MinSize()
	// MD3's minimum touch target is 48 dp.
	return geo.NewSize(r.iconSize()+pad/2+ls.Width+pad*2, geo.Max(48, geo.Max(ls.Height, r.iconSize())+pad))
}

func (r *pressTileRenderer) Refresh() {
	if r.tile.Disabled() {
		r.bg.FillColor = r.color(theme.ColorNameDisabledButton)
		r.icon.Resource = theme.NewDisabledResource(r.tile.icon)
	} else {
		r.bg.FillColor = r.color(theme.ColorNameSecondaryContainer)
		r.icon.Resource = theme.NewColoredResource(r.tile.icon, theme.ColorNameOnSecondaryContainer)
	}
	r.label.Text = r.tile.text
	r.label.Color = r.color(r.contentColor())
	r.label.TextSize = r.tile.Theme().Size(theme.SizeNameText)
	r.bg.Refresh()
	r.icon.Refresh()
	r.label.Refresh()
	r.Layout(r.tile.Size())
}

func (r *pressTileRenderer) Objects() []core.CanvasObject {
	return []core.CanvasObject{r.bg, r.ripple, r.icon, r.label}
}

func (r *pressTileRenderer) Destroy() {}
