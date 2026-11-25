//go:build !android

package pages

import (
	"strings"
	"testing"

	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	xeditor "github.com/belluikit/bellui/x/editor"
)

// The code editor page draws each of its diagnostics in the gutter: the
// editor's DiagnosticGutter is in the page's composite gutter, and a line
// with an error gets the error's dot, one with a warning the warning's
// triangle (gallery#47).
func TestFycoderGutterShowsTheDiagnostics(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	ed := fycoderEditor(t, FycoderPage(w))

	gutter, ok := ed.Gutter.(*xeditor.CompositeGutter)
	if !ok {
		t.Fatalf("the editor's gutter is a %T, not a composite", ed.Gutter)
	}
	inGutter := false
	for _, p := range gutter.Providers() {
		if p == xeditor.GutterProvider(ed.DiagnosticGutter()) {
			inGutter = true
		}
	}
	if !inGutter {
		t.Fatal("the diagnostic gutter is not one of the editor's gutter columns")
	}

	cell := geo.NewSize(20, 20)
	marker := func(row int) []core.CanvasObject {
		return ed.DiagnosticGutter().Render(row+1, cell, 14) // the gutter's lines are 1-based
	}
	marked := map[int]bool{}
	for _, d := range fycoderDiagnostics() {
		marked[d.Range.StartRow] = true
		objs := marker(d.Range.StartRow)
		if len(objs) == 0 {
			t.Errorf("no gutter marker on line %d for %q", d.Range.StartRow+1, d.Message)
			continue
		}
		switch d.Severity {
		case xeditor.SeverityError:
			if _, ok := objs[0].(*canvas.Circle); !ok {
				t.Errorf("the error on line %d is marked with a %T, not a dot", d.Range.StartRow+1, objs[0])
			}
		case xeditor.SeverityWarning:
			if _, ok := objs[0].(*canvas.Polygon); !ok {
				t.Errorf("the warning on line %d is marked with a %T, not a triangle", d.Range.StartRow+1, objs[0])
			}
		}
	}
	for row := 0; row < 5; row++ {
		if !marked[row] && len(marker(row)) != 0 {
			t.Errorf("line %d has no diagnostic but a gutter marker", row+1)
		}
	}
}

// Typing in the code editor page's field marks every row that contains the
// text, and only those, through a SearchHighlightDecorator (gallery#47).
func TestFycoderMarksTheRowsThatMatch(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	page := FycoderPage(w)
	ed := fycoderEditor(t, page)

	var marked *xeditor.SearchHighlightDecorator
	for _, d := range ed.Decorators {
		if s, ok := d.(*xeditor.SearchHighlightDecorator); ok {
			marked = s
		}
	}
	if marked == nil {
		t.Fatal("the editor has no SearchHighlightDecorator")
	}
	var field *widget.Entry
	walkPage(page, func(o core.CanvasObject) {
		if e, ok := o.(*widget.Entry); ok && e.PlaceHolder == "Mark rows containing…" {
			field = e
		}
	})
	if field == nil {
		t.Fatal("no field to mark rows with")
	}

	const query = "momentum"
	want := map[int]bool{}
	for i, line := range strings.Split(ed.Text(), "\n") {
		if strings.Contains(line, query) {
			want[i+1] = true
		}
	}
	if len(want) < 2 {
		t.Fatalf("the sample has %d rows with %q; the test needs a few", len(want), query)
	}

	test.Type(field, query)
	if len(marked.Rows) != len(want) {
		t.Errorf("%d rows marked, want %d", len(marked.Rows), len(want))
	}
	for line := range want {
		if !marked.Rows[line] {
			t.Errorf("line %d contains %q but isn't marked", line, query)
		}
		if marked.Decorate(line, geo.Position{}, geo.NewSize(100, 20)) == nil {
			t.Errorf("line %d is marked but the decorator draws nothing on it", line)
		}
	}
	for line := range marked.Rows {
		if !want[line] {
			t.Errorf("line %d is marked but doesn't contain %q", line, query)
		}
	}

	field.SetText("")
	if len(marked.Rows) != 0 {
		t.Errorf("%d rows still marked after the field was cleared", len(marked.Rows))
	}
}

// fycoderEditor finds the editor on the code editor page.
func fycoderEditor(t *testing.T, page core.CanvasObject) *xeditor.Editor {
	t.Helper()
	var ed *xeditor.Editor
	walkPage(page, func(o core.CanvasObject) {
		if e, ok := o.(*xeditor.Editor); ok {
			ed = e
		}
	})
	if ed == nil {
		t.Fatal("no editor on the page")
	}
	return ed
}
