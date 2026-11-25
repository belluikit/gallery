//go:build !android && !js

package pages

import (
	"bytes"
	"fmt"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/mmedia"
)

// VideoPlayerPage shows x/mmedia's VideoPlayer, the bare video surface
// (gallery#47). Unlike MediaPlayer it has no transport controls of its own, so
// the page makes them: Play, Pause and Stop buttons and a clock fed by
// WithOnTimeUpdate. It reads the bundled sample from memory, the io.Reader form,
// where the Media Player page opens it by path.
func VideoPlayerPage(_ bellui.Window) core.CanvasObject {
	clock := widget.NewLabel("0:00 / 0:00")
	vp, err := mmedia.NewVideoPlayer(bytes.NewReader(sampleVideo),
		mmedia.WithAutoPlay(), mmedia.WithLoop(),
		mmedia.WithOnTimeUpdate(func(current, total time.Duration) {
			bellui.Do(func() { clock.SetText(clockText(current) + " / " + clockText(total)) })
		}))
	if err != nil {
		return widget.NewLabel("Cannot play the sample video: " + err.Error())
	}
	OnLeave(vp.Close)

	play := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), vp.Play)
	pause := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), vp.Pause)
	stop := widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), vp.Stop)
	controls := container.NewHBox(play, pause, stop, clock)

	return container.NewBorder(
		note("VideoPlayer draws only the video. These buttons and the clock are the page's own, calling Play, Pause and Stop and listening with WithOnTimeUpdate."),
		controls, nil, nil, vp)
}

// clockText is d as minutes and seconds.
func clockText(d time.Duration) string {
	s := int(d.Round(time.Second) / time.Second)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
