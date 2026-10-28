//go:build !js

package pages

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"

	"github.com/belluikit/bellui/x/mmedia"
)

// Every bundled clip plays through x/mmedia's audio player to its end: opened
// from memory, as the page does, sought to a second before the end, and the
// end reported. A clip that would not decode, or a format the bundled FFmpeg
// lacks, fails here instead of on a device.
func TestAudioClipsPlayToTheirEnd(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	for _, f := range audioFormats {
		t.Run(f.label, func(t *testing.T) {
			data, err := audioClips.ReadFile("media/audio/" + f.file)
			if err != nil {
				t.Fatal(err)
			}
			ended := make(chan struct{}, 1)
			ap, err := mmedia.NewAudioPlayer(bytes.NewReader(data),
				mmedia.WithOnEnd(func() { ended <- struct{}{} }))
			if err != nil {
				if strings.Contains(err.Error(), "audio device") {
					t.Skipf("no audio device: %v", err)
				}
				t.Fatal(err)
			}
			defer ap.Close()
			if d := ap.Duration(); d < 11*time.Second || d > 13*time.Second {
				t.Fatalf("duration %v, want the clip's 12 s", d)
			}
			ap.Play()
			ap.Seek(11 * time.Second)
			select {
			case <-ended:
			case <-time.After(5 * time.Second):
				t.Fatalf("no end within 5 s of seeking to 11 s (at %v)", ap.CurrentTime())
			}
		})
	}
}

// The page builds on its first clip, and leaving it closes the player.
func TestAudioPlayerPageBuildsAndLeaves(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(900, 700))

	w.SetContent(AudioPlayerPage(w))
	LeavePage()
	LeavePage()
}
