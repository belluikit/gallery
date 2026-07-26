package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// NavigationPage shows Navigation Bar, Rail, Drawer, Tabs, Top/Bottom App Bar.
func NavigationPage(w bellui.Window) core.CanvasObject {
	// --- Navigation Bar ---
	navBar := widget.NewNavigationBar([]*widget.NavigationBarItem{
		{Label: "Home", Icon: theme.HomeIcon()},
		{Label: "Search", Icon: theme.SearchIcon()},
		{Label: "Profile", Icon: theme.AccountIcon()},
	}, func(int) {})

	navBarCard := widget.NewCard("Navigation Bar", "Bottom navigation (3-5 items)", navBar)
	navBarCard.Style = widget.CardStyleOutlined

	// --- Navigation Rail ---
	navRail := widget.NewNavigationRail([]*widget.NavigationRailItem{
		{Label: "Home", Icon: theme.HomeIcon()},
		{Label: "Search", Icon: theme.SearchIcon()},
		{Label: "Settings", Icon: theme.SettingsIcon()},
	}, func(int) {})

	navRailCard := widget.NewCard("Navigation Rail", "Side navigation", container.NewHBox(navRail, widget.NewLabel("Content area")))
	navRailCard.Style = widget.CardStyleOutlined

	// --- Top App Bar ---
	topBar := widget.NewTopAppBar("Page Title", theme.NavigateBackIcon(), func() {})
	topBar.Actions = []widget.TopAppBarAction{
		{Icon: theme.SearchIcon(), OnTapped: func() {}},
		{Icon: theme.SettingsIcon(), OnTapped: func() {}},
	}

	topBarCard := widget.NewCard("Top App Bar", "Header with title and actions", topBar)
	topBarCard.Style = widget.CardStyleOutlined

	// --- Bottom App Bar ---
	fab := widget.NewFAB(theme.ContentAddIcon(), func() {})
	bottomBar := widget.NewBottomAppBar([]*widget.BottomAppBarAction{
		{Icon: theme.MailComposeIcon(), OnTapped: func() {}},
		{Icon: theme.SearchIcon(), OnTapped: func() {}},
	}, fab)

	bottomBarCard := widget.NewCard("Bottom App Bar", "Toolbar with FAB", bottomBar)
	bottomBarCard.Style = widget.CardStyleOutlined

	// --- Tabs ---
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Content 1")),
		container.NewTabItem("Tab 2", widget.NewLabel("Content 2")),
		container.NewTabItem("Tab 3", widget.NewLabel("Content 3")),
	)

	tabsCard := widget.NewCard("Tabs", "Tabbed content", tabs)
	tabsCard.Style = widget.CardStyleOutlined

	// --- Navigation Drawer ---
	drawerCard := widget.NewCard("Navigation Drawer", "Slide-out side panel", container.NewVBox(
		widget.NewButton("Show Drawer", func() {
			drawer := widget.NewNavigationDrawer([]*widget.NavigationDrawerItem{
				{Label: "Home", Icon: theme.HomeIcon()},
				{Label: "Settings", Icon: theme.SettingsIcon()},
				{Label: "", IsDivider: true},
				{Label: "Help", Icon: theme.HelpIcon()},
			}, func(int) {})
			drawer.Show(w.Canvas())
		}),
	))
	drawerCard.Style = widget.CardStyleOutlined

	// --- Toolbar ---
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentUndoIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentRedoIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)
	toolbarCard := widget.NewCard("Toolbar", "Action icons with separators", toolbar)
	toolbarCard.Style = widget.CardStyleOutlined

	return container.NewVBox(navBarCard, navRailCard, topBarCard, bottomBarCard, tabsCard, drawerCard, toolbarCard)
}
