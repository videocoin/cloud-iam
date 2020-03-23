NAME := iam
BIN := iamd
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)
LD_FLAGS = -X main.Version=${VERSION} -s -w
#BUILD_FLAGS = -mod=vendor -ldflags "$(LD_FLAGS)"
BUILD_FLAGS = -mod=vendor

.PHONY: default
default: build

.PHONY: all
all: build

.PHONY: iamd
iamd:
	go build $(BUILD_FLAGS) -o build/bin/iamd ./cmd/iamd

.PHONY: metadatad
metadatad:
	go build $(BUILD_FLAGS) -o build/bin/metadatad ./cmd/metadatad

.PHONY: e2e
e2e:
	cd test/e2e && docker-compose up -d --build mysql
	sleep 30
	cd test/e2e && docker-compose up -d --build migrate ${BIN}
	sleep 3
	go test -mod=vendor ./test/e2e/...

.PHONY: e2e-nobuild
e2e-nobuild:
	cd test/e2e && docker-compose up -d --no-build mysql 
	sleep 30
	cd test/e2e && docker-compose up -d --no-build migrate ${BIN}
	sleep 3
	go test -mod=vendor ./test/e2e/...
	
.PHONY: e2e-stop
e2e-stop:
	cd test/e2e && docker-compose down --volumes


.phony: proto
proto:
	protoc -I . -I ${GOPATH}/src/github.com/googleapis/api-common-protos --go_out=plugins=grpc:. api/admin/v1/rbac.proto
	