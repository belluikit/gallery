# Material Design 3 Gallery

A showcase app for [Bellui](https://github.com/belluikit/bellui), a Go UI
toolkit forked from Fyne. Each page demonstrates a group of Material Design 3
components (actions, inputs, navigation, containment, typography, theming and
more) plus the Bellui extension modules: the code editor (`x/editor`), charts
(`x/charts`) and the data grid (`x/tables`).

Homepage: https://bellui.com

## Run it

You need the Go version named in `go.mod` and the platform prerequisites
listed in the Bellui README.

```sh
go run .
```

The top app bar has toggles for light/dark, dense theme, right-to-left layout
and an About box. On desktop the Bellui dev inspector is enabled.

### In a browser

The gallery also runs as WebAssembly. Install the `bellui` command once, then
serve the gallery from this folder:

```sh
go install github.com/belluikit/bellui/cmd/bellui@latest
bellui serve
```

Open <http://localhost:8080> in a browser with WebGPU (Chrome or Edge;
verified in Chromium). The first load downloads a large `.wasm` file (about
64 MB uncompressed), so give it a moment. A few pages are desktop-only; the
table below says which.

To build static files for a web server instead, run
`bellui package -os web --release` and copy the `wasm` folder it creates.
WebGPU needs HTTPS (or `localhost`); see
[Bellui on the web](https://bellui.com/docs/platforms/web/) for the WebGL
alternative and other details.

## Pages and platforms

`main.go` registers the core pages, which build on every platform including
Android and the web; among them the DAG Editor (`x/dag`). Three more pages are
appended by `pages.NativePages()`:

| Page | Desktop | Android | Web |
| --- | --- | --- | --- |
| Clair | yes | no | no |
| Code Assistant | yes | no | no |
| Media Player (`x/mmedia`, FFmpeg) | yes | no | no |

The selection is done by build tags: `pages/native_desktop.go`,
`pages/clair_demo.go`, `pages/code_assistant.go` and `pages/media_player.go`
carry `//go:build !android && !js`, and `pages/native_mobile.go`
(`//go:build android || js`) returns no extra pages. Clair and the Code
Assistant keep their knowledge base in bleve, whose storage (bbolt, mmap-go) does not build for
js/wasm; the Media Player links FFmpeg through goffmpeg (cgo). It plays a
bundled six-second test video (`pages/media/sample.mp4`, made for the gallery)
or any local file opened with "Open file…".

Some pages need something outside the app to do anything useful:

- **Code Assistant** asks for an Anthropic API key, entered on the page and
  kept in the app's preferences.
- **File Transfer** talks to an uploader server (default `localhost:9090`).
- **Clair** builds a search index of Bellui's widget package in the system
  temp directory on first use.

## Android

`scripts/android-build.sh` packages an arm64 APK, repacks it so Android 11+
(API 30) accepts it (`resources.arsc` stored and aligned), signs it with the
debug key, and with `--install` installs and launches it through `adb`. It
expects the `bellui` command on `PATH`
(`go install github.com/belluikit/bellui/cmd/bellui@latest`), the Android SDK
in `$ANDROID_HOME` (default `~/Library/Android/sdk`) with build-tools 34.0.0,
and an NDK in `$ANDROID_NDK_HOME`.

```sh
scripts/android-build.sh [--install]
```

## License

BSD 3-Clause, see [LICENSE](LICENSE). Third-party material is listed in
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
