package main

import (
	"net/url"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/app"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dev"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
	"github.com/belluikit/gallery/pages"
)

// corePages are pure-Go destinations that build and run on every platform,
// including Android. Native-only pages (FFmpeg/FAISS) are appended from
// pages.NativePages(), which is empty on mobile builds.
func corePages() []pages.PageDef {
	return []pages.PageDef{
		{Name: "Actions", Icon: theme.ConfirmIcon(), New: pages.ActionsPage},
		{Name: "Animated GIF", Icon: theme.MediaPlayIcon(), New: pages.AnimatedGifPage},
		{Name: "Auth", Icon: theme.AccountIcon(), New: pages.AuthPage},
		{Name: "Charts", Icon: theme.HistoryIcon(), New: pages.FychartsPage},
		{Name: "Code Editor", Icon: theme.DocumentCreateIcon(), New: pages.FycoderPage, Fill: true},
		{Name: "Communication", Icon: theme.InfoIcon(), New: pages.CommunicationPage},
		{Name: "Containment", Icon: theme.ContentCopyIcon(), New: pages.ContainmentPage},
		{Name: "DAG Editor", Icon: theme.ViewFullScreenIcon(), New: pages.DAGEditorPage, Fill: true},
		{Name: "Data Grid", Icon: theme.GridIcon(), New: pages.FytablesPage},
		{Name: "File Transfer", Icon: theme.UploadIcon(), New: pages.UploaderPage},
		{Name: "File Tree", Icon: theme.FolderOpenIcon(), New: pages.FileTreePage, Fill: true},
		{Name: "Font Picker", Icon: theme.DocumentIcon(), New: pages.FontPickerPage},
		{Name: "Inputs", Icon: theme.DocumentCreateIcon(), New: pages.InputsPage},
		{Name: "Navigation", Icon: theme.ComputerIcon(), New: pages.NavigationPage},
		{Name: "Responsive", Icon: theme.ViewFullScreenIcon(), New: pages.ResponsivePage},
		{Name: "Scrolling", Icon: theme.MoveDownIcon(), New: pages.ScrollingPage, Fill: true},
		{Name: "Theme", Icon: theme.ColorPaletteIcon(), New: pages.ThemePage},
		{Name: "Typography", Icon: theme.DocumentIcon(), New: pages.TypographyPage},
	}
}

func main() {
	a := app.NewWithID("io.bellui.md3gallery")

	// Restore the appearance chosen on the Theme page (stored in the
	// gallery's own Preferences; applies to this app only).
	if stored, ok := widget.LoadThemeSelection(a.Preferences()); ok {
		stored.Apply()
	}

	w := a.NewWindow("Material Design 3 Gallery")

	// On mobile the window is always fullscreen, so the desktop preferred
	// size and the on-canvas dev inspector only make sense on desktop.
	isMobile := bellui.CurrentDevice().IsMobile()
	if !isMobile {
		w.Resize(geo.NewSize(900, 650))
		dev.Enable(w)
	}

	// Page registry: core pages everywhere, native-only pages appended on
	// desktop. The drawer and content are both driven off this single slice.
	registry := append(corePages(), pages.NativePages()...)

	// Content area. pageSlot hosts the active page; how it is hosted depends on
	// the page:
	//   - Fill pages (e.g. the code editor) are placed directly so they manage
	//     their own scrolling and pointer handling. Wrapping such a widget in an
	//     outer container.Scroll lets the scroll intercept presses/drags, which
	//     breaks focus, gutter taps, and modifier-clicks.
	//   - Other pages get the gallery's outer scroll. Phones can't horizontally
	//     reflow the wide component pages, so both axes scroll there; desktop keeps
	//     vertical-only. Paired with pages.pageMinSize, which drops its forced
	//     desktop width on mobile.
	current := 0
	pageSlot := container.NewStack()
	setPage := func() {
		pages.LeavePage() // stop whatever the previous page left running
		page := registry[current].New(w)
		if registry[current].Fill {
			pageSlot.Objects = []core.CanvasObject{page}
		} else {
			inner := container.NewPadded(page)
			if isMobile {
				pageSlot.Objects = []core.CanvasObject{container.NewScroll(inner)}
			} else {
				pageSlot.Objects = []core.CanvasObject{container.NewVScroll(inner)}
			}
		}
		pageSlot.Refresh()
	}
	setPage()

	rebuild := func() { setPage() }

	// onUI runs fn on Bellui's render goroutine. The mobile driver can dispatch
	// some tap callbacks from a non-render goroutine (e.g. double-tap
	// detection), so mutating widgets/overlays from a handler is marshalled
	// onto the render loop via bellui.Do. On desktop this is an effectively free
	// hop. Paired with the manifest.pxf belluiDo migration flag.
	onUI := func(fn func()) { bellui.Do(fn) }

	// Navigation drawer built from the registry.
	drawerItems := make([]*widget.NavigationDrawerItem, len(registry))
	for i, p := range registry {
		drawerItems[i] = &widget.NavigationDrawerItem{Label: p.Name, Icon: p.Icon}
	}

	var appBar *widget.TopAppBar
	drawer := widget.NewNavigationDrawer(drawerItems, func(idx int) {
		onUI(func() {
			current = idx
			rebuild()
			appBar.Title = registry[idx].Name
			appBar.Refresh()
		})
	})

	// TopAppBar with hamburger menu and theme actions.
	appBar = widget.NewTopAppBar(registry[current].Name, theme.MenuIcon(), func() {
		onUI(func() { drawer.Toggle(w.Canvas()) })
	})

	isDense := false
	isRTL := false
	// Chromatic in dark mode ("switch to light"), achromatic in light —
	// seeded from the current variant since a stored Theme-page selection
	// may already have applied light mode.
	variantIcon := theme.ColorChromaticIcon()
	if a.Settings().ThemeVariant() == theme.VariantLight {
		variantIcon = theme.ColorAchromaticIcon()
	}
	appBar.Actions = []widget.TopAppBarAction{
		{
			Icon: theme.AccountIcon(),
			OnTapped: func() {
				onUI(func() { pages.ShowAuthDialog(w) })
			},
		},
		{
			Icon: variantIcon,
			OnTapped: func() {
				onUI(func() {
					if a.Settings().ThemeVariant() == theme.VariantDark {
						a.Settings().SetThemeVariant(theme.VariantLight)
						appBar.Actions[1].Icon = theme.ColorAchromaticIcon()
					} else {
						a.Settings().SetThemeVariant(theme.VariantDark)
						appBar.Actions[1].Icon = theme.ColorChromaticIcon()
					}
					appBar.Refresh()
					rebuild()
				})
			},
		},
		{
			Icon: theme.GridIcon(),
			OnTapped: func() {
				onUI(func() {
					isDense = !isDense
					if isDense {
						a.Settings().SetTheme(theme.DenseTheme())
					} else {
						a.Settings().SetTheme(theme.DefaultTheme())
					}
					rebuild()
				})
			},
		},
		{
			// Toggle the process-wide layout direction to preview the whole
			// gallery right-to-left. SetLayoutDirection re-lays-out and repaints
			// the running window (mirroring containers, alignment, scroll, focus
			// order and directional icons); rebuild() refreshes the current page
			// to match. MoreHorizontal (LTR) / MoreVertical (RTL) flag the state.
			Icon: theme.MoreHorizontalIcon(),
			OnTapped: func() {
				onUI(func() {
					isRTL = !isRTL
					if isRTL {
						bellui.SetLayoutDirection(bellui.DirectionRTL)
						appBar.Actions[3].Icon = theme.MoreVerticalIcon()
					} else {
						bellui.SetLayoutDirection(bellui.DirectionLTR)
						appBar.Actions[3].Icon = theme.MoreHorizontalIcon()
					}
					appBar.Refresh()
					rebuild()
				})
			},
		},
		{
			// dialog.About (fyne #446/#453): MD3-native about box built from the
			// app metadata, with the parallax icon underlay behind the content.
			Icon: theme.InfoIcon(),
			OnTapped: func() {
				onUI(func() { showAbout(a, w) })
			},
		},
	}

	body := container.NewBorder(appBar, nil, nil, nil, pageSlot)

	w.SetContent(body)
	bellui.ShowAndRun(w)
}

// showAbout opens the MD3 about dialog for the gallery itself: metadata
// header (icon, name, version), markdown body over the parallax underlay,
// and a links footer.
func showAbout(a bellui.App, w bellui.Window) {
	gallerySite, _ := url.Parse("https://bellui.com/gallery")
	belluiSite, _ := url.Parse("https://bellui.com")
	links := []*widget.Hyperlink{
		widget.NewHyperlink("Gallery", gallerySite),
		widget.NewHyperlink("Bellui", belluiSite),
	}
	dialog.ShowAbout(`## Material Design 3 Gallery

A showcase for the **Bellui** toolkit: MD3 components, the embedded
code editor, charts, data grids, and the widgets added in each release.

Scroll this text to see the parallax underlay drift behind it.

---

Built against the Bellui toolkit — every page renders through the
same code you are changing.`, links, a, w)
}
