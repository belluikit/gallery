//go:build !js

package pages

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/container"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/dialog"
	"github.com/belluikit/bellui/theme"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/mmedia"
)

// audioClips are a 12-second rising whole-tone scale from A3, one note per
// second, made for this gallery (no third-party audio), in each format the
// page offers. A note per second makes a seek audible: after a seek to 0:05
// the sixth note plays.
//
//go:embed media/audio/scale.*
var audioClips embed.FS

// audioFormats lists the bundled clips in the order the page offers them.
var audioFormats = []struct{ label, file string }{
	{"MP3", "scale.mp3"},
	{"AAC (m4a)", "scale.m4a"},
	{"Ogg Vorbis", "scale.ogg"},
	{"Opus", "scale.opus"},
	{"FLAC", "scale.flac"},
	{"WAV", "scale.wav"},
}

// AudioPages returns the audio-only player, which runs wherever FFmpeg links:
// desktop, iOS and Android. The web build gets the stub in audio_pages_js.go.
func AudioPages() []PageDef {
	return []PageDef{{Name: "Audio Player", Icon: theme.VolumeUpIcon(), New: AudioPlayerPage}}
}

// AudioPlayerPage demos x/mmedia's audio-only player: the bundled clip in each
// format, any local file, and network streams by URL, with loop and auto-play
// switches and a log of what the player reports.
func AudioPlayerPage(w bellui.Window) core.CanvasObject {
	slot := container.NewStack()
	position := widget.NewLabel("")
	logLines := []string{}
	events := widget.NewLabel("")
	events.Wrapping = bellui.TextWrapWord
	logf := func(format string, args ...any) {
		line := time.Now().Format("15:04:05") + "  " + fmt.Sprintf(format, args...)
		logLines = append([]string{line}, logLines...)
		if len(logLines) > 8 {
			logLines = logLines[:8]
		}
		events.SetText(strings.Join(logLines, "\n"))
	}

	loop := widget.NewCheck("Loop", nil)
	autoPlay := widget.NewCheck("Auto-play", nil)
	autoPlay.SetChecked(true)

	var current *mmedia.AudioPlayer
	var source io.Closer // a reader the current player reads through, if any
	closeCurrent := func() {
		if current != nil {
			current.Close()
			current = nil
		}
		if source != nil {
			_ = source.Close()
			source = nil
		}
	}
	loop.OnChanged = func(on bool) {
		if current != nil {
			current.SetLoop(on)
		}
	}
	OnLeave(closeCurrent)

	// options are rebuilt for each source so the switches apply to it.
	options := func() []mmedia.Option {
		opts := []mmedia.Option{
			mmedia.WithOnEnd(func() { bellui.Do(func() { logf("ended") }) }),
			mmedia.WithOnError(func(err error) { bellui.Do(func() { logf("error: %v", err) }) }),
			mmedia.WithOnTimeUpdate(func(cur, total time.Duration) {
				bellui.Do(func() { position.SetText(formatClock(cur) + " / " + formatClock(total)) })
			}),
		}
		if loop.Checked {
			opts = append(opts, mmedia.WithLoop())
		}
		if autoPlay.Checked {
			opts = append(opts, mmedia.WithAutoPlay())
		}
		return opts
	}
	show := func(name string, open func([]mmedia.Option) (*mmedia.AudioPlayer, error)) {
		closeCurrent()
		position.SetText("")
		ap, err := open(options())
		if err != nil {
			slot.Objects = []core.CanvasObject{widget.NewLabel("Cannot play " + name + ": " + err.Error())}
			slot.Refresh()
			logf("cannot open %s: %v", name, err)
			return
		}
		current = ap
		slot.Objects = []core.CanvasObject{ap}
		slot.Refresh()
		logf("opened %s (%s)", name, formatClock(ap.Duration()))
	}
	playClip := func(i int) {
		f := audioFormats[i]
		data, err := audioClips.ReadFile("media/audio/" + f.file)
		if err != nil {
			logf("cannot read %s: %v", f.file, err)
			return
		}
		show(f.label+" clip", func(opts []mmedia.Option) (*mmedia.AudioPlayer, error) {
			return mmedia.NewAudioPlayer(bytes.NewReader(data), opts...)
		})
	}

	labels := make([]string, len(audioFormats))
	for i, f := range audioFormats {
		labels[i] = f.label
	}
	format := widget.NewSelect(labels, func(label string) {
		for i, f := range audioFormats {
			if f.label == label {
				playClip(i)
				return
			}
		}
	})

	url := widget.NewEntry()
	url.SetPlaceHolder("http://… or https://… (a file or an HLS playlist)")
	playURL := widget.NewButtonWithIcon("Play URL", theme.MediaPlayIcon(), func() {
		u := strings.TrimSpace(url.Text)
		if u == "" {
			return
		}
		show(u, func(opts []mmedia.Option) (*mmedia.AudioPlayer, error) {
			return mmedia.NewAudioPlayerURL(u, opts...)
		})
	})
	url.OnSubmitted = func(string) { playURL.OnTapped() }

	openFile := widget.NewButtonWithIcon("Open file…", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(r bellui.URIReadCloser, err error) {
			if err != nil {
				logf("cannot open the file: %v", err)
				return
			}
			if r == nil {
				logf("no file chosen")
				return
			}
			// Read it through the dialog's reader, which seeks on every
			// platform, and close it with the player: on iOS the reader holds
			// the access a file picked from outside the app needs.
			show(r.URI().Name(), func(opts []mmedia.Option) (*mmedia.AudioPlayer, error) {
				return mmedia.NewAudioPlayer(r, opts...)
			})
			if current != nil {
				source = r
			} else {
				_ = r.Close()
			}
		}, w)
	})

	controls := container.NewVBox(
		widget.NewForm(widget.NewFormItem("Bundled clip", format)),
		container.NewBorder(nil, nil, nil, playURL, url),
		container.NewHBox(openFile, loop, autoPlay),
	)
	log := widget.NewCard("Player events", "Newest first", events)

	format.SetSelected(audioFormats[0].label)
	return container.NewVBox(controls, slot, position, log)
}

// formatClock renders d as m:ss.
func formatClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	s := int(d.Round(time.Second) / time.Second)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
