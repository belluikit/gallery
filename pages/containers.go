package pages

import (
	"fmt"
	"image/color"
	"math"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// ContainersPage shows the containers the other pages do not: DocTabs,
// InnerWindow in MultipleWindows, HPaned and VPaned, Navigation,
// DirectionScope, Clip and PathClip, Opacity, RepaintBoundary and
// AdaptiveGrid (gallery#38).
func ContainersPage(_ bellui.Window) core.CanvasObject {
	box := func(w, h float32, o core.CanvasObject) core.CanvasObject {
		return container.NewGridWrap(geo.NewSize(w, h), o)
	}
	section := func(title, subtitle string, objs ...core.CanvasObject) *widget.Card {
		c := widget.NewCard(title, subtitle, container.NewVBox(objs...))
		c.Style = widget.CardStyleOutlined
		return c
	}

	// --- Doc tabs ---
	docStatus := widget.NewLabel("3 documents open")
	docNo := 3
	docs := container.NewDocTabs(
		container.NewTabItem("Draft 1", widget.NewLabel("The first draft.")),
		container.NewTabItem("Draft 2", widget.NewLabel("The second draft.")),
		container.NewTabItem("Draft 3", widget.NewLabel("The third draft.")),
	)
	docs.CreateTab = func() *container.TabItem {
		docNo++
		return container.NewTabItem(fmt.Sprintf("Draft %d", docNo), widget.NewLabel("A new draft."))
	}
	count := func(*container.TabItem) { docStatus.SetText(fmt.Sprintf("%d documents open", len(docs.Items))) }
	docs.OnClosed = count
	docs.OnSelected = count
	docCard := section("Doc Tabs", "Tabs you can close and add",
		note("Close a tab with its ×; the + makes a new one."), box(300, 120, docs), docStatus)

	// --- Inner windows ---
	win1 := container.NewInnerWindow("Notes", widget.NewLabel("Drag me by the title bar."))
	win1.Resize(geo.NewSize(200, 110))
	win2 := container.NewInnerWindow("Palette", widget.NewLabel("Resize me from the corner."))
	win2.Resize(geo.NewSize(200, 110))
	win2.Move(geo.NewPos(120, 110))
	windows := container.NewMultipleWindows(win1, win2)
	winCard := section("Inner Window", "Windows inside the window",
		note("Each one moves, resizes and comes to the front when tapped, inside the area it belongs to."),
		box(300, 260, windows))

	// --- Paned ---
	hpaned := container.NewHPaned(
		widget.NewCard("", "", widget.NewLabel("Left")),
		widget.NewCard("", "", widget.NewLabel("Middle")),
		widget.NewCard("", "", widget.NewLabel("Right")),
	)
	hpaned.Weights = []float64{1, 2, 1}
	vpaned := container.NewVPaned(
		widget.NewCard("", "", widget.NewLabel("Top")),
		widget.NewCard("", "", widget.NewLabel("Bottom")),
	)
	panedCard := section("Paned", "Panels with dividers you drag",
		note("Side by side, in proportion 1 : 2 : 1, and one above the other."),
		box(300, 110, hpaned), box(300, 180, vpaned))

	// --- Navigation ---
	var nav *container.Navigation
	list := container.NewVBox()
	for _, city := range []string{"Lisbon", "Kyoto", "Nairobi"} {
		list.Add(widget.NewButton(city, func() {
			nav.PushWithTitle(widget.NewLabel("All about "+city+"."), city)
		}))
	}
	nav = container.NewNavigationWithTitle(list, "Cities")
	navCard := section("Navigation", "Pages pushed and popped",
		note("Tap a city to push its page; the back arrow pops it."), box(300, 220, nav))

	// --- Direction scope ---
	form := func() core.CanvasObject {
		e := widget.NewEntry()
		e.SetText("Ada")
		return container.NewVBox(
			widget.NewForm(widget.NewFormItem("Name", e)),
			container.NewHBox(widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), nil), widget.NewCheck("Remember", nil)),
		)
	}
	rtl := container.NewDirectionScope(form(), bellui.DirectionRTL)
	dirCard := section("Direction Scope", "Right-to-left for one part",
		note("The same form twice: as the app lays it out, and inside a scope forced right-to-left."),
		form(), widget.NewSeparator(), rtl)

	// --- Clip and path clip ---
	gradient := func() core.CanvasObject {
		g := canvas.NewLinearGradient(color.NRGBA{R: 0x67, G: 0x50, B: 0xa4, A: 0xff}, color.NRGBA{R: 0xe0, G: 0x8a, B: 0x3c, A: 0xff}, 45)
		g.SetMinSize(geo.NewSize(120, 120))
		return g
	}
	// A clip only cuts what is drawn past its edge, so the content is a
	// larger gradient moved half out of it.
	big := gradient()
	big.Resize(geo.NewSize(200, 200))
	big.Move(geo.NewPos(-30, -30))
	clipped := container.NewClip(container.NewWithoutLayout(big))
	star := canvas.NewPath(func(p *canvas.Path) {
		const cx, cy, outer, inner = 60, 60, 58, 24
		for i := 0; i < 10; i++ {
			r := float64(outer)
			if i%2 == 1 {
				r = inner
			}
			a := math.Pi/2 + float64(i)*math.Pi/5
			x, y := cx+float32(r*math.Cos(a)), cy-float32(r*math.Sin(a))
			if i == 0 {
				p.MoveTo(x, y)
			} else {
				p.LineTo(x, y)
			}
		}
		p.Close()
	})
	starClip := container.NewPathClip(star, gradient())
	clipCard := section("Clip & Path Clip", "Cut content to a shape",
		note("A rectangle that cuts what spills past it, and a star-shaped path."),
		container.NewHBox(box(120, 120, clipped), box(120, 120, starClip)))

	// --- Opacity ---
	faded := container.NewOpacity(widget.NewCard("Faded card", "Its whole content fades as one", widget.NewButton("A button inside", nil)), 0.6)
	opacityLabel := widget.NewLabel("Opacity: 0.6")
	slider := widget.NewSlider(0, 1)
	slider.Step = 0.05
	slider.SetValue(0.6)
	slider.OnChanged = func(v float64) {
		faded.SetOpacity(float32(v))
		opacityLabel.SetText(fmt.Sprintf("Opacity: %.2f", v))
	}
	opacityCard := section("Opacity", "A group drawn see-through",
		note("The card and everything in it fade together, so overlapping parts don't show through each other."),
		faded, slider, opacityLabel)

	// --- Repaint boundary ---
	spinner := widget.NewProgressBarInfinite()
	boundary := container.NewRepaintBoundary(container.NewVBox(widget.NewLabel("Animated, and redrawn on its own:"), spinner))
	boundaryCard := section("Repaint Boundary", "Redraw one part on its own",
		note("It looks like nothing: the bar below is drawn from a cache, so its animation doesn't redraw the rest of the page."),
		boundary)

	// --- Adaptive grid ---
	var tiles []core.CanvasObject
	for i := 1; i <= 6; i++ {
		tiles = append(tiles, widget.NewCard("", "", widget.NewLabelWithStyle(fmt.Sprintf("Tile %d", i), bellui.TextAlignCenter, bellui.TextStyle{})))
	}
	adaptiveCard := section("Adaptive Grid", "Columns that turn into rows",
		note("Three columns on a desktop or a phone held sideways; three rows on a phone held upright."),
		container.NewAdaptiveGrid(3, tiles...))

	return container.NewVBox(docCard, winCard, panedCard, navCard, dirCard, clipCard, opacityCard, boundaryCard, adaptiveCard)
}
