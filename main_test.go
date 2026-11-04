//go:build !js

package main

import (
	"slices"
	"testing"
)

// The drawer lists pages A–Z, including the ones only some platforms add
// (the desktop-only pages, the audio player). Actions stays first, so the
// gallery still opens on it.
func TestPageRegistrySorted(t *testing.T) {
	var names []string
	for _, p := range pageRegistry() {
		names = append(names, p.Name)
	}
	if !slices.IsSorted(names) {
		t.Errorf("drawer pages are not in A–Z order: %q", names)
	}
	if names[0] != "Actions" {
		t.Errorf("first page is %q, want Actions", names[0])
	}
	if !slices.Contains(names, "Audio Player") {
		t.Errorf("Audio Player is missing from this build's pages: %q", names)
	}
}
