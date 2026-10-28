//go:build js

package pages

// AudioPages returns no destinations on the web: x/mmedia links FFmpeg
// through cgo. See audio_player.go for the other platforms.
func AudioPages() []PageDef { return nil }
