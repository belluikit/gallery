//go:build !android

package pages

import "github.com/belluikit/bellui/theme"

// NativePages returns the gallery destinations that depend on host-only native
// libraries — Clair and the Code Assistant link FAISS (go-faiss) for vector
// search, and the Media Player links FFmpeg (goffmpeg) plus the oto audio
// backend. They are only available in desktop builds; the Android build uses
// the stub in native_mobile.go, which returns no extra pages.
func NativePages() []PageDef {
	return []PageDef{
		{Name: "Clair", Icon: theme.SearchIcon(), New: ClairPage},
		{Name: "Code Assistant", Icon: theme.ComputerIcon(), New: CodeAssistantPage},
	}
}
