package pages

import (
	"fmt"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"
)

// FontPickerPage returns the Font Picker demo page.
func FontPickerPage(w bellui.Window) core.CanvasObject {
	// Inline picker
	config := widget.FontPickerConfig{
		Family: "Roboto",
		Size:   16,
		Weight: bellui.WeightNormal,
	}
	resultLabel := widget.NewLabel("Selected: Roboto, 16pt, Normal")
	picker := widget.NewFontPickerWithSystemFonts(config, func(c widget.FontPickerConfig) {
		resultLabel.SetText(fmt.Sprintf("Selected: %s, %.0fpt, %s, bold=%v, italic=%v, underline=%v",
			c.Family, c.Size, nameForWeight(c.Weight), c.Bold, c.Italic, c.Underline))
	})
	inlineCard := widget.NewCard("Inline Font Picker", "Select font properties inline", container.NewVBox(picker, resultLabel))
	inlineCard.Style = widget.CardStyleOutlined

	// Popup picker
	popupLabel := widget.NewLabel("No font selected yet")
	popupBtn := &widget.Button{
		Text:  "Open Font Picker",
		Style: widget.ButtonStyleOutlined,
		OnTapped: func() {
			widget.ShowFontPickerPopup(config, nil, func(c widget.FontPickerConfig) {
				popupLabel.SetText(fmt.Sprintf("%s, %.0fpt, %s", c.Family, c.Size, nameForWeight(c.Weight)))
			}, w.Canvas())
		},
	}
	popupCard := widget.NewCard("Font Picker Popup", "Modal dialog variant", container.NewVBox(popupBtn, popupLabel))
	popupCard.Style = widget.CardStyleOutlined

	return container.NewVBox(inlineCard, popupCard)
}

func nameForWeight(w bellui.TextWeight) string {
	switch w {
	case bellui.WeightThin:
		return "Thin"
	case bellui.WeightExtraLight:
		return "ExtraLight"
	case bellui.WeightLight:
		return "Light"
	case bellui.WeightNormal:
		return "Normal"
	case bellui.WeightMedium:
		return "Medium"
	case bellui.WeightSemiBold:
		return "SemiBold"
	case bellui.WeightBold:
		return "Bold"
	case bellui.WeightExtraBold:
		return "ExtraBold"
	case bellui.WeightBlack:
		return "Black"
	default:
		return fmt.Sprintf("W%.0f", float32(w))
	}
}
