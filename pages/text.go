package pages

import (
	"net/url"
	"strings"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// TextPage shows the text widgets the other pages do not: SelectEntry,
// SearchSelect, CompletionPopup, RichEntry, RichText built from HTML and from
// segments, and TextGrid (gallery#38).
func TextPage(w bellui.Window) core.CanvasObject {
	var disableables []core.Disableable

	// --- Select entry ---
	selEntryStatus := widget.NewLabel("Fruit: none")
	selEntry := widget.NewSelectEntry([]string{"Apple", "Banana", "Cherry", "Mango", "Orange"})
	selEntry.PlaceHolder = "Type a fruit, or pick one"
	selEntry.OnChanged = func(s string) { selEntryStatus.SetText("Fruit: " + orNone(s)) }
	disableables = append(disableables, selEntry)
	selEntryCard := widget.NewCard("Select Entry", "Type, or pick from a list", container.NewVBox(
		note("A text field with a list of suggestions behind its button; any text is allowed."),
		selEntry, selEntryStatus,
	))
	selEntryCard.Style = widget.CardStyleOutlined

	// --- Search select ---
	searchStatus := widget.NewLabel("Country: none")
	search := widget.NewSearchSelect([]string{
		"Argentina", "Australia", "Brazil", "Canada", "France", "Germany", "India",
		"Japan", "Kenya", "Mexico", "Netherlands", "Norway", "Portugal", "Spain",
	}, func(s string) { searchStatus.SetText("Country: " + orNone(s)) })
	search.SetPreferred([]string{"Japan", "Portugal"})
	search.PlaceHolder = "Choose a country"
	disableables = append(disableables, search)
	searchCard := widget.NewCard("Search Select", "A dropdown you can search", container.NewVBox(
		note("Open it and type to filter. Japan and Portugal are preferred: they come first, above a divider."),
		search, searchStatus,
	))
	searchCard.Style = widget.CardStyleOutlined

	// --- Completion popup ---
	words := []widget.CompletionItem{
		{Text: "bellui", Detail: "the toolkit"},
		{Text: "button", Detail: "widget"},
		{Text: "canvas", Detail: "package"},
		{Text: "card", Detail: "widget"},
		{Text: "carousel", Detail: "widget"},
		{Text: "checkbox", Detail: "widget"},
		{Text: "chip", Detail: "widget"},
		{Text: "container", Detail: "package"},
	}
	compEntry := widget.NewEntry()
	compEntry.PlaceHolder = "Type b or c"
	compStatus := widget.NewLabel("Picked: none")
	popup := widget.NewCompletionPopup(w.Canvas(), nil)
	popup.KeepHostFocus = true
	popup.OnSelected = func(_ int, item widget.CompletionItem) {
		compEntry.SetText(item.Text)
		compStatus.SetText("Picked: " + item.Text)
		popup.Hide()
	}
	compEntry.OnChanged = func(s string) {
		var matches []widget.CompletionItem
		if s != "" {
			for _, it := range words {
				if strings.HasPrefix(it.Text, strings.ToLower(s)) && it.Text != s {
					matches = append(matches, it)
				}
			}
		}
		if len(matches) == 0 {
			popup.Hide()
			return
		}
		popup.SetItems(matches)
		pos := bellui.CurrentApp().Driver().AbsolutePositionForObject(compEntry)
		popup.ShowAtCaret(geo.Rectangle{Pos: pos, Size: compEntry.Size()})
	}
	disableables = append(disableables, compEntry)
	compCard := widget.NewCard("Completion Popup", "Suggestions as you type", container.NewVBox(
		note("The app supplies the suggestions; here, words starting with what you typed. Arrow keys and Enter pick one, Escape closes it."),
		compEntry, compStatus,
	))
	compCard.Style = widget.CardStyleOutlined

	// --- Rich entry ---
	rich := widget.NewRichEntry()
	rich.Wrapping = bellui.TextWrapWord // a long line is otherwise drawn past the box
	rich.SetSegments([]widget.RichTextSegment{
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: "Select some of this text, then make it "},
		&widget.TextSegment{Style: widget.RichTextStyleStrong, Text: "bold"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: ", "},
		&widget.TextSegment{Style: widget.RichTextStyleEmphasis, Text: "italic"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: " or "},
		&widget.TextSegment{Style: widget.RichTextStyleCodeInline, Text: "code"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: "."},
	})
	richBox := container.NewGridWrap(geo.NewSize(320, 110), rich)
	styleButton := func(label string, style widget.RichTextStyle) *widget.Button {
		b := widget.NewButton(label, func() { rich.ApplyStyle(style) })
		b.Importance = widget.LowImportance
		return b
	}
	richTools := container.NewHBox(
		styleButton("Bold", widget.RichTextStyleStrong),
		styleButton("Italic", widget.RichTextStyleEmphasis),
		styleButton("Code", widget.RichTextStyleCodeInline),
		widget.NewButtonWithIcon("", theme.ContentUndoIcon(), rich.Undo),
		widget.NewButtonWithIcon("", theme.ContentRedoIcon(), rich.Redo),
	)
	disableables = append(disableables, rich)
	richCard := widget.NewCard("Rich Entry", "Editable text with styles", container.NewVBox(richTools, richBox))
	richCard.Style = widget.CardStyleOutlined

	// --- Rich text: HTML ---
	fromHTML := widget.NewRichTextFromHTML(`<h3>From HTML</h3>
<p>Paragraphs with <b>bold</b>, <i>italic</i> and <code>code</code>, and a <a href="https://bellui.io">link</a>.</p>
<ul><li>A list item</li><li>Another one</li></ul>
<blockquote>Remote content, such as images, is off unless the app allows it.</blockquote>`)
	fromHTML.Wrapping = bellui.TextWrapWord

	// --- Rich text: segments ---
	link, _ := url.Parse("https://bellui.io")
	fromSegments := widget.NewRichText(
		&widget.TextSegment{Style: widget.RichTextStyleSubHeading, Text: "From segments"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: "Built in Go: "},
		&widget.TextSegment{Style: widget.RichTextStyleStrong, Text: "strong"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: ", "},
		&widget.TextSegment{Style: widget.RichTextStyleEmphasis, Text: "emphasis"},
		&widget.TextSegment{Style: widget.RichTextStyleInline, Text: " and "},
		&widget.HyperlinkSegment{Text: "a link", URL: link},
		&widget.TextSegment{Style: widget.RichTextStyleParagraph, Text: "."},
		&widget.ListSegment{Ordered: true, Items: []widget.RichTextSegment{
			&widget.TextSegment{Style: widget.RichTextStyleParagraph, Text: "An ordered list"},
			&widget.TextSegment{Style: widget.RichTextStyleParagraph, Text: "with two items"},
		}},
		&widget.CodeSegment{Language: "go", Text: "fmt.Println(\"code block\")"},
		&widget.SeparatorSegment{},
		&widget.TableSegment{
			Header: []string{"Widget", "Since"},
			Rows:   [][]string{{"RichText", "1.0"}, {"ChipGroup", "1.3"}},
		},
	)
	fromSegments.Wrapping = bellui.TextWrapWord
	richTextCard := widget.NewCard("Rich Text", "Styled text, read-only", container.NewVBox(
		fromHTML, widget.NewSeparator(), fromSegments,
	))
	richTextCard.Style = widget.CardStyleOutlined

	// --- Text grid ---
	grid := widget.NewTextGridFromString("package main\n\nfunc main() {\n\tprintln(\"hello\")  \n}\n")
	grid.ShowLineNumbers = true
	grid.ShowWhitespace = true
	grid.SetRowStyle(2, &widget.CustomTextGridStyle{BGColor: theme.Color(theme.ColorNameSelection)})
	gridCard := widget.NewCard("Text Grid", "A monospaced grid of cells", container.NewVBox(
		note("For code previews and terminals: line numbers, visible whitespace, and a style per cell (line 3 is highlighted)."),
		grid,
	))
	gridCard.Style = widget.CardStyleOutlined

	disable := widget.NewCheck("Disable everything on this page", func(on bool) {
		for _, d := range disableables {
			if on {
				d.Disable()
			} else {
				d.Enable()
			}
		}
	})

	return container.NewVBox(disable, selEntryCard, searchCard, compCard, richCard, richTextCard, gridCard)
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}
