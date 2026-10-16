//go:build !android && !js

// Excluded from Android builds: clair's knowledge base links FAISS
// (go-faiss, -lfaiss_c) for vector search, a host-only native library not
// part of the mobile build. Omitted on mobile via pages/native_mobile.go.
package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/clair/assistant"
	"github.com/belluikit/clair/kb"
	"github.com/belluikit/clair/kb/builder"
	"github.com/belluikit/clair/parser/goparser"
)

// clairAssistant is lazily initialized on first use.
var (
	clairOnce      sync.Once
	clairAssistant *assistant.Assistant
	clairErr       error
)

func getClairAssistant() (*assistant.Assistant, error) {
	clairOnce.Do(func() {
		// Build a framework index from Bellui's widget package.
		// In production, this would be a pre-built index shipped with the IDE.
		indexDir := filepath.Join(os.TempDir(), "clair_gallery_demo.bleve")

		// Try to open an existing index first.
		fw, err := kb.OpenFramework(indexDir)
		if err != nil {
			// Build a new one.
			b := builder.New()
			p := goparser.New("github.com/belluikit/bellui")

			// Find the Bellui widget source directory relative to the GOPATH.
			// For this demo, we use a known path.
			belluiRoot := findBelluiRoot()
			if belluiRoot == "" {
				clairErr = fmt.Errorf("could not locate Bellui source tree")
				return
			}

			widgetDir := filepath.Join(belluiRoot, "widget")
			if err := b.AddFromParser(p, widgetDir); err != nil {
				clairErr = fmt.Errorf("indexing widgets: %w", err)
				return
			}

			// Also index container package.
			containerDir := filepath.Join(belluiRoot, "container")
			if info, err := os.Stat(containerDir); err == nil && info.IsDir() {
				b.AddFromParser(p, containerDir)
			}

			if err := b.Build(indexDir); err != nil {
				clairErr = fmt.Errorf("building index: %w", err)
				return
			}

			fw, err = kb.OpenFramework(indexDir)
			if err != nil {
				clairErr = fmt.Errorf("opening index: %w", err)
				return
			}
		}

		clairAssistant = assistant.New(fw, nil, assistant.WithSystemPrompt(belluiSystemPrompt))
	})

	return clairAssistant, clairErr
}

func findBelluiRoot() string {
	// Try common locations.
	candidates := []string{
		"../bellui",    // relative to gallery
		"../../bellui", // other layouts
	}

	// Check GOPATH-based paths.
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		candidates = append(candidates, filepath.Join(gopath, "src/github.com/belluikit/bellui"))
	}

	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(abs, "widget")); err == nil {
			return abs
		}
	}
	return ""
}

// ClairPage demonstrates the Clair knowledge base search.
func ClairPage(_ bellui.Window) core.CanvasObject {
	assist, err := getClairAssistant()
	if err != nil {
		errCard := widget.NewCard("Clair", "Knowledge Base Error",
			widget.NewLabel(fmt.Sprintf("Failed to initialize: %v", err)))
		errCard.Style = widget.CardStyleOutlined
		return container.NewVBox(errCard)
	}

	// Search input.
	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "Search Bellui API (e.g. \"button\", \"list widget\", \"text input\")"

	// Results display.
	resultItems := []core.CanvasObject{}
	resultsBox := container.NewVBox(resultItems...)

	statusLabel := widget.NewLabel("Type a query and press Enter to search.")
	statusLabel.TextStyle = bellui.TextStyle{Italic: true}

	doSearch := func(query string) {
		if query == "" {
			return
		}

		results, searchErr := assist.Builtin.Search(query, 8)
		if searchErr != nil {
			statusLabel.SetText(fmt.Sprintf("Search error: %v", searchErr))
			return
		}

		statusLabel.SetText(fmt.Sprintf("Found %d results for %q", len(results), query))

		resultsBox.Objects = nil
		for _, r := range results {
			r := r // capture

			// Build result card content.
			kindLabel := widget.NewLabel(r.Kind.String())
			kindLabel.TextStyle = bellui.TextStyle{Bold: true}

			idLabel := widget.NewLabel(r.ID)
			idLabel.TextStyle = bellui.TextStyle{Bold: true}

			header := container.NewHBox(idLabel, kindLabel)

			var lines []core.CanvasObject
			lines = append(lines, header)

			if r.Signature != "" {
				sig := widget.NewLabel(r.Signature)
				sig.TextStyle = bellui.TextStyle{Monospace: true}
				sig.Wrapping = bellui.TextWrapWord
				lines = append(lines, sig)
			}

			if r.Doc != "" {
				doc := r.Doc
				if len(doc) > 200 {
					doc = doc[:197] + "..."
				}
				docLabel := widget.NewLabel(doc)
				docLabel.Wrapping = bellui.TextWrapWord
				lines = append(lines, docLabel)
			}

			scoreLabel := widget.NewLabel(fmt.Sprintf("Score: %.2f  |  %s", r.Score, r.Tier))
			scoreLabel.TextStyle = bellui.TextStyle{Italic: true}
			lines = append(lines, scoreLabel)

			card := widget.NewCard("", "", container.NewVBox(lines...))
			card.Style = widget.CardStyleOutlined
			resultsBox.Objects = append(resultsBox.Objects, card)
		}

		if len(results) == 0 {
			noResults := widget.NewLabel("No results found.")
			noResults.Alignment = bellui.TextAlignCenter
			resultsBox.Objects = append(resultsBox.Objects, noResults)
		}

		resultsBox.Refresh()
	}

	searchEntry.OnSubmitted = doSearch

	searchBtn := &widget.Button{
		Text:  "Search",
		Icon:  theme.SearchIcon(),
		Style: widget.ButtonStyleFilled,
		OnTapped: func() {
			doSearch(searchEntry.Text)
		},
	}

	searchRow := container.NewBorder(nil, nil, nil, searchBtn, searchEntry)

	// Info card.
	infoCard := widget.NewCard("Clair — Knowledge Base Demo", "Search the Bellui framework API using Clair's two-tier knowledge base",
		container.NewVBox(
			widget.NewLabel("This demo indexes Bellui's widget and container packages into a Bleve full-text search index,\nthen provides instant offline search via the BuiltinAssistant."),
			widget.NewSeparator(),
			searchRow,
			statusLabel,
		))
	infoCard.Style = widget.CardStyleOutlined

	return container.NewVBox(infoCard, resultsBox)
}
