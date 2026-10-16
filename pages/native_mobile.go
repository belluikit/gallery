//go:build android || js

package pages

// NativePages returns no extra destinations on Android and on the web. The
// Clair and Code Assistant pages depend on bleve's on-disk storage (bbolt,
// mmap-go), which does not build for js/wasm and is kept out of the Android
// build. See native_desktop.go for the desktop set.
func NativePages() []PageDef { return nil }
