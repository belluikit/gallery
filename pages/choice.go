package pages

import (
	"fmt"
	"strings"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"
)

// ChoicePage shows the widgets for picking from a set of options: ChipGroup,
// with single, multiple and three-state (selected and excluded) chips,
// CheckGroup and SpinBox (gallery#38).
func ChoicePage(_ bellui.Window) core.CanvasObject {
	var disableables []core.Disableable

	// --- Chip groups ---
	single := widget.NewChipGroup([]string{"Small", "Medium", "Large"}, false, nil)
	single.Required = true
	single.SetSelected([]string{"Medium"})
	singleStatus := widget.NewLabel("Size: Medium")
	single.OnChanged = func(sel []string) { singleStatus.SetText("Size: " + joinOrNone(sel)) }

	multi := widget.NewChipGroup([]string{"Running", "Cycling", "Swimming", "Hiking", "Climbing", "Rowing"}, true, nil)
	multiStatus := widget.NewLabel("Activities: none")
	multi.OnChanged = func(sel []string) { multiStatus.SetText("Activities: " + joinOrNone(sel)) }

	three := widget.NewChipGroup([]string{"Rock", "Jazz", "Pop", "Folk", "Classical", "Electronic"}, true, nil)
	three.Excludable = true
	threeStatus := widget.NewLabel("")
	threeStatus.Wrapping = bellui.TextWrapWord
	last := "nothing yet"
	showThree := func() {
		threeStatus.SetText(fmt.Sprintf("Include: %s\nExclude: %s\nLast change: %s",
			joinOrNone(three.Selected), joinOrNone(three.Excluded), last))
	}
	three.OnStateChanged = func(opt string, s widget.ChipState) {
		last = opt + " is now " + chipStateName(s)
		showThree()
	}
	showThree()

	align := widget.NewSelect([]string{"Start", "Center", "End"}, func(s string) {
		a := map[string]bellui.TextAlign{
			"Start": bellui.TextAlignLeading, "Center": bellui.TextAlignCenter, "End": bellui.TextAlignTrailing,
		}[s]
		for _, g := range []*widget.ChipGroup{single, multi, three} {
			g.Alignment = a
			g.Refresh()
		}
	})
	align.SetSelected("Start")
	disableables = append(disableables, single, multi, three)

	chipCard := widget.NewCard("Chip Group", "Filter chips to pick from, wrapping onto the next row", container.NewVBox(
		widget.NewLabel("One at a time, one always picked (Required):"), single, singleStatus,
		widget.NewSeparator(),
		widget.NewLabel("Any number (Multiple):"), multi, multiStatus,
		widget.NewSeparator(),
		widget.NewLabel("Three states (Multiple + Excludable): tap a chip to move it off → included → excluded → off"),
		three, threeStatus,
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel("Rows aligned at"), align),
	))
	chipCard.Style = widget.CardStyleOutlined

	// --- Check groups ---
	vertical := widget.NewCheckGroup([]string{"Email", "SMS", "Push notifications"}, nil)
	vertical.SetSelected([]string{"Email"})
	horizontal := widget.NewCheckGroup([]string{"Mon", "Tue", "Wed", "Thu", "Fri"}, nil)
	horizontal.Horizontal = true
	horizontal.SetSelected([]string{"Mon", "Wed", "Fri"})
	required := widget.NewCheckGroup([]string{"Terms of service", "Privacy policy"}, nil)
	required.Required = true
	required.SetSelected([]string{"Terms of service"})
	checkStatus := widget.NewLabel("Notify by: Email")
	vertical.OnChanged = func(sel []string) { checkStatus.SetText("Notify by: " + joinOrNone(sel)) }
	disableables = append(disableables, vertical, horizontal, required)

	checkCard := widget.NewCard("Check Group", "Checkboxes for a set of options", container.NewVBox(
		widget.NewLabel("Vertical:"), vertical, checkStatus,
		widget.NewLabel("Horizontal:"), horizontal,
		widget.NewLabel("Required (the last one can't be unticked):"), required,
	))
	checkCard.Style = widget.CardStyleOutlined

	// --- Spin boxes ---
	count := widget.NewSpinBox(0, 10, 1)
	count.SetValue(3)
	countStatus := widget.NewLabel("Guests: 3")
	count.OnChanged = func(v float64) { countStatus.SetText(fmt.Sprintf("Guests: %.0f", v)) }
	speed := widget.NewSpinBox(0, 5, 0.25)
	speed.SetValue(1.5)
	speedStatus := widget.NewLabel("Speed: 1.50×")
	speed.OnChanged = func(v float64) { speedStatus.SetText(fmt.Sprintf("Speed: %.2f×", v)) }
	disableables = append(disableables, count, speed)

	spinCard := widget.NewCard("Spin Box", "A number with step buttons; arrow keys step it too", container.NewVBox(
		widget.NewLabel("Whole numbers, 0 to 10:"), container.NewHBox(count, countStatus),
		widget.NewLabel("Steps of 0.25, 0 to 5:"), container.NewHBox(speed, speedStatus),
	))
	spinCard.Style = widget.CardStyleOutlined

	disable := widget.NewCheck("Disable everything on this page", func(on bool) {
		for _, d := range disableables {
			if on {
				d.Disable()
			} else {
				d.Enable()
			}
		}
	})

	return container.NewVBox(disable, chipCard, checkCard, spinCard)
}

func joinOrNone(s []string) string {
	if len(s) == 0 {
		return "none"
	}
	return strings.Join(s, ", ")
}

func chipStateName(s widget.ChipState) string {
	switch s {
	case widget.ChipStateSelected:
		return "included"
	case widget.ChipStateExcluded:
		return "excluded"
	}
	return "off"
}
