.PHONY: build
iam:
	go build -mod=vendor -o build/bin/iamd ./cmd/iamd

metadata:
	go build -mod=vendor -o build/bin/metadatad ./cmd/metadatad	