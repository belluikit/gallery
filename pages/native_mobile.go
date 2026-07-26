//go:build android

package pages

// NativePages returns no extra destinations on Android. The Clair / Code
// Assistant (FAISS) and Media Player (FFmpeg) pages link host-only native
// libraries that aren't part of the mobile build, so they — and their heavy
// dependencies — are excluded here. See native_desktop.go for the desktop set.
func NativePages() []PageDef { return nil }
