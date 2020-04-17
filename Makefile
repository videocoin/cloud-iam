.PHONY: iam proxy metadata
iam:
	go build -mod=vendor -o build/bin/iamd ./cmd/iamd

proxy:
	go build -mod=vendor -o build/bin/proxyd ./cmd/proxyd	

metadata:
	go build -mod=vendor -o build/bin/metadatad ./cmd/metadatad