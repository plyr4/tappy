BIN := tapper
PKG := ./cmd/tapper

.PHONY: build test clean release

build:
	go build -o bin/$(BIN) $(PKG)

test:
	go test ./...

release:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o dist/$(BIN)-darwin-arm64 $(PKG)
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/$(BIN)-darwin-amd64 $(PKG)

clean:
	rm -rf bin dist
