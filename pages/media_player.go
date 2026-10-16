//go:build !android && !js

package pages

import (
	_ "embed"
	"os"
	"path/filepath"
	"sync"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/mmedia"
)

// sampleVideo is a 6-second 640x360 test pattern with a silent audio track,
// made for this gallery (no third-party footage).
//
//go:embed media/sample.mp4
var sampleVideo []byte

var (
	samplePathOnce sync.Once
	samplePath     string
	samplePathErr  error
)

// sampleVideoPath writes the embedded video to a temporary file once: the
// player opens media by path or URL.
func sampleVideoPath() (string, error) {
	samplePathOnce.Do(func() {
		dir, err := os.MkdirTemp("", "bellui-gallery-")
		if err != nil {
			samplePathErr = err
			return
		}
		samplePath = filepath.Join(dir, "sample.mp4")
		samplePathErr = os.WriteFile(samplePath, sampleVideo, 0o600)
	})
	return samplePath, samplePathErr
}

// MediaPlayerPage demos x/mmedia: a video player with transport controls,
// decoded by FFmpeg through goffmpeg. It starts on a bundled test video; "Open
// file…" plays any local file FFmpeg can read.
func MediaPlayerPage(w bellui.Window) core.CanvasObject {
	slot := container.NewStack()
	status := widget.NewLabel("")
	var current *mmedia.MediaPlayer

	play := func(path, name string) {
		if current != nil {
			current.Close()
			current = nil
		}
		mp, err := mmedia.NewMediaPlayerURL(path, mmedia.WithAutoPlay(), mmedia.WithLoop())
		if err != nil {
			slot.Objects = []core.CanvasObject{widget.NewLabel("Cannot play " + name + ": " + err.Error())}
			slot.Refresh()
			return
		}
		current = mp
		slot.Objects = []core.CanvasObject{mp}
		slot.Refresh()
		status.SetText("Playing " + name)
	}
	OnLeave(func() {
		if current != nil {
			current.Close()
			current = nil
		}
	})

	if p, err := sampleVideoPath(); err == nil {
		play(p, "the sample video")
	} else {
		status.SetText("Sample video unavailable: " + err.Error())
	}

	open := widget.NewButtonWithIcon("Open file…", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(r bellui.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			uri := r.URI()
			_ = r.Close()
			play(uri.Path(), uri.Name())
		}, w)
	})

	return container.NewBorder(container.NewHBox(open, status), nil, nil, nil, slot)
}
