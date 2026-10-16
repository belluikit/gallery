// Ported from the former fytables cmd/demo app.
package pages

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/tables"
	"github.com/belluikit/bellui/x/tables/pipeline"
	"github.com/belluikit/bellui/x/tables/sparkline"
)

// FytablesPage demonstrates the fytables DataGrid — an AG Grid equivalent for Bellui.
func FytablesPage(w bellui.Window) core.CanvasObject {
	tabs := container.NewAppTabs(
		&container.TabItem{Text: "Sort & Filter", Content: makeSortFilterTab(w)},
		&container.TabItem{Text: "Sparklines", Content: makeSparklineTab()},
		&container.TabItem{Text: "Columns", Content: makeColumnsTab()},
		&container.TabItem{Text: "Grouping", Content: makeGroupingTab()},
		&container.TabItem{Text: "Pivot", Content: makePivotTab()},
		&container.TabItem{Text: "Tree Data", Content: makeTreeTab()},
		&container.TabItem{Text: "Editing", Content: makeEditingTab()},
		&container.TabItem{Text: "Large Data", Content: makeLargeDataTab()},
	)

	return pageMinSize(1000, 560, tabs)
}

// ─── Tab 1: Sort & Filter ─────────────────────────────────────────────────

func makeSortFilterTab(w bellui.Window) core.CanvasObject {
	data := generateStocks(200)

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "symbol", HeaderName: "Symbol", Width: 80},
			{Field: "name", HeaderName: "Company", Width: 180},
			{Field: "sector", HeaderName: "Sector", Width: 120},
			{
				Field: "price", HeaderName: "Price", Type: tables.ColumnNumber, Width: 90,
				ValueFormatter: dollarFmt,
			},
			{
				Field: "change", HeaderName: "Change %", Type: tables.ColumnNumber, Width: 90,
				ValueFormatter: pctFmt,
			},
			{
				Field: "volume", HeaderName: "Volume", Type: tables.ColumnNumber, Width: 100,
				ValueFormatter: volumeFmt,
			},
		},
		func() int { return len(data) },
		func(row int) any { return data[row] },
	)
	grid.SelectionMode = tables.SelectionMultiRow
	grid.Options.RowAlternation = true
	grid.Options.ShowStatusBar = true

	// Quick filter input
	filterEntry := widget.NewEntry()
	filterEntry.PlaceHolder = "Quick filter (searches all columns)..."
	filterEntry.OnChanged = func(text string) {
		grid.API().SetQuickFilter(text)
	}

	// Export button
	exportBtn := widget.NewButton("Export CSV", func() {
		var buf bytes.Buffer
		if err := grid.API().ExportCSV(&buf); err != nil {
			dialog.ShowError(err, w)
			return
		}
		_ = os.WriteFile("export.csv", buf.Bytes(), 0o644)
		dialog.ShowInformation("Exported", fmt.Sprintf("Wrote %d bytes to export.csv", buf.Len()), w)
	})

	toolbar := container.NewBorder(nil, nil, nil, exportBtn, filterEntry)
	return container.NewBorder(toolbar, nil, nil, nil, grid)
}

// ─── Tab 2: Sparklines ────────────────────────────────────────────────────

func makeSparklineTab() core.CanvasObject {
	data := generateStocks(50)

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "symbol", HeaderName: "Symbol", Width: 80},
			{Field: "name", HeaderName: "Company", Width: 160},
			{
				Field: "trend", HeaderName: "Line Trend", Width: 140,
				CellRenderer: sparkline.NewCellRenderer(sparkline.Line, func(o *sparkline.Options) {
					o.Stroke = color.NRGBA{R: 33, G: 150, B: 243, A: 255}
					o.Fill = color.NRGBA{R: 33, G: 150, B: 243, A: 30}
					o.ShowMarkers = true
					o.MarkerSize = 2
				}),
			},
			{
				Field: "trend", HeaderName: "Area Trend", Width: 140,
				ColID: "area_trend",
				CellRenderer: sparkline.NewCellRenderer(sparkline.Area, func(o *sparkline.Options) {
					o.Stroke = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
					o.Fill = color.NRGBA{R: 76, G: 175, B: 80, A: 60}
				}),
			},
			{
				Field: "volumes", HeaderName: "Volume Bars", Width: 140,
				CellRenderer: sparkline.NewCellRenderer(sparkline.Column, func(o *sparkline.Options) {
					o.Stroke = color.NRGBA{R: 156, G: 39, B: 176, A: 255}
					o.NegativeColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255}
					o.ShowAxis = true
				}),
			},
			{
				Field: "changes", HeaderName: "Change Bars", Width: 140,
				CellRenderer: sparkline.NewCellRenderer(sparkline.Bar, func(o *sparkline.Options) {
					o.Stroke = color.NRGBA{R: 255, G: 152, B: 0, A: 255}
					o.NegativeColor = color.NRGBA{R: 244, G: 67, B: 54, A: 255}
					o.ShowAxis = true
				}),
			},
		},
		func() int { return len(data) },
		func(row int) any {
			m := data[row]
			// Add volume and change sparkline data if not present
			if _, ok := m["volumes"]; !ok {
				volumes := make([]float64, 7)
				changes := make([]float64, 7)
				for i := range volumes {
					volumes[i] = float64(rand.Intn(100))
					changes[i] = (rand.Float64() - 0.5) * 20
				}
				m["volumes"] = volumes
				m["changes"] = changes
			}
			return m
		},
	)
	grid.Options.RowAlternation = true

	label := widget.NewLabel("Hover over sparklines to see tooltips. All four chart types: Line, Area, Column, Bar.")
	return container.NewBorder(label, nil, nil, nil, grid)
}

// ─── Tab 3: Column Features ──────────────────────────────────────────────

func makeColumnsTab() core.CanvasObject {
	data := generateStocks(100)

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "symbol", HeaderName: "Symbol", Width: 80, Pinned: tables.PinLeft},
			{
				HeaderName: "Company Info",
				ColID:      "company_group",
				Children: []*tables.ColumnDef{
					{Field: "name", HeaderName: "Name", Width: 160},
					{Field: "sector", HeaderName: "Sector", Width: 120},
				},
			},
			{
				Field: "price", HeaderName: "Price", Type: tables.ColumnNumber, Width: 90,
				ValueFormatter: dollarFmt,
			},
			{
				Field: "change", HeaderName: "Change %", Type: tables.ColumnNumber, Width: 90,
				ValueFormatter: pctFmt,
			},
			{
				Field: "volume", HeaderName: "Volume", Type: tables.ColumnNumber, Width: 100,
				ValueFormatter: volumeFmt, Pinned: tables.PinRight,
			},
		},
		func() int { return len(data) },
		func(row int) any { return data[row] },
	)
	grid.Options.RowAlternation = true

	// Controls — all toggles
	collapseBtn := widget.NewButton("Collapse Group", nil)
	collapseBtn.OnTapped = func() {
		grid.ToggleColumnGroup("company_group")
		if grid.IsGroupCollapsed("company_group") {
			collapseBtn.SetText("Expand Group")
		} else {
			collapseBtn.SetText("Collapse Group")
		}
	}

	sectorBtn := widget.NewButton("Hide Sector", nil)
	sectorVisible := true
	sectorBtn.OnTapped = func() {
		sectorVisible = !sectorVisible
		grid.SetColumnVisible("sector", sectorVisible)
		if sectorVisible {
			sectorBtn.SetText("Hide Sector")
		} else {
			sectorBtn.SetText("Show Sector")
		}
	}

	symbolPinned := true
	pinBtn := widget.NewButton("Unpin Symbol", nil)
	pinBtn.OnTapped = func() {
		symbolPinned = !symbolPinned
		if symbolPinned {
			grid.API().SetColumnPinned("symbol", tables.PinLeft)
			pinBtn.SetText("Unpin Symbol")
		} else {
			grid.API().SetColumnPinned("symbol", tables.PinNone)
			pinBtn.SetText("Pin Symbol")
		}
	}

	fitBtn := widget.NewButton("Size to Fit", func() {
		grid.API().SizeColumnsToFit()
	})

	toolbar := container.NewHBox(collapseBtn, sectorBtn, pinBtn, fitBtn)
	info := widget.NewLabel("Symbol pinned left, Volume pinned right. Right-click headers for context menu. Drag header borders to resize.")
	top := container.NewVBox(info, toolbar)
	return container.NewBorder(top, nil, nil, nil, grid)
}

// ─── Tab 4: Row Grouping & Pagination ────────────────────────────────────

func makeGroupingTab() core.CanvasObject {
	data := generateStocks(500)

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "sector", HeaderName: "Sector", Width: 120},
			{Field: "symbol", HeaderName: "Symbol", Width: 80},
			{Field: "name", HeaderName: "Company", Width: 160},
			{
				Field: "price", HeaderName: "Price", Type: tables.ColumnNumber, Width: 90,
				ValueFormatter: dollarFmt,
			},
		},
		func() int { return len(data) },
		func(row int) any { return data[row] },
	)
	grid.Options.RowAlternation = true
	grid.Options.ShowStatusBar = true

	groupBtn := widget.NewButton("Group by Sector", func() {
		grid.API().SetRowGroupColumn("sector")
	})
	ungroupBtn := widget.NewButton("Clear Grouping", func() {
		grid.API().SetRowGroupColumn("")
	})
	expandBtn := widget.NewButton("Expand All", func() {
		grid.API().ExpandAllGroups()
	})
	collapseBtn := widget.NewButton("Collapse All", func() {
		grid.API().CollapseAllGroups()
	})
	paginateBtn := widget.NewButton("Paginate (50/page)", func() {
		grid.API().EnablePagination(50)
	})
	unpaginateBtn := widget.NewButton("Show All", func() {
		grid.API().DisablePagination()
	})
	nextPage := widget.NewButton("Next Page", func() {
		grid.API().SetPage(grid.CurrentPage() + 1)
	})
	prevPage := widget.NewButton("Prev Page", func() {
		p := grid.CurrentPage() - 1
		if p < 0 {
			p = 0
		}
		grid.API().SetPage(p)
	})

	toolbar := container.NewHBox(groupBtn, ungroupBtn, expandBtn, collapseBtn)
	pageBar := container.NewHBox(paginateBtn, unpaginateBtn, prevPage, nextPage)
	top := container.NewVBox(toolbar, pageBar)
	return container.NewBorder(top, nil, nil, nil, grid)
}

// ─── Tab 5: Pivot Table ──────────────────────────────────────────────────

func makePivotTab() core.CanvasObject {
	// Sales data for pivot
	regions := []string{"North", "South", "East", "West"}
	products := []string{"Widget A", "Widget B", "Widget C", "Gadget X", "Gadget Y"}
	var salesData []map[string]any
	for i := 0; i < 200; i++ {
		salesData = append(salesData, map[string]any{
			"region":  regions[rand.Intn(len(regions))],
			"product": products[rand.Intn(len(products))],
			"revenue": float64(rand.Intn(10000)) + rand.Float64()*100,
			"units":   rand.Intn(500),
		})
	}

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "region", HeaderName: "Region", Width: 100},
			{Field: "product", HeaderName: "Product", Width: 120},
			{
				Field: "revenue", HeaderName: "Revenue", Type: tables.ColumnNumber, Width: 100,
				ValueFormatter: dollarFmt,
			},
			{Field: "units", HeaderName: "Units", Type: tables.ColumnNumber, Width: 80},
		},
		func() int { return len(salesData) },
		func(row int) any { return salesData[row] },
	)
	grid.Options.RowAlternation = true

	pivotBtn := widget.NewButton("Pivot: Revenue by Region × Product", func() {
		grid.API().SetPivot("region", "product", "revenue", pipeline.AggSum)
		cols := grid.API().GetPivotColumns()
		fmt.Printf("Pivot columns: %v\n", cols)
	})
	clearBtn := widget.NewButton("Clear Pivot", func() {
		grid.API().SetPivot("", "", "", 0)
	})

	info := widget.NewLabel("Click 'Pivot' to aggregate revenue by region (rows) × product (columns). The grid dynamically creates columns for each product.")
	toolbar := container.NewHBox(pivotBtn, clearBtn)
	top := container.NewVBox(info, toolbar)
	return container.NewBorder(top, nil, nil, nil, grid)
}

// ─── Tab 6: Tree Data ────────────────────────────────────────────────────

func makeTreeTab() core.CanvasObject {
	// Org chart as tree data
	orgData := []map[string]any{
		{"id": 1, "parent": -1, "name": "CEO", "title": "Chief Executive Officer", "dept": "Executive"},
		{"id": 2, "parent": 1, "name": "CTO", "title": "Chief Technology Officer", "dept": "Engineering"},
		{"id": 3, "parent": 1, "name": "CFO", "title": "Chief Financial Officer", "dept": "Finance"},
		{"id": 4, "parent": 1, "name": "CMO", "title": "Chief Marketing Officer", "dept": "Marketing"},
		{"id": 5, "parent": 2, "name": "VP Eng", "title": "VP of Engineering", "dept": "Engineering"},
		{"id": 6, "parent": 2, "name": "VP Data", "title": "VP of Data Science", "dept": "Engineering"},
		{"id": 7, "parent": 3, "name": "Controller", "title": "Financial Controller", "dept": "Finance"},
		{"id": 8, "parent": 5, "name": "Alice", "title": "Senior Engineer", "dept": "Engineering"},
		{"id": 9, "parent": 5, "name": "Bob", "title": "Senior Engineer", "dept": "Engineering"},
		{"id": 10, "parent": 5, "name": "Charlie", "title": "Engineer", "dept": "Engineering"},
		{"id": 11, "parent": 6, "name": "Diana", "title": "Data Scientist", "dept": "Engineering"},
		{"id": 12, "parent": 6, "name": "Eve", "title": "ML Engineer", "dept": "Engineering"},
		{"id": 13, "parent": 4, "name": "Frank", "title": "Marketing Manager", "dept": "Marketing"},
		{"id": 14, "parent": 4, "name": "Grace", "title": "Content Lead", "dept": "Marketing"},
		{"id": 15, "parent": 7, "name": "Hank", "title": "Accountant", "dept": "Finance"},
	}

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "name", HeaderName: "Name", Width: 160},
			{Field: "title", HeaderName: "Title", Width: 200},
			{Field: "dept", HeaderName: "Department", Width: 120},
		},
		func() int { return len(orgData) },
		func(row int) any { return orgData[row] },
	)

	info := widget.NewLabel("Hierarchical org chart. Tree data uses parent-child IDs. Click expand/collapse icons.")
	return container.NewBorder(info, nil, nil, nil, grid)
}

// ─── Tab 7: Cell Editing ─────────────────────────────────────────────────

func makeEditingTab() core.CanvasObject {
	editable := true
	data := []map[string]any{
		{"name": "Alice", "role": "Engineer", "active": true, "salary": 120000, "started": "2020-03-15"},
		{"name": "Bob", "role": "Designer", "active": true, "salary": 95000, "started": "2021-06-01"},
		{"name": "Charlie", "role": "Manager", "active": false, "salary": 140000, "started": "2019-01-10"},
		{"name": "Diana", "role": "Engineer", "active": true, "salary": 130000, "started": "2022-09-20"},
		{"name": "Eve", "role": "Analyst", "active": true, "salary": 85000, "started": "2023-02-14"},
	}

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "name", HeaderName: "Name", Width: 120, Editable: &editable},
			{Field: "role", HeaderName: "Role", Width: 120, Editable: &editable},
			{Field: "active", HeaderName: "Active", Type: tables.ColumnBoolean, Width: 80, Editable: &editable},
			{
				Field: "salary", HeaderName: "Salary", Type: tables.ColumnNumber, Width: 100,
				Editable: &editable,
				ValueFormatter: func(v any) string {
					if n, ok := v.(int); ok {
						return fmt.Sprintf("$%d", n)
					}
					return fmt.Sprint(v)
				},
			},
			{Field: "started", HeaderName: "Start Date", Type: tables.ColumnDate, Width: 120, Editable: &editable},
		},
		func() int { return len(data) },
		func(row int) any { return data[row] },
	)
	grid.SelectionMode = tables.SelectionSingleCell
	grid.Options.RowAlternation = true
	grid.Options.CheckboxSelection = true

	grid.OnCellValueChanged = func(id tables.CellID, old, new any) {
		fmt.Printf("Cell [%d,%d] changed: %v → %v\n", id.Row, id.Col, old, new)
	}

	undoBtn := widget.NewButton("Undo (Ctrl+Z)", func() { grid.Undo() })
	redoBtn := widget.NewButton("Redo", func() { grid.Redo() })

	info := widget.NewLabel("Double-click or press F2 to edit. Enter to commit, Escape to cancel. Ctrl+C/V for clipboard. Undo/Redo supported.")
	toolbar := container.NewHBox(undoBtn, redoBtn)
	top := container.NewVBox(info, toolbar)
	return container.NewBorder(top, nil, nil, nil, grid)
}

// ─── Tab 8: Large Data Performance ────────────────────────────────────────

func makeLargeDataTab() core.CanvasObject {
	// 10,000 rows (5,000 on phones and tablets): plenty to show virtual
	// scrolling, at about 4 MB rather than the ~42 MB the 100,000-row sample
	// cost every time the page opened. Phones, and phone browsers running the
	// web demo, have far less memory to spare than a desktop.
	n := 10000
	if bellui.CurrentDevice().IsMobile() {
		n = 5000
	}
	start := time.Now()
	data := make([]map[string]any, n)
	for i := range data {
		data[i] = map[string]any{
			"id":     i + 1,
			"name":   fmt.Sprintf("Row %d", i+1),
			"value":  rand.Float64() * 1000,
			"group":  fmt.Sprintf("Group %c", 'A'+rune(i%26)),
			"active": i%3 != 0,
		}
	}
	genTime := time.Since(start)

	grid := tables.NewDataGrid(
		[]*tables.ColumnDef{
			{Field: "id", HeaderName: "ID", Type: tables.ColumnNumber, Width: 70},
			{Field: "name", HeaderName: "Name", Width: 120},
			{
				Field: "value", HeaderName: "Value", Type: tables.ColumnNumber, Width: 100,
				ValueFormatter: func(v any) string {
					if f, ok := v.(float64); ok {
						return fmt.Sprintf("%.2f", f)
					}
					return ""
				},
			},
			{Field: "group", HeaderName: "Group", Width: 100},
			{Field: "active", HeaderName: "Active", Type: tables.ColumnBoolean, Width: 70},
		},
		func() int { return len(data) },
		func(row int) any { return data[row] },
	)
	grid.SelectionMode = tables.SelectionMultiRow
	grid.Options.RowAlternation = true
	grid.Options.ShowStatusBar = true

	info := widget.NewLabel(fmt.Sprintf(
		"%d rows generated in %v. Virtual scrolling renders only visible cells. "+
			"Click headers to sort. Shift+click for multi-sort. Right-click for context menu.",
		n, genTime.Round(time.Millisecond),
	))

	paginateBtn := widget.NewButton("Paginate (1000/page)", func() {
		grid.API().EnablePagination(1000)
	})
	showAllBtn := widget.NewButton("Show All", func() {
		grid.API().DisablePagination()
	})
	filterBtn := widget.NewButton("Filter: Active Only", func() {
		grid.API().SetColumnFilter("active", func(v any) bool {
			b, _ := v.(bool)
			return b
		})
	})
	clearBtn := widget.NewButton("Clear Filter", func() {
		grid.API().ClearFilters()
	})

	toolbar := container.NewHBox(paginateBtn, showAllBtn, filterBtn, clearBtn)
	top := container.NewVBox(info, toolbar)
	return container.NewBorder(top, nil, nil, nil, grid)
}

// ─── Data Generators ─────────────────────────────────────────────────────

func generateStocks(n int) []map[string]any {
	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "AMD", "INTC"}
	names := []string{
		"Apple Inc.", "Alphabet Inc.", "Microsoft Corp.", "Amazon.com Inc.",
		"Tesla Inc.", "Meta Platforms", "NVIDIA Corp.", "Netflix Inc.", "AMD Inc.", "Intel Corp.",
	}
	sectors := []string{
		"Technology", "Technology", "Technology", "Consumer", "Automotive",
		"Technology", "Technology", "Entertainment", "Technology", "Technology",
	}

	data := make([]map[string]any, n)
	for i := range data {
		idx := i % len(symbols)
		data[i] = map[string]any{
			"symbol": symbols[idx],
			"name":   names[idx],
			"sector": sectors[idx],
			"price":  math.Round((100+rand.Float64()*400)*100) / 100,
			"change": math.Round((rand.Float64()-0.5)*10*100) / 100,
			"volume": rand.Intn(50000000),
			"trend":  generateTrend(14),
		}
	}
	return data
}

func generateTrend(points int) []float64 {
	trend := make([]float64, points)
	v := 100 + rand.Float64()*50
	for i := range trend {
		v += (rand.Float64() - 0.5) * 10
		trend[i] = math.Round(v*100) / 100
	}
	return trend
}

// ─── Value Formatters ────────────────────────────────────────────────────

func dollarFmt(v any) string {
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("$%.2f", n)
	case int:
		return fmt.Sprintf("$%d", n)
	}
	return ""
}

func pctFmt(v any) string {
	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%+.2f%%", f)
	}
	return ""
}

func volumeFmt(v any) string {
	if n, ok := v.(int); ok {
		if n >= 1000000 {
			return fmt.Sprintf("%.1fM", float64(n)/1e6)
		}
		if n >= 1000 {
			return fmt.Sprintf("%.1fK", float64(n)/1e3)
		}
		return fmt.Sprintf("%d", n)
	}
	return ""
}
