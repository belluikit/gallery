# Gallery

## Checking a change

The gallery has no CI, so each PR carries this evidence from a local run:

```sh
go build ./...
go vet ./...
go test -count=1 ./...
go test -count=1 -tags mobile -v ./... 2>&1 | grep -E -- '--- (PASS|SKIP|FAIL): Test(ChartHeightsOnPhone|DAGPhoneLayout)'
CGO_ENABLED=0 GOOS=js GOARCH=wasm go vet ./...
CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o /dev/null .
```

**The phone tests need `-tags mobile`.** Without the tag, bellui reports a
desktop device, so `TestChartHeightsOnPhone` and `TestDAGPhoneLayout` skip and
`go test ./...` still prints `ok`. Run the suite a second time with the tag,
and check the `-v` output shows those two tests as `--- PASS`, not `--- SKIP`.
A new test that checks the phone layout should skip the same way, through
`bellui.CurrentDevice().IsMobile()`, and should be added to the line above.

The web vet needs `CGO_ENABLED=0`, because `CGO_ENABLED=1` is often set
globally and the js/wasm target has no cgo.
