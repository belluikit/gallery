package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
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
	v := bellui.CurrentApp().Settings().ThemeVariant()

	var rows []core.CanvasObject
	for _, tok := range tokens {
		size := th.Size(tok.size)

		// Token name and specs.
		specs := fmt.Sprintf("%s — %.0fpx, weight %d", tok.name, size, int(tok.style.ResolvedWeight()))
		specLabel := canvas.NewText(specs, th.Color(theme.ColorNameOnSurfaceVariant, v))
		specLabel.TextSize = 11

		// Sample text with the actual style applied.
		sample := canvas.NewText(tok.sample, th.Color(theme.ColorNameOnSurface, v))
		sample.TextSize = size
		sample.TextStyle = tok.style

		card := widget.NewCard("", "", container.NewVBox(specLabel, sample))
		card.Style = widget.CardStyleOutlined
		rows = append(rows, card)
	}

	return container.NewVBox(rows...)
}
