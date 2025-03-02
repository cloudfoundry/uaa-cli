.PHONY: all build ci clean format ginkgo test install

BUILD_DEST = build/uaa
INSTALL_DEST = $(GOPATH)/bin/uaa
COMMIT_HASH=$(shell git rev-parse --short HEAD)
GOFILES=$(shell find . -type f -name '*.go')

ifndef VERSION
	VERSION = DEV
endif

GOFLAGS := -v -ldflags "-X code.cloudfoundry.org/uaa-cli/version.Version=${VERSION} -X code.cloudfoundry.org/uaa-cli/version.Commit=${COMMIT_HASH}"

all: test clean build

clean:
		rm -rf build

format:
		go fmt ./...

ginkgo:
		go run github.com/onsi/ginkgo/v2/ginkgo -v -r --randomize-suites --randomize-all -race

test: format ginkgo

ci: ginkgo

build:
		mkdir -p build
		CGO_ENABLED=0 go build $(GOFLAGS) -o $(BUILD_DEST)

build_all:
		mkdir -p build
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DEST)-darwin-arm64
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-darwin-amd64
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-linux-amd64
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-windows-amd64

install:
		rm -rf $(INSTALL_DEST)
		cp $(BUILD_DEST) $(INSTALL_DEST)

cve: build
		grype file:$(BUILD_DEST)