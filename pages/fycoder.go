// Ported from the former fycoder cmd/demo app.
package pages

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/layout"
	"github.com/belluikit/bellui/widget"

	xeditor "github.com/belluikit/bellui/x/editor"
	"github.com/belluikit/bellui/x/editor/highlight/starlark"
)

const fycoderSampleCode = `# Strategy: momentum_scalper
#
# Buys when price momentum exceeds threshold and sells on reversal.

META = {
    "version": "1.0",
    "use_engine": True,
    "api_key": "sk-live-7Qd9f2KpZ1xR8vN3",
}

PARAMS = {
    "entry_threshold": 0.15,
    "exit_threshold": -0.05,
    "max_hold_games": 3,
}

def evaluate(ctx):
    tick = ctx.tick
    score = ctx.score
    position = ctx.position
    params = ctx.params
    state = ctx.state
    mid = tick.mid

    if score == None or score.completed_rounds < 1:
        return {"action": "hold", "reason": "waiting-for-score"}

    # Calculate momentum from recent price changes.
    prev_mid = state.get("prev_mid", mid)
    momentum = mid - prev_mid
    state["prev_mid"] = mid

    # Entry logic: buy when momentum exceeds threshold.
    if position == None:
        if momentum > params["entry_threshold"]:
            edge = momentum - params["entry_threshold"]
            kelly = calc_kelly(edge, mid)
            return {
                "action": "buy",
                "side": "yes",
                "price": mid + 0.02,
                "kelly": kelly,
                "reason": "momentum-entry",
            }
        return {"action": "hold", "reason": "no-signal"}

    # Exit logic: sell on reversal or max hold.
    games_held = score.completed_rounds - state.get("entry_round", 0)
    if momentum < params["exit_threshold"] or games_held >= params["max_hold_games"]:
        return {
            "action": "sell",
            "price": mid - 0.01,
            "reason": "momentum-exit",
        }

    return {"action": "hold"}
`

// newFycoderHighlighter creates a Starlark highlighter with domain-specific
// builtins registered for the demo strategy script.
func newFycoderHighlighter() *starlark.Highlighter {
	h := starlark.New()
	h.AddBuiltins(
		"calc_kelly",
		"cost_per_contract",
		"point_ordinal",
		"undamped_prob",
		"trailing_regret_floor",
		"trailing_regret_threshold",
		"trailing_regret_reversals",
		"competitiveness",
		"volume_rate",
		"get_position",
		"min_f", "max_f",
	)
	return h
}

// spanOf locates the first occurrence of substr in text and returns its
// character range in 0-based (row, col) coordinates (rune columns). The bool
// reports whether the substring was found.
func spanOf(text, substr string) (xeditor.TextRange, bool) {
	for row, line := range strings.Split(text, "\n") {
		if idx := strings.Index(line, substr); idx >= 0 {
			startCol := len([]rune(line[:idx]))
			return xeditor.TextRange{
				StartRow: row,
				StartCol: startCol,
				EndRow:   row,
				EndCol:   startCol + len([]rune(substr)),
			}, true
		}
	}
	return xeditor.TextRange{}, false
}

// fycoderDiagnostics builds the Phase 3 diagnostics demo: one diagnostic per
// severity, anchored to real spans in the sample. Each draws a wavy squiggle,
// a severity-distinct gutter marker (circle/triangle/square/diamond, #56), and
// a hover card showing its message — and is reachable by clicking its marker.
func fycoderDiagnostics() []xeditor.Diagnostic {
	var diags []xeditor.Diagnostic
	add := func(substr string, sev xeditor.Severity, msg, source string) {
		if r, ok := spanOf(fycoderSampleCode, substr); ok {
			diags = append(diags, xeditor.Diagnostic{Range: r, Severity: sev, Message: msg, Source: source})
		}
	}
	add("calc_kelly", xeditor.SeverityError, "undefined name 'calc_kelly'", "typecheck")
	add("score == None", xeditor.SeverityWarning, "comparison to None should use 'is None'", "lint")
	add("max_hold_games", xeditor.SeverityInfo, "parameter read from PARAMS — document its units", "lint")
	add("def evaluate", xeditor.SeverityHint, "strategy entry point — consider adding a docstring", "lint")
	return diags
}

// fycoderCompletionWords backs the Ctrl+Space completion popup: Starlark
// keywords, common builtins, and the demo strategy's own identifiers.
var fycoderCompletionWords = []string{
	"def", "return", "if", "elif", "else", "for", "while", "and", "or", "not",
	"None", "True", "False", "pass", "break", "continue", "lambda", "in",
	"len", "range", "dict", "list", "min", "max", "abs", "round",
	"calc_kelly", "cost_per_contract", "get_position", "min_f", "max_f",
	"evaluate", "ctx", "tick", "score", "position", "params", "state", "mid",
	"momentum", "entry_threshold", "exit_threshold", "max_hold_games",
}

// fycoderHoverDocs backs the hover cards (mouse dwell, or Ctrl+K Ctrl+I) shown
// over keywords, builtins, and strategy identifiers.
var fycoderHoverDocs = map[string]string{
	"def":        "**def** — define a function.\n\n```\ndef name(args):\n    ...\n```",
	"return":     "**return** — exit a function, optionally yielding a value.",
	"None":       "**None** — the null value. Test identity with `is None`, not `==`.",
	"calc_kelly": "**calc_kelly(edge, mid)** *(builtin)*\n\nKelly-criterion position size for the given `edge` at price `mid`.",
	"evaluate":   "**evaluate(ctx)**\n\nStrategy entry point — called once per tick, returns an action dict.",
	"momentum":   "Local: recent price change, `mid - prev_mid`. Drives entry/exit.",
}

// fycoderSignatures and fycoderOverloads back signature help (typing "(" / ","
// or Ctrl+Shift+Space). calc_kelly carries a second overload so the ▲ N/M ▾
// overload control (Alt+Up/Down, or click) is exercised.
var fycoderSignatures = map[string]xeditor.SignatureInformation{
	"calc_kelly": {
		Label:         "calc_kelly(edge, mid) -> float",
		Documentation: "Kelly-criterion position size for an edge at price mid.",
		Parameters:    []xeditor.ParameterInformation{{Label: "edge"}, {Label: "mid"}},
	},
	"evaluate": {
		Label:         "evaluate(ctx) -> dict",
		Documentation: "Strategy entry point, called once per tick.",
		Parameters:    []xeditor.ParameterInformation{{Label: "ctx"}},
	},
}

var fycoderOverloads = map[string][]xeditor.SignatureInformation{
	"calc_kelly": {{
		Label:         "calc_kelly(edge, mid, cap) -> float",
		Documentation: "Kelly size capped at fraction `cap` of bankroll.",
		Parameters:    []xeditor.ParameterInformation{{Label: "edge"}, {Label: "mid"}, {Label: "cap"}},
	}},
}

// overloadSignatureProvider wraps a [xeditor.FuncSignatureProvider] (which finds
// the enclosing call and active parameter) and expands the matched function to
// its full overload set, so the demo exercises the multi-signature UI.
type overloadSignatureProvider struct {
	base      *xeditor.FuncSignatureProvider
	overloads map[string][]xeditor.SignatureInformation
}

func (p *overloadSignatureProvider) Signatures(ctx context.Context, row, col int) (*xeditor.SignatureHelp, error) {
	help, err := p.base.Signatures(ctx, row, col)
	if err != nil || help == nil || len(help.Signatures) == 0 {
		return help, err
	}
	label := help.Signatures[0].Label
	name := label
	if i := strings.IndexByte(label, '('); i >= 0 {
		name = label[:i]
	}
	if extra := p.overloads[name]; len(extra) > 0 {
		help.Signatures = append(help.Signatures, extra...)
	}
	return help, nil
}

// wordAtPos returns the identifier covering (row, col), or "" if none — used to
// resolve the Cmd/Ctrl-click navigation target.
func wordAtPos(text string, row, col int) string {
	lines := strings.Split(text, "\n")
	if row < 0 || row >= len(lines) {
		return ""
	}
	runes := []rune(lines[row])
	if col < 0 || col > len(runes) {
		return ""
	}
	isIdent := func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }
	start, end := col, col
	for start > 0 && isIdent(runes[start-1]) {
		start--
	}
	for end < len(runes) && isIdent(runes[end]) {
		end++
	}
	if start == end {
		return ""
	}
	return string(runes[start:end])
}

// fycoderSampleFolds derives the sample script's foldable regions from its
// content (not hard-coded line numbers, so it survives edits to the sample):
// the header comment block, the META and PARAMS dict literals, and the
// evaluate() function body. PARAMS loads pre-collapsed with a summary
// placeholder to show auto-fold + the ellipsis/summary on first paint.
func fycoderSampleFolds() []xeditor.FoldRegion {
	lines := strings.Split(fycoderSampleCode, "\n")
	var regions []xeditor.FoldRegion

	// Header comment block: the first contiguous run of "#" lines.
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		end := 0
		for end+1 < len(lines) && strings.HasPrefix(lines[end+1], "#") {
			end++
		}
		if end > 0 {
			regions = append(regions, xeditor.FoldRegion{
				Range: xeditor.TextRange{StartRow: 0, EndRow: end},
				Kind:  xeditor.FoldKindComment,
			})
		}
	}

	// A top-level "NAME = {" dict literal folds down to its closing "}".
	dict := func(header string, collapsed bool, summary string) {
		start := -1
		for i, l := range lines {
			if strings.HasPrefix(l, header) {
				start = i
				break
			}
		}
		if start < 0 {
			return
		}
		for i := start + 1; i < len(lines); i++ {
			if strings.HasPrefix(lines[i], "}") {
				regions = append(regions, xeditor.FoldRegion{
					Range:     xeditor.TextRange{StartRow: start, EndRow: i},
					Kind:      xeditor.FoldKindBlock,
					Collapsed: collapsed,
					Summary:   summary,
				})
				return
			}
		}
	}
	dict("META = {", false, "")
	dict("PARAMS = {", true, " { ... 3 params }")

	// The evaluate() function: header line down to the last non-blank line.
	for i, l := range lines {
		if strings.HasPrefix(l, "def evaluate") {
			end := len(lines) - 1
			for end > i && strings.TrimSpace(lines[end]) == "" {
				end--
			}
			if end > i {
				regions = append(regions, xeditor.FoldRegion{
					Range: xeditor.TextRange{StartRow: i, EndRow: end},
					Kind:  xeditor.FoldKindBlock,
				})
			}
			break
		}
	}
	return regions
}

// bigFycoderFile generates an n-line document to exercise large-file
// virtualisation + viewport-windowed gutter/decoration rendering (Phase 5
// deliverable 5).
func bigFycoderFile(n int) string {
	var b strings.Builder
	b.Grow(n * 40)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "def rule_%05d(ctx):  # generated rule %d — scales with the viewport, not the file\n", i, i)
	}
	return b.String()
}

// FycoderPage demonstrates the fycoder code editor widget: syntax
// highlighting, line numbers, breakpoints, current-line decorator,
// span decorations (squiggle + redaction box), code folding, find bar and
// save handling.
func FycoderPage(_ bellui.Window) core.CanvasObject {
	// Status bar.
	cursorLabel := widget.NewLabel("Ln 1, Col 1")
	cursorLabel.TextStyle = bellui.TextStyle{Monospace: true}
	breakpointLabel := widget.NewLabel("0 breakpoints")
	breakpointLabel.TextStyle = bellui.TextStyle{Monospace: true}
	statusLabel := widget.NewLabel("Ready")

	// Breakpoint gutter with status update.
	bpGutter := xeditor.NewBreakpointGutter(func(lines map[int]bool) {
		count := len(lines)
		if count == 1 {
			breakpointLabel.SetText("1 breakpoint")
		} else {
			breakpointLabel.SetText(fmt.Sprintf("%d breakpoints", count))
		}
	})

	// Current line decorator — highlights the cursor's row.
	currentLine := xeditor.NewCurrentLineDecorator(0)

	// Editor.
	editor := xeditor.NewEditor()
	editor.Highlighter = newFycoderHighlighter()
	// Composite gutter: line numbers, Phase 5 fold chevrons, then Phase 3
	// diagnostic markers, then breakpoints. Column-aware tap routing (#95) keeps
	// the interactive columns from interfering — a tap folds, reveals, or toggles,
	// never several at once.
	editor.Gutter = xeditor.NewCompositeGutter(
		xeditor.NewLineNumberGutter(),
		editor.FoldGutter(),
		editor.DiagnosticGutter(),
		bpGutter,
	)
	editor.Decorators = []xeditor.RowDecorator{currentLine}

	// Phase 3 diagnostics: squiggles + severity-shaped gutter markers + hover.
	// Click a marker to reveal its diagnostic; repeated clicks cycle a line's
	// diagnostics (#57, #97, #99).
	editor.SetDiagnostics(fycoderDiagnostics())
	editor.SetDiagnosticClickToReveal(true)
	editor.DiagnosticGutter().CycleOnTap = true

	// Phase 2 floating surfaces: completion, hover, and signature help.
	editor.SetCompletionProvider(&xeditor.WordListCompletionProvider{
		Words: fycoderCompletionWords,
		Kind:  xeditor.CompletionKeyword,
	})
	editor.SetHoverProvider(&xeditor.WordHoverProvider{
		Docs: fycoderHoverDocs,
		Text: editor.Text,
	})
	editor.SetSignatureHelpProvider(&overloadSignatureProvider{
		base:      &xeditor.FuncSignatureProvider{Funcs: fycoderSignatures, Text: editor.Text},
		overloads: fycoderOverloads,
	})

	// Phase 3 navigation: Cmd/Ctrl-hover lights a link over navigable symbols;
	// Cmd/Ctrl-click jumps to the definition (or reports a builtin).
	editor.SetLinkProvider(&xeditor.WordLinkProvider{
		Navigable: map[string]bool{"evaluate": true, "calc_kelly": true},
		Text:      editor.Text,
	})
	editor.OnRequestNavigation(func(row, col int) {
		word := wordAtPos(editor.Text(), row, col)
		if word == "" {
			return
		}
		if r, ok := spanOf(editor.Text(), "def "+word); ok {
			// Reveal the function name (skip the leading "def ").
			editor.RevealRange(xeditor.TextRange{
				StartRow: r.StartRow, StartCol: r.StartCol + len("def "),
				EndRow: r.EndRow, EndCol: r.EndCol,
			})
			statusLabel.SetText("Navigated to def " + word)
			return
		}
		statusLabel.SetText(word + ": builtin (no local definition)")
	})

	editor.OnCursorChanged = func(row, col int) {
		cursorLabel.SetText(fmt.Sprintf("Ln %d, Col %d", row+1, col+1))
		currentLine.Row = row + 1
		editor.RefreshDecorators()
	}

	cleanText := fycoderSampleCode

	editor.OnSave = func() {
		cleanText = editor.Text()
		statusLabel.SetText("Saved (demo — no actual file written)")
	}

	editor.OnChanged = func(s string) {
		if s == cleanText {
			statusLabel.SetText("Ready")
		} else {
			statusLabel.SetText("Modified")
		}
	}

	editor.SetText(fycoderSampleCode)

	// Phase 4 redaction: mask the secret API key for real. Unlike a cosmetic
	// overlay, a masked redaction bars the caret inside its range and silently
	// vetoes edits, so typing can't push the hidden text out from under the bar.
	redactSample := func() {
		if r, ok := spanOf(fycoderSampleCode, "sk-live-7Qd9f2KpZ1xR8vN3"); ok {
			editor.SetRedactions([]xeditor.TextRange{r})
		}
	}
	redactSample()

	// Phase 5 folding: register the sample's foldable regions (chevrons appear in
	// the fold gutter); PARAMS loads pre-collapsed with a summary placeholder.
	sampleFolds := fycoderSampleFolds()
	editor.SetFoldRegions(sampleFolds)

	// Fold controls.
	foldAllBtn := widget.NewButton("Fold all", func() {
		for _, r := range sampleFolds {
			editor.Fold(r.Range)
		}
	})
	unfoldAllBtn := widget.NewButton("Unfold all", func() {
		for _, r := range sampleFolds {
			editor.Unfold(r.Range)
		}
	})

	// Large-file toggle: swap between the sample and a 50k-line generated file to
	// show virtualisation + the viewport-windowed gutter/decoration rendering. The
	// sample-only diagnostics/redaction are cleared for the big file and restored
	// on return.
	const bigLines = 50000
	var showingBig bool
	var bigBtn *widget.Button
	bigBtn = widget.NewButton(fmt.Sprintf("Load %d-line file", bigLines), func() {
		if showingBig {
			editor.SetText(fycoderSampleCode)
			editor.SetDiagnostics(fycoderDiagnostics())
			redactSample()
			editor.SetFoldRegions(sampleFolds)
			bigBtn.SetText(fmt.Sprintf("Load %d-line file", bigLines))
			statusLabel.SetText("Ready")
		} else {
			editor.SetText(bigFycoderFile(bigLines))
			editor.SetDiagnostics(nil)
			editor.SetRedactions(nil)
			editor.SetFoldRegions(nil)
			// Auto-fold a huge interior region — pxfed's "collapse the giant block"
			// case — leaving the file responsive.
			editor.AutoFold([]xeditor.TextRange{{StartRow: 10, EndRow: bigLines - 10}})
			bigBtn.SetText("Load sample")
			statusLabel.SetText(fmt.Sprintf("%d lines — folded, viewport-windowed", bigLines))
		}
		showingBig = !showingBig
	})

	toolbar := container.NewHBox(foldAllBtn, unfoldAllBtn, widget.NewSeparator(), bigBtn)

	// Status bar layout.
	statusBar := container.NewHBox(
		cursorLabel,
		widget.NewSeparator(),
		breakpointLabel,
		layout.NewSpacer(),
		statusLabel,
	)

	// Help text.
	helpLabel := widget.NewLabel(
		"Ctrl+Space: Complete  |  Hover or Ctrl+K Ctrl+I: Docs  |  Ctrl+Shift+Space: Signature (Alt+↑/↓ overloads)  |  Cmd/Ctrl+Click: Go to def\n" +
			"Cmd+F: Find  |  Cmd+S: Save  |  Click a gutter chevron: Fold  |  Click marker: Reveal diagnostic  |  Gray box = redacted secret",
	)
	helpLabel.TextStyle = bellui.TextStyle{Monospace: true}
	helpLabel.Importance = widget.LowImportance

	content := container.NewBorder(
		container.NewVBox(toolbar, helpLabel), // top
		statusBar,                             // bottom
		nil, nil,
		editor,
	)

	return pageMinSize(900, 600, content)
}
