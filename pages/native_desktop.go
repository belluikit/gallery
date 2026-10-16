//go:build !android && !js

package pages

import "github.com/belluikit/bellui/theme"

// NativePages returns the gallery destinations that only build for the
// desktop. Clair and the Code Assistant keep their knowledge base in bleve,
// whose storage (bbolt, mmap-go) has no js/wasm implementation, and the Media
// Player links FFmpeg through goffmpeg (cgo). Android and web builds use the
// stub in native_mobile.go, which returns no extra pages.
func NativePages() []PageDef {
	return []PageDef{
		{Name: "Clair", Icon: theme.SearchIcon(), New: ClairPage},
		{Name: "Code Assistant", Icon: theme.ComputerIcon(), New: CodeAssistantPage},
		{Name: "Media Player", Icon: theme.MediaPlayIcon(), New: MediaPlayerPage, Fill: true},
	}
}
