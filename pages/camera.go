//go:build !android && !js

package pages

import (
	"sync"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/mmedia"
)

// capturePreview is a live device preview the Camera page shows and closes.
type capturePreview interface {
	core.CanvasObject
	Close()
}

// openCapture opens a capture device as a player. A variable so tests can
// stand in for a camera: opening a real one lights it and, on macOS, asks the
// user for permission.
var openCapture = func(device string) (capturePreview, error) {
	return mmedia.NewCapturePlayer(device, mmedia.WithAutoPlay())
}

// CameraPage shows x/mmedia's CaptureSelect, a camera picker, and a live
// preview of the camera chosen, through NewCapturePlayer (gallery#47).
//
// Building the page opens no device, so it asks for no camera permission.
// Where the platform lists its cameras (Linux, Windows) the picker fills at
// once; on macOS, which cannot list them, it offers a Scan button, and only
// pressing it opens each camera in turn. Choosing a camera opens it.
func CameraPage(_ bellui.Window) core.CanvasObject {
	slot := container.NewStack(widget.NewLabel("No camera open"))
	status := widget.NewLabel("")
	status.Wrapping = bellui.TextWrapWord

	var mu sync.Mutex
	var current capturePreview
	left := false
	// replace makes p the page's preview and closes the one before it. After
	// the page is left it closes p instead: an open that was still under way
	// must not leave a camera running.
	replace := func(p capturePreview) bool {
		mu.Lock()
		old := current
		keep := !left
		if keep {
			current = p
		} else {
			current = nil
		}
		mu.Unlock()
		if old != nil {
			old.Close()
		}
		if !keep && p != nil {
			p.Close()
		}
		return keep
	}
	OnLeave(func() {
		mu.Lock()
		left = true
		mu.Unlock()
		replace(nil)
	})

	picker := mmedia.NewCaptureSelect(func(d mmedia.CaptureDevice) {
		status.SetText("Opening " + d.Label() + "…")
		// Opening a camera takes about a second: off the main goroutine.
		go func() {
			p, err := openCapture(d.Name)
			bellui.Do(func() {
				if err != nil {
					status.SetText("Cannot open " + d.Label() + ": " + err.Error())
					return
				}
				if replace(p) {
					Show(slot, p)
					status.SetText("Showing " + d.Label())
				}
			})
		}()
	})
	if err := picker.Err(); err != nil {
		status.SetText(err.Error())
	}

	return container.NewBorder(
		container.NewVBox(
			note("Pick a camera to see it live. On macOS, which cannot list its cameras, Scan finds them by opening each one, so it turns the camera light on and asks for permission."),
			picker, status,
		),
		nil, nil, nil, slot)
}
