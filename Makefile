VERSION ?= `git rev-parse HEAD`
DOCKER_REGISTRY ?= "registry.videocoin.net"

.PHONY: build
build:
		@go build -o build/bin/iamd ./cmd/iamd

.PHONY: vendor
vendor:
		@go mod tidy
		@go mod vendor
		@modvendor -copy="**/*.c **/*.h" -v

.PHONY: docker-build
docker-build:
		@printf "\033[32m\xE2\x9c\x93 Building IAM docker image...\n\033[0m"
		@export DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t ${DOCKER_REGISTRY}/cloud-iam/iam:${VERSION} .

.PHONY: docker-push-to-harbor
docker-push-to-harbor:
		@docker push ${DOCKER_REGISTRY}/cloud-iam/iam:${VERSION}