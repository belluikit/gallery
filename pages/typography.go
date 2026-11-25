package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"
)

// TypographyPage shows all 15 MD3 type scale tokens.
func TypographyPage(_ bellui.Window) core.CanvasObject {
	type typeToken struct {
		name   string
		size   bellui.ThemeSizeName
		style  bellui.TextStyle
		sample string
	}

	tokens := []typeToken{
		{"Display Large", theme.SizeNameDisplayLarge, theme.TypographyDisplayLarge, "Display Large"},
		{"Display Medium", theme.SizeNameDisplayMedium, theme.TypographyDisplayMedium, "Display Medium"},
		{"Display Small", theme.SizeNameDisplaySmall, theme.TypographyDisplaySmall, "Display Small"},
		{"Headline Large", theme.SizeNameHeadlineLarge, theme.TypographyHeadlineLarge, "Headline Large"},
		{"Headline Medium", theme.SizeNameHeadlineMedium, theme.TypographyHeadlineMedium, "Headline Medium"},
		{"Headline Small", theme.SizeNameHeadlineSmall, theme.TypographyHeadlineSmall, "Headline Small"},
		{"Title Large", theme.SizeNameTitleLarge, theme.TypographyTitleLarge, "Title Large"},
		{"Title Medium", theme.SizeNameTitleMedium, theme.TypographyTitleMedium, "Title Medium"},
		{"Title Small", theme.SizeNameTitleSmall, theme.TypographyTitleSmall, "Title Small"},
		{"Body Large", theme.SizeNameBodyLarge, theme.TypographyBodyLarge, "Body Large — The quick brown fox jumps over the lazy dog"},
		{"Body Medium", theme.SizeNameBodyMedium, theme.TypographyBodyMedium, "Body Medium — The quick brown fox jumps over the lazy dog"},
		{"Body Small", theme.SizeNameBodySmall, theme.TypographyBodySmall, "Body Small — The quick brown fox jumps over the lazy dog"},
		{"Label Large", theme.SizeNameLabelLarge, theme.TypographyLabelLarge, "LABEL LARGE"},
		{"Label Medium", theme.SizeNameLabelMedium, theme.TypographyLabelMedium, "LABEL MEDIUM"},
		{"Label Small", theme.SizeNameLabelSmall, theme.TypographyLabelSmall, "LABEL SMALL"},
	}

	th := theme.Current()

	var rows []core.CanvasObject
	for _, tok := range tokens {
		size := th.Size(tok.size)

		// The token's name and specs, then a sample in the style itself. One
		// RichText holds both, so its padding goes round the pair once. The
		// sample wraps, so one wider than a phone runs onto a second line
		// instead of making the page scroll sideways; that also shows the
		// style's line height (gallery#51).
		specs := fmt.Sprintf("%s — %.0fpx, weight %d", tok.name, size, int(tok.style.ResolvedWeight()))
		text := widget.NewRichText(
			&widget.TextSegment{Text: specs, Style: widget.RichTextStyle{
				SizeName: theme.SizeNameLabelSmall, ColorName: theme.ColorNameOnSurfaceVariant,
			}},
			&widget.TextSegment{Text: tok.sample, Style: widget.RichTextStyle{
				SizeName: tok.size, ColorName: theme.ColorNameOnSurface, TextStyle: tok.style,
			}},
		)
		text.Wrapping = bellui.TextWrapWord

		card := widget.NewCard("", "", text)
		card.Style = widget.CardStyleOutlined
		rows = append(rows, card)
	}

	return container.NewVBox(rows...)
}
