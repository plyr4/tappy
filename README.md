# tappy

Terminal clicking-speed game for macOS. Tap 5 times quickly to start a 10-second challenge.

## Install

```
brew install plyr4/tap/tappy
```

Or from source:

```
go install github.com/plyr4/tappy/cmd/tappy@latest
```

## Usage

Run `tappy` in a terminal with mouse reporting support (e.g. Ghostty).

- Left-click / trackpad tap 5 times within 1 second to start.
- Click as fast as you can for 10 seconds.
- `esc` quits.

Personal best is stored in `~/Library/Application Support/tappy/best.json`.

## Development

```
make build    # bin/tappy
make test
make lint     # requires golangci-lint
make release  # dist/tappy-darwin-{arm64,amd64}
make clean
```
