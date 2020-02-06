-include ./test/e2e/e2e.mk

NAME := iam
BIN := iamd
VERSION=$$(git describe --abbrev=0 --always)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)
LD_FLAGS = -X main.Version=${VERSION} -s -w
BUILD_FLAGS = -mod=vendor -ldflags "$(LD_FLAGS)"
OUTPUT ?= build/bin/${BIN}

GCP_PROJECT ?= videocoin-network
DOCKER_REGISTRY ?= gcr.io/${GCP_PROJECT}

.PHONY: default
default: build

.PHONY: all
all: build

.PHONY: version
version:
		@echo ${VERSION}

.PHONY: build
build:
		go build $(BUILD_FLAGS) -o $(OUTPUT) ./cmd/${BIN}

.PHONY: install
install:
		go install $(BUILD_FLAGS) ./cmd/${BIN}

.PHONY: deps		
deps:
		GO111MODULE=on go mod vendor

.PHONY: docker-build
docker-build:
		docker build -t ${DOCKER_REGISTRY}/${NAME}:${VERSION} .

.PHONY: docker-push
docker-push:
		docker push gcr.io/${GCP_PROJECT}/${NAME}:${VERSION}

.PHONY: release
release: docker-build docker-push

.PHONY: deploy
deploy:
		ENV=${ENV} deployments/deploy.sh
