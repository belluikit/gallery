//go:build !android && !js

package pages

import (
	"bytes"
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

// sampleSubtitles are three WebVTT cues timed to the sample video, shown
// through mmedia.SubtitleOverlay (gallery#47).
//
//go:embed media/sample.vtt
var sampleSubtitles []byte

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
// decoded by FFmpeg through goffmpeg. It starts on a bundled test video with
// subtitles from a bundled WebVTT track; "Open file…" plays any local file
// FFmpeg can read, without subtitles.
func MediaPlayerPage(w bellui.Window) core.CanvasObject {
	slot := container.NewStack()
	status := widget.NewLabel("")
	var current *mmedia.MediaPlayer

	subtitles := mmedia.NewSubtitleOverlay()
	if err := subtitles.LoadWebVTT(bytes.NewReader(sampleSubtitles)); err != nil {
		status.SetText("Subtitles unavailable: " + err.Error())
	}
	showSubtitles := widget.NewCheck("Subtitles", func(on bool) {
		if on {
			subtitles.Show()
		} else {
			subtitles.Hide()
		}
	})
	showSubtitles.SetChecked(true)

	// play opens path. withSubtitles lays the overlay over the video, with
	// its position following playback; only the sample has a track.
	play := func(path, name string, withSubtitles bool) {
		if current != nil {
			current.Close()
			current = nil
		}
		mp, err := mmedia.NewMediaPlayerURL(path, mmedia.WithAutoPlay(), mmedia.WithLoop())
		if err != nil {
			Show(slot, widget.NewLabel("Cannot play "+name+": "+err.Error()))
			return
		}
		current = mp
		if withSubtitles {
			Show(slot, mmedia.MediaPlayerWithSubtitles(mp, subtitles))
			showSubtitles.Enable()
		} else {
			Show(slot, mp)
			showSubtitles.Disable()
		}
		status.SetText("Playing " + name)
	}
	OnLeave(func() {
		if current != nil {
			current.Close()
			current = nil
		}
	})

	if p, err := sampleVideoPath(); err == nil {
		play(p, "the sample video", true)
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
			play(uri.Path(), uri.Name(), false)
		}, w)
	})

	return container.NewBorder(container.NewHBox(open, showSubtitles, status), nil, nil, nil, slot)
}
