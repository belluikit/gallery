// Ported from the former uploader cmd/demo app.
package pages

import (
	"context"
	"fmt"
	"sync"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/data/binding"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/widget"

	"github.com/trendvidia/uploader/sdk"
)

// UploaderPage demonstrates file upload/download with progress using the
// uploader SDK. It needs an uploader server to transfer files against
// (default localhost:9090).
func UploaderPage(w bellui.Window) core.CanvasObject {
	// State: connected client and last uploaded file ID for download.
	var (
		mu         sync.Mutex
		client     *sdk.Client
		lastFileID string
	)

	// Use data bindings to avoid the Bellui textRenderer stale-bounds panic
	// that occurs when SetText changes text length during a layout pass.

	statusText := binding.NewString()
	statusText.Set("Not connected — start an uploader server and click Connect")
	statusLabel := widget.NewLabelWithData(statusText)
	statusLabel.Truncation = bellui.TextTruncateEllipsis

	fileIDText := binding.NewString()
	fileIDText.Set("File ID: (none)")
	fileIDLabel := widget.NewLabelWithData(fileIDText)
	fileIDLabel.Truncation = bellui.TextTruncateEllipsis

	progressVal := binding.NewFloat()
	progressBar := widget.NewProgressBarWithData(progressVal)

	serverEntry := widget.NewEntry()
	serverEntry.SetText("localhost:9090")
	serverEntry.SetPlaceHolder("server:port")

	connectBtn := widget.NewButton("Connect", func() {
		mu.Lock()
		if client != nil {
			client.Close()
			client = nil
		}
		mu.Unlock()
		c, err := sdk.New(serverEntry.Text)
		if err != nil {
			statusText.Set(fmt.Sprintf("Connection failed: %v", err))
			return
		}
		mu.Lock()
		client = c
		mu.Unlock()
		statusText.Set("Connected to " + serverEntry.Text)
	})

	// --- Upload ---

	uploadBtn := widget.NewButton("Upload File...", func() {
		mu.Lock()
		c := client
		mu.Unlock()
		if c == nil {
			statusText.Set("Not connected")
			return
		}
		dialog.ShowFileOpen(func(reader bellui.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return // cancelled
			}
			filePath := reader.URI().Path()
			fileName := reader.URI().Name()
			reader.Close()

			statusText.Set("Uploading: " + fileName)
			progressVal.Set(0)

			go func() {
				result, err := c.Upload(
					context.Background(), filePath,
					sdk.WithProgress(func(current, total int64) {
						if total > 0 {
							progressVal.Set(float64(current) / float64(total))
						}
					}),
				)
				if err != nil {
					bellui.Do(func() {
						statusText.Set(fmt.Sprintf("Upload failed: %v", err))
					})
					return
				}

				mu.Lock()
				lastFileID = result.FileID
				mu.Unlock()

				bellui.Do(func() {
					fileIDText.Set("File ID: " + result.FileID)
					if result.Complete {
						statusText.Set(fmt.Sprintf("Upload complete (%d bytes)", result.CommittedSize))
					} else {
						statusText.Set(fmt.Sprintf("Upload partial (%d bytes, session: %s)", result.CommittedSize, result.SessionID))
					}
				})
			}()
		}, w)
	})

	// --- Download ---

	downloadBtn := widget.NewButton("Download File...", func() {
		mu.Lock()
		c := client
		fid := lastFileID
		mu.Unlock()
		if c == nil {
			statusText.Set("Not connected")
			return
		}
		if fid == "" {
			statusText.Set("No file to download — upload one first")
			return
		}

		dialog.ShowFileSave(func(writer bellui.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				return // cancelled
			}
			destPath := writer.URI().Path()
			destName := writer.URI().Name()
			writer.Close()

			statusText.Set("Downloading to: " + destName)
			progressVal.Set(0)

			go func() {
				err := c.Download(
					context.Background(), fid, destPath,
					sdk.WithProgress(func(current, total int64) {
						if total > 0 {
							progressVal.Set(float64(current) / float64(total))
						}
					}),
				)
				if err != nil {
					bellui.Do(func() {
						statusText.Set(fmt.Sprintf("Download failed: %v", err))
					})
					return
				}
				bellui.Do(func() {
					statusText.Set("Download complete: " + destPath)
				})
			}()
		}, w)
	})

	// --- Layout ---

	serverRow := container.NewBorder(nil, nil, widget.NewLabel("Server:"), connectBtn, serverEntry)

	buttons := container.NewHBox(uploadBtn, downloadBtn)

	return container.NewVBox(
		serverRow,
		widget.NewSeparator(),
		buttons,
		fileIDLabel,
		progressBar,
		statusLabel,
	)
}
