//go:build !android && !js

package pages

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/belluikit/bellui"
	"github.com/belluikit/bellui/canvas"
	"github.com/belluikit/bellui/core"
	"github.com/belluikit/bellui/geo"
	"github.com/belluikit/bellui/test"
	"github.com/belluikit/bellui/widget"

	"github.com/belluikit/bellui/x/mmedia"
)

// walkTree calls f for every object under o, through containers and widget
// renderers.
func walkTree(o core.CanvasObject, f func(core.CanvasObject)) {
	f(o)
	switch v := o.(type) {
	case *bellui.Container:
		for _, c := range v.Objects {
			walkTree(c, f)
		}
	case core.Widget:
		for _, c := range test.WidgetRenderer(v).Objects() {
			walkTree(c, f)
		}
	}
}

// eventually polls cond on the main goroutine for up to five seconds.
func eventually(t *testing.T, what string, cond func() bool) {
	t.Helper()
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(20 * time.Millisecond) {
		ok := false
		bellui.DoAndWait(func() { ok = cond() })
		if ok {
			return
		}
	}
	t.Fatal(what)
}

// The Media Player page lays a SubtitleOverlay over the sample video, loaded
// with the gallery's WebVTT track: at a cue's time the overlay shows that
// cue, between cues nothing, and the Subtitles box hides it (gallery#47).
func TestMediaPlayerSubtitles(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(900, 600))
	page := MediaPlayerPage(w)
	w.SetContent(page)
	defer LeavePage()

	var overlay *mmedia.SubtitleOverlay
	var box *widget.Check
	walkTree(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *mmedia.SubtitleOverlay:
			overlay = v
		case *widget.Check:
			if v.Text == "Subtitles" {
				box = v
			}
		}
	})
	if overlay == nil || box == nil {
		t.Fatalf("overlay %v, Subtitles box %v: the sample plays without subtitles", overlay != nil, box != nil)
	}
	// shown is the overlay's text at pos. Set and read in one turn of the
	// main goroutine, so the playing video's own position updates, which
	// also arrive there, can't land in between.
	shown := func(pos time.Duration) string {
		text := ""
		bellui.DoAndWait(func() {
			overlay.SetPosition(pos)
			walkTree(overlay, func(o core.CanvasObject) {
				if c, ok := o.(*canvas.Text); ok {
					text = c.Text
				}
			})
		})
		return text
	}
	for _, c := range []struct {
		at   time.Duration
		want string
	}{
		{1 * time.Second, "A test pattern made for this gallery"},
		{2100 * time.Millisecond, ""},
		{3 * time.Second, "These lines come from a WebVTT track"},
		{5 * time.Second, "drawn over the video by SubtitleOverlay"},
	} {
		if got := shown(c.at); got != c.want {
			t.Errorf("at %v the overlay shows %q, want %q", c.at, got, c.want)
		}
	}

	bellui.DoAndWait(func() { box.SetChecked(false) })
	if overlay.Visible() {
		t.Error("unticking Subtitles leaves the overlay showing")
	}
}

// The Video Player page's own buttons drive the bare VideoPlayer, its clock
// follows playback, and leaving the page closes the player (gallery#47).
func TestVideoPlayerPage(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(900, 600))
	page := VideoPlayerPage(w)
	w.SetContent(page)
	defer LeavePage()

	var vp *mmedia.VideoPlayer
	buttons := map[string]*widget.Button{}
	var clock *widget.Label
	walkTree(page, func(o core.CanvasObject) {
		switch v := o.(type) {
		case *mmedia.VideoPlayer:
			vp = v
		case *widget.Button:
			buttons[v.Text] = v
		case *widget.Label:
			if strings.Contains(v.Text, " / ") {
				clock = v
			} else if strings.HasPrefix(v.Text, "Cannot play") && strings.Contains(v.Text, "audio device") {
				t.Skipf("no audio device: %s", v.Text)
			}
		}
	})
	if vp == nil || clock == nil || buttons["Play"] == nil || buttons["Pause"] == nil || buttons["Stop"] == nil {
		t.Fatalf("player %v, clock %v, buttons %v", vp != nil, clock != nil, buttons)
	}

	eventually(t, "the sample doesn't start playing", vp.IsPlaying)
	eventually(t, "the clock never shows the sample's 6 s length", func() bool {
		return strings.HasSuffix(clock.Text, "/ 0:06")
	})
	test.Tap(buttons["Pause"])
	eventually(t, "Pause leaves it playing", func() bool { return !vp.IsPlaying() })
	test.Tap(buttons["Play"])
	eventually(t, "Play doesn't resume it", vp.IsPlaying)

	LeavePage()
	eventually(t, "leaving the page leaves the video playing", func() bool { return !vp.IsPlaying() })
}

// fakeCamera stands in for a capture player: opening a real camera lights it
// and asks for permission.
type fakeCamera struct {
	*widget.Label
	mu     sync.Mutex
	closed bool
}

func (f *fakeCamera) Close() {
	f.mu.Lock()
	f.closed = true
	f.mu.Unlock()
}

func (f *fakeCamera) isClosed() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closed
}

// The Camera page opens nothing until a camera is chosen; then it opens that
// device and shows it, closes the one it replaces, says why an open failed,
// and closes the camera when the page is left, even one whose open was still
// under way (gallery#47).
func TestCameraPage(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	var mu sync.Mutex
	var opened []string
	var cams []*fakeCamera
	fail := false
	gate := make(chan struct{}, 10)
	saved := openCapture
	defer func() { openCapture = saved }()
	openCapture = func(device string) (capturePreview, error) {
		<-gate
		mu.Lock()
		defer mu.Unlock()
		opened = append(opened, device)
		if fail {
			return nil, errors.New("device busy")
		}
		c := &fakeCamera{Label: widget.NewLabel("live " + device)}
		cams = append(cams, c)
		return c, nil
	}
	count := func() int { mu.Lock(); defer mu.Unlock(); return len(opened) }

	w := test.NewWindow(nil)
	defer w.Close()
	w.Resize(geo.NewSize(900, 600))
	page := CameraPage(w)
	w.SetContent(page)

	findPicker := func() *mmedia.CaptureSelect {
		var picker *mmedia.CaptureSelect
		walkTree(page, func(o core.CanvasObject) {
			if p, ok := o.(*mmedia.CaptureSelect); ok {
				picker = p
			}
		})
		return picker
	}
	picker := findPicker()
	if picker == nil {
		t.Fatal("the page has no camera picker")
	}
	// Let one open through, so an open the page started on its own shows.
	gate <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	if n := count(); n != 0 {
		t.Fatalf("building the page opened %d devices", n)
	}
	<-gate

	shows := func(text string) bool {
		found := false
		walkTree(page, func(o core.CanvasObject) {
			switch v := o.(type) {
			case *widget.Label:
				found = found || v.Text == text
			case *fakeCamera:
				found = found || v.Text == text
			}
		})
		return found
	}
	choose := func(name, desc string) {
		bellui.DoAndWait(func() {
			picker.OnChosen(mmedia.CaptureDevice{Name: name, Description: desc, HasVideo: true})
		})
		gate <- struct{}{}
	}

	choose("cam0", "Front camera")
	eventually(t, "choosing Front camera shows no preview", func() bool {
		return shows("live cam0") && shows("Showing Front camera")
	})
	choose("cam1", "Back camera")
	eventually(t, "choosing Back camera shows no preview", func() bool { return shows("live cam1") })
	if !cams[0].isClosed() {
		t.Error("the Front camera stays open after switching to the Back one")
	}

	mu.Lock()
	fail = true
	mu.Unlock()
	choose("cam2", "Broken camera")
	eventually(t, "a failed open isn't reported", func() bool {
		return shows("Cannot open Broken camera: device busy")
	})
	mu.Lock()
	fail = false
	mu.Unlock()

	LeavePage()
	if !cams[1].isClosed() {
		t.Error("leaving the page leaves the Back camera open")
	}

	// An open still under way when the page is left is closed when it
	// lands, not shown.
	page = CameraPage(w)
	w.SetContent(page)
	picker = findPicker()
	bellui.DoAndWait(func() {
		picker.OnChosen(mmedia.CaptureDevice{Name: "late", Description: "Slow camera", HasVideo: true})
	})
	LeavePage()
	gate <- struct{}{}
	eventually(t, "a camera that opened after the page was left stays open", func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(cams) == 3 && cams[2].isClosed()
	})
	if shows("live late") {
		t.Error("a camera that opened after the page was left is shown")
	}
}
