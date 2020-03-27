
BIN := iamd
BUILD_FLAGS = -mod=vendor

.PHONY: default
default: build

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o build/bin/${BIN} ./cmd/${BIN}

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