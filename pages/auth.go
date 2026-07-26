package pages

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/auth"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/widget"
)

// ShowAuthDialog opens a dialog with branded sign-in buttons for all
// supported OIDC/OAuth providers.
func ShowAuthDialog(w bellui.Window) {
	statusLabel := widget.NewLabel("")

	onAuth := func(result auth.OIDCResult) {
		statusLabel.SetText("Welcome, " + result.Name + " (" + result.Email + ")")
	}

	google := widget.NewAuthButton("google", "demo-client-id", onAuth)
	microsoft := widget.NewAuthButton("microsoft", "demo-client-id", onAuth)
	apple := widget.NewAuthButton("apple", "demo-client-id", onAuth)
	github := widget.NewAuthButton("github", "demo-client-id", onAuth)
	okta := widget.NewAuthButton("okta", "demo-client-id", onAuth)
	okta.Issuer = "https://dev-example.okta.com"

	custom := widget.NewAuthButton("custom", "demo-client-id", onAuth)
	custom.Text = "Sign in with Corporate SSO"
	custom.Issuer = "https://auth.example.com"

	var dialog *widget.PopUp
	closeBtn := &widget.Button{Text: "Cancel", Style: widget.ButtonStyleText, OnTapped: func() {
		dialog.Hide()
	}}

	content := container.NewVBox(
		widget.NewLabel("Choose a sign-in provider:"),
		google,
		microsoft,
		apple,
		github,
		okta,
		widget.NewSeparator(),
		custom,
		statusLabel,
		widget.NewSeparator(),
		container.NewHBox(closeBtn),
	)

	dialog = widget.NewModalPopUp(
		container.NewPadded(content),
		w.Canvas(),
	)
	dialog.Resize(geo.NewSize(360, 500))
	dialog.Show()
}

// AuthPage is kept for direct page navigation (optional).
func AuthPage(w bellui.Window) core.CanvasObject {
	btn := widget.NewButton("Open Auth Dialog", func() {
		ShowAuthDialog(w)
	})
	btn.Style = widget.ButtonStyleFilled

	card := widget.NewCard("Authentication", "OIDC/OAuth sign-in with branded provider buttons", container.NewVBox(
		widget.NewLabel("Click the account icon in the top bar, or:"),
		btn,
	))
	card.Style = widget.CardStyleOutlined

	return container.NewVBox(card)
}
