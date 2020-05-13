# Root makefile, will define all the necessary variables and main calls.

APP=gtfo
PROJECT=github.com/nobe4/${APP}

GO?=go
GOOS?=darwin
GOARCH?=amd64

RELEASE?=0.0.4
COMMIT?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

MAIN_PATH=${PROJECT}/cmd/${APP}
BUILD_PATH=./bin/${GOOS}-${GOARCH}/${APP}

default: | run

.PHONY: version
version:
	@echo -n ${RELEASE}

.PHONY: bump
bump:
	./build/bump.sh

.PHONY: build
build:
	${GO} build \
		-a \
		-race \
		-installsuffix cgo \
		-o ${BUILD_PATH} \
		${MAIN_PATH}

.PHONY: run
run:
	${GO} run ./cmd/${APP}

.PHONY: test
test:
	# Run all the short tests from the host
	CGO_ENABLED=1 ${GO} test -race -short -cover ./...

.PHONY: test
gtfo:
	# Run all the short tests from the host
	CGO_ENABLED=1 ${GO} test -json -race -short -cover ./... | go run ./cmd/gtfo/gtfo.go

.PHONY: cover
cover:
	${GO} test -race -cover -coverprofile=coverprofile.out -short ./... && \
		${GO} tool cover -html=coverprofile.out && \
		rm coverprofile.out

.PHONY: lint
lint:
	golangci-lint run \
	--config build/golangci.yml \
	./...
