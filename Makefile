.PHONY: build
build:
	go build -mod=vendor -o build/bin/iamd ./cmd/iamd