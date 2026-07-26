package pages

import (
	"testing"

	xeditor "github.com/belluikit/bellui/x/editor"
)

func TestSpanOf(t *testing.T) {
	text := "alpha\nbeta gamma\ndelta"
	r, ok := spanOf(text, "gamma")
	if !ok {
		t.Fatal("expected to find 'gamma'")
	}
	want := xeditor.TextRange{StartRow: 1, StartCol: 5, EndRow: 1, EndCol: 10}
	if r != want {
		t.Errorf("spanOf = %+v, want %+v", r, want)
	}
	if _, ok := spanOf(text, "zzz"); ok {
		t.Error("expected 'zzz' not to be found")
	}
}

func TestFycoderDiagnostics(t *testing.T) {
	diags := fycoderDiagnostics()
	// One diagnostic per severity, each anchored to a real span in the sample.
	wantSev := map[xeditor.Severity]bool{
		xeditor.SeverityError:   false,
		xeditor.SeverityWarning: false,
		xeditor.SeverityInfo:    false,
		xeditor.SeverityHint:    false,
	}
	for _, d := range diags {
		if _, ok := wantSev[d.Severity]; !ok {
			t.Errorf("unexpected severity %v", d.Severity)
		}
		wantSev[d.Severity] = true
		if d.Message == "" {
			t.Errorf("diagnostic at %+v has no message", d.Range)
		}
		if d.Range.IsEmpty() {
			t.Errorf("diagnostic %q resolved to an empty range (substring not found?)", d.Message)
		}
	}
	for sev, seen := range wantSev {
		if !seen {
			t.Errorf("missing a diagnostic of severity %v", sev)
		}
	}
}

func TestWordAtPos(t *testing.T) {
	text := "def evaluate(ctx):\n    return None"
	if got := wordAtPos(text, 0, 5); got != "evaluate" {
		t.Errorf("wordAtPos on 'evaluate' = %q, want %q", got, "evaluate")
	}
	if got := wordAtPos(text, 1, 11); got != "None" {
		t.Errorf("wordAtPos on 'None' = %q, want %q", got, "None")
	}
	if got := wordAtPos(text, 1, 2); got != "" { // inside the leading indentation
		t.Errorf("wordAtPos on whitespace = %q, want empty", got)
	}
}
