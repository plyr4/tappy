# tapper

Terminal clicking-speed game for macOS. Tap 5 times quickly to start a 10-second challenge.

## Install

```
brew install plyr4/tap/tapper
```

Or from source:

```
go install github.com/plyr4/tappy/cmd/tapper@latest
```

## Usage

Run `tapper` in a terminal with mouse reporting support (e.g. Ghostty).

- Left-click / trackpad tap 5 times within 1 second to start.
- Click as fast as you can for 10 seconds.
- `esc` quits.

Personal best is stored in `~/Library/Application Support/tapper/best.json`.

## Development

```
make build    # bin/tapper
make test
make release  # dist/tapper-darwin-{arm64,amd64}
make clean
```
