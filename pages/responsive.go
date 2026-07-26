package pages

import (
	"fmt"
	"image/color"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/layout"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// ResponsivePage demos layout.Responsive (fyne #444): children declare a
// width ratio per MD3 window size class and rows reflow as the window
// resizes. Resize the gallery window to watch the grid reorganize.
func ResponsivePage(_ bellui.Window) core.CanvasObject {
	intro := widget.NewLabel("Resize the window: each tile declares a width ratio per MD3 size class " +
		"(compact <600, medium <840, expanded <1200, large <1600, extra-large). " +
		"Full width when compact, halves when medium, quarters when expanded and beyond.")
	intro.Wrapping = bellui.TextWrapWord

	classLabel := widget.NewLabelWithStyle("", bellui.TextAlignCenter, theme.TypographyTitleMedium)

	tile := func(title string, c color.Color) core.CanvasObject {
		bg := canvas.NewRectangle(c)
		bg.CornerRadius = theme.Size(theme.SizeNameShapeMedium)
		bg.SetMinSize(geo.NewSize(0, 90))
		label := widget.NewLabelWithStyle(title, bellui.TextAlignCenter, bellui.TextStyle{Bold: true})
		return container.NewStack(bg, container.NewCenter(label))
	}

	grid := layout.NewResponsiveLayout(
		layout.Responsive(tile("1 / ½ / ¼", color.NRGBA{R: 0x5c, G: 0x7c, B: 0xa8, A: 0xff}), 1, .5, .25),
		layout.Responsive(tile("1 / ½ / ¼", color.NRGBA{R: 0x6b, G: 0x9a, B: 0x7c, A: 0xff}), 1, .5, .25),
		layout.Responsive(tile("1 / ½ / ¼", color.NRGBA{R: 0xa8, G: 0x8a, B: 0x5c, A: 0xff}), 1, .5, .25),
		layout.Responsive(tile("1 / ½ / ¼", color.NRGBA{R: 0x9a, G: 0x6b, B: 0x8a, A: 0xff}), 1, .5, .25),
		layout.Responsive(tile("always ⅓", color.NRGBA{R: 0x7c, G: 0x6b, B: 0x9a, A: 0xff}), 1.0/3, 1.0/3, 1.0/3),
		layout.Responsive(tile("always ⅓", color.NRGBA{R: 0x5c, G: 0x8a, B: 0x8a, A: 0xff}), 1.0/3, 1.0/3, 1.0/3),
		layout.Responsive(tile("always ⅓", color.NRGBA{R: 0x8a, G: 0x7c, B: 0x6b, A: 0xff}), 1.0/3, 1.0/3, 1.0/3),
		tile("plain object: full row in every class", color.NRGBA{R: 0x60, G: 0x60, B: 0x68, A: 0xff}),
	)

	// classProbe reports the responsive grid's own width and size class; its
	// Layout hook runs whenever the window (and therefore the grid) resizes.
	probe := &classProbe{target: grid, label: classLabel}
	probeHost := container.New(probe, grid)

	return container.NewBorder(container.NewVBox(intro, classLabel), nil, nil, nil,
		container.NewVScroll(probeHost))
}

// classProbe is a pass-through layout that mirrors the hosted object to the
// full size and updates a label with the active MD3 window size class.
type classProbe struct {
	target core.CanvasObject
	label  *widget.Label
}

func (p *classProbe) Layout(objects []core.CanvasObject, size geo.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(geo.Position{})
	}
	class := "extra-large"
	for i, bound := range layout.DefaultResponsiveBreakpoints {
		if size.Width < bound {
			class = []string{"compact", "medium", "expanded", "large"}[i]
			break
		}
	}
	text := fmt.Sprintf("width %.0f dp — %s", size.Width, class)
	if p.label.Text != text {
		p.label.SetText(text)
	}
}

func (p *classProbe) MinSize(objects []core.CanvasObject) geo.Size {
	if len(objects) == 0 {
		return geo.Size{}
	}
	return objects[0].MinSize()
}
