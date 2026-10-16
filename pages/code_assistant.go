//go:build !android && !js

// Excluded from Android builds: the assistant uses clair's knowledge base,
// which links FAISS (go-faiss) for vector search — a host-only native
// library. Omitted on mobile via pages/native_mobile.go.
package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/clair/assistant"
	"github.com/belluikit/clair/llm"
	"github.com/belluikit/clair/llm/anthropic"
)

const (
	apiKeyPref = "clair_anthropic_api_key"
	modelPref  = "clair_anthropic_model"
)

var availableModels = []string{
	"claude-sonnet-4-20250514",
	"claude-haiku-4-20250414",
	"claude-opus-4-20250514",
	"claude-3-7-sonnet-20250219",
	"claude-3-5-haiku-20241022",
}

// belluiSystemPrompt gives the LLM context about the Bellui framework.
const belluiSystemPrompt = `You are a code assistant for the Bellui GUI toolkit (github.com/belluikit/bellui), a cross-platform Go framework.

Key patterns:
- App: app.New() or app.NewWithID("id") creates the application
- Window: a.NewWindow("title") creates a window, w.SetContent(obj) sets the root widget
- Widgets: widget.NewButton, widget.NewLabel, widget.NewEntry, widget.NewList, etc.
- Containers: container.NewVBox, NewHBox, NewBorder, NewHSplit(leading, trailing), NewGridWithColumns
- Dialogs: dialog.ShowInformation, dialog.ShowConfirm — these are in the dialog package, not widget
- List widget: widget.NewList(length func() int, createItem func() core.CanvasObject, updateItem func(int, core.CanvasObject))
- Entry: Set MultiLine=true and PlaceHolder for multi-line text input
- Size: geo.NewSize(w, h) is the constructor

When the user asks about code, reference the relevant Bellui API documentation provided in context.
Give concise, practical answers with corrected code snippets. Use Go code blocks.`

// Broken sample code for the editor — has several issues a developer might ask about.
const sampleCode = `package main

import (
	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/app"
	"github.com/belluikit/bellui/widget"
	"github.com/belluikit/bellui/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("My App")

	// BUG: Entry is missing PlaceHolder and MultiLine config
	input := widget.NewEntry()

	// BUG: List requires 3 callbacks, not just data
	items := []string{"Apple", "Banana", "Cherry"}
	list := widget.NewList(items)

	// BUG: wrong function — ShowInformation doesn't exist on widget
	submit := widget.NewButton("Submit", func() {
		widget.ShowInformation("Result", input.Text, w)
	})

	// BUG: Split expects exactly 2 objects, not a container
	split := container.NewHSplit(
		container.NewVBox(input, submit),
	)

	w.SetContent(split)
	w.Resize(geo.Size{400, 300})
	w.ShowAndRun()
}
`

// CodeAssistantPage shows a code editor with a Claude-powered AI assistant.
func CodeAssistantPage(w bellui.Window) core.CanvasObject {
	a := bellui.CurrentApp()

	// Code editor.
	editor := widget.NewEntry()
	editor.MultiLine = true
	editor.TextStyle = bellui.TextStyle{Monospace: true}
	editor.SetText(sampleCode)
	editor.SetMinRowsVisible(20)

	// Chat history — dynamic VBox, adding labels per message.
	chatBox := container.NewVBox(
		widget.NewLabel("Chat messages will appear here."),
		widget.NewSeparator(),
	)

	appendChat := func(role, text string) {
		header := widget.NewLabel("— " + role + " —")
		header.TextStyle = bellui.TextStyle{Bold: true}
		body := widget.NewLabel(text)
		body.Wrapping = bellui.TextWrapWord
		chatBox.Add(header)
		chatBox.Add(body)
		chatBox.Add(widget.NewSeparator())
	}

	// Prompt input.
	promptEntry := widget.NewEntry()
	promptEntry.PlaceHolder = "Ask about the code (e.g. \"fix the bugs\", \"explain the List widget\")"

	// Status.
	statusLabel := widget.NewLabel("Configure API key to enable the assistant.")
	statusLabel.TextStyle = bellui.TextStyle{Italic: true}
	statusLabel.Wrapping = bellui.TextWrapWord

	// Get or create the assistant with the existing Clair KB.
	assist, err := getClairAssistant()
	if err != nil {
		errCard := widget.NewCard("Code Assistant", "Knowledge Base Error",
			widget.NewLabel(fmt.Sprintf("Failed to initialize: %v", err)))
		errCard.Style = widget.CardStyleOutlined
		return container.NewVBox(errCard)
	}

	// Check if API key is already stored.
	apiKey := bellui.PrefOr(a.Preferences(), apiKeyPref, "")
	if apiKey != "" {
		model := bellui.PrefOr(a.Preferences(), modelPref, "")
		if model == "" {
			model = availableModels[0]
		}
		configureProvider(assist, apiKey, model)
		statusLabel.SetText(fmt.Sprintf("Connected to %s.", model))
	}

	// API key configuration button.
	configBtn := &widget.Button{
		Text:  "API Key",
		Icon:  theme.SettingsIcon(),
		Style: widget.ButtonStyleOutlined,
		OnTapped: func() {
			showAPIKeyDialog(w, a, assist, statusLabel)
		},
	}

	// Send handler.
	var sending bool
	doSend := func() {
		query := promptEntry.Text
		if query == "" || sending {
			return
		}

		if !assist.HasAI() {
			showAPIKeyDialog(w, a, assist, statusLabel)
			return
		}

		sending = true
		promptEntry.Disable()
		statusLabel.SetText("Thinking...")

		// Include current editor content as context.
		editorCtx := assistant.EditorContext{
			CurrentFile: "main.go",
			VisibleCode: editor.Text,
			Language:    "go",
		}

		// Prepend the code to the query so Claude can see it.
		fullQuery := query + "\n\nHere is the code I'm working on:\n```go\n" + editor.Text + "\n```"

		appendChat("you", query)
		promptEntry.SetText("")

		go func() {
			stream, err := assist.AI.Ask(context.Background(), fullQuery, editorCtx)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Error: %v", err))
				sending = false
				promptEntry.Enable()
				return
			}

			var response strings.Builder
			for {
				event, err := stream.Next()
				if err != nil {
					break
				}
				if event.Type == llm.EventTextDelta {
					response.WriteString(event.Text)
				} else if event.Type == llm.EventDone || event.Type == llm.EventError {
					break
				}
			}
			stream.Close()

			if finalText := response.String(); finalText != "" {
				appendChat("assistant", finalText)
			}
			statusLabel.SetText("Ready.")
			sending = false
			promptEntry.Enable()
		}()
	}

	promptEntry.OnSubmitted = func(_ string) { doSend() }

	sendBtn := &widget.Button{
		Text:     "Send",
		Icon:     theme.MailSendIcon(),
		Style:    widget.ButtonStyleFilled,
		OnTapped: doSend,
	}

	// Layout.
	promptRow := container.NewBorder(nil, nil, nil, container.NewHBox(configBtn, sendBtn), promptEntry)

	editorCard := widget.NewCard("Code Editor", "Sample Go code with several bugs", editor)
	editorCard.Style = widget.CardStyleOutlined

	chatCard := widget.NewCard(
		"AI Assistant", "Powered by Claude via Clair",
		container.NewVBox(chatBox, statusLabel, promptRow),
	)
	chatCard.Style = widget.CardStyleOutlined

	return container.NewGridWithColumns(2, editorCard, chatCard)
}

func showAPIKeyDialog(w bellui.Window, a bellui.App, assist *assistant.Assistant, status *widget.Label) {
	keyEntry := widget.NewPasswordEntry()
	keyEntry.PlaceHolder = "sk-ant-..."

	existing := bellui.PrefOr(a.Preferences(), apiKeyPref, "")
	if existing != "" {
		keyEntry.SetText(existing)
	}

	// Model selector — starts with fallback list, loads from API on demand.
	modelSelect := widget.NewSelect(availableModels, nil)
	savedModel := bellui.PrefOr(a.Preferences(), modelPref, "")
	if savedModel != "" {
		modelSelect.SetSelected(savedModel)
	} else {
		modelSelect.SetSelected(availableModels[0])
	}

	keyLabel := widget.NewLabel("API Key:")
	keyLabel.TextStyle = bellui.TextStyle{Bold: true}
	modelLabel := widget.NewLabel("Model:")
	modelLabel.TextStyle = bellui.TextStyle{Bold: true}

	testStatus := widget.NewLabel("Enter your API key, then load available models.")
	testStatus.Wrapping = bellui.TextWrapWord

	var dlg *dialog.CustomDialog

	loadModelsBtn := widget.NewButton("Load Models", func() {
		key := keyEntry.Text
		if key == "" {
			testStatus.SetText("Enter a key first.")
			return
		}
		testStatus.SetText("Loading models...")

		go func() {
			models, err := listModels(key)
			log.Printf("Loaded %d models, err: %v", len(models), err)
			if err != nil {
				testStatus.SetText(fmt.Sprintf("Failed to load models: %v", err))
				return
			}
			prev := modelSelect.Selected
			modelSelect.Options = models
			found := false
			for _, m := range models {
				if m == prev {
					found = true
					break
				}
			}
			if found {
				modelSelect.SetSelected(prev)
			} else if len(models) > 0 {
				modelSelect.SetSelected(models[0])
			}
			modelSelect.Refresh()
			testStatus.SetText(fmt.Sprintf("Loaded %d models.", len(models)))
		}()
	})

	testBtn := widget.NewButton("Validate", func() {
		key := keyEntry.Text
		if key == "" {
			testStatus.SetText("Enter a key first.")
			return
		}
		model := modelSelect.Selected
		if model == "" {
			testStatus.SetText("Select a model first.")
			return
		}
		testStatus.SetText(fmt.Sprintf("Testing %s...", model))

		go func() {
			err := validateAPIKey(key, model)
			log.Printf("Validation result: %v", err)
			msg := fmt.Sprintf("Valid! %s works.", model)
			if err != nil {
				msg = fmt.Sprintf("Failed: %v", err)
			}
			testStatus.SetText(msg)
		}()
	})
	testBtn.Importance = widget.MediumImportance

	saveBtn := widget.NewButton("Save", func() {
		key := keyEntry.Text
		if key == "" {
			return
		}
		model := modelSelect.Selected
		if model == "" {
			model = availableModels[0]
		}
		bellui.SetPref(a.Preferences(), apiKeyPref, key)
		bellui.SetPref(a.Preferences(), modelPref, model)
		configureProvider(assist, key, model)
		status.SetText(fmt.Sprintf("Connected to %s.", model))
		if dlg != nil {
			dlg.Hide()
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancel", func() {
		if dlg != nil {
			dlg.Hide()
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Configure the Anthropic API connection."),
		keyLabel,
		keyEntry,
		modelLabel,
		container.NewBorder(nil, nil, nil, loadModelsBtn, modelSelect),
		testStatus,
		container.NewHBox(testBtn, saveBtn, cancelBtn),
	)

	dlg = dialog.NewCustomWithoutButtons("Configure Claude", content, w)
	dlg.Resize(geo.NewSize(500, 320))
	dlg.Show()
}

// validateAPIKey sends a minimal request to verify the key and model work.
func validateAPIKey(apiKey, model string) error {
	provider := anthropic.New(anthropic.Config{APIKey: apiKey, Model: model})
	defer provider.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := provider.Complete(ctx, llm.Request{
		Messages:  []llm.Message{{Role: llm.RoleUser, Content: "Hi"}},
		MaxTokens: 1,
	})
	return err
}

func configureProvider(assist *assistant.Assistant, apiKey, model string) {
	provider := anthropic.New(anthropic.Config{
		APIKey: apiKey,
		Model:  model,
	})
	assist.AI.SetProvider(provider)
}

// listModels fetches available models from the Anthropic API.
func listModels(apiKey string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.anthropic.com/v1/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}

	// Sort so newest models come first.
	sort.Sort(sort.Reverse(sort.StringSlice(models)))
	return models, nil
}
