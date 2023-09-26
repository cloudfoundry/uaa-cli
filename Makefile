.PHONY: all build ci clean dep format ginkgo test install

BUILD_DEST = build/uaa
INSTALL_DEST = $(GOPATH)/bin/uaa
COMMIT_HASH=`git rev-parse --short HEAD`
GOFILES=`find . -type f -name '*.go'`

ifndef VERSION
	VERSION = DEV
endif

GOFLAGS := -v -ldflags "-X code.cloudfoundry.org/uaa-cli/version.Version=${VERSION} -X code.cloudfoundry.org/uaa-cli/version.Commit=${COMMIT_HASH}"

all: dep test clean build

clean:
		rm -rf build

format:
		gofmt -l -s -w ${GOFILES}

ginkgo:
		bin/test 2>&1

test: format ginkgo

dep:
		go install github.com/onsi/ginkgo/ginkgo@latest

ci: ginkgo

build:
		mkdir -p build
		go build $(GOFLAGS) -o $(BUILD_DEST)

build_all:
		mkdir -p build
		GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DEST)-darwin-arm64
		GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-darwin-amd64
		GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-linux-amd64
		GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DEST)-windows-amd64

install:
		rm -rf $(INSTALL_DEST)
		cp $(BUILD_DEST) $(INSTALL_DEST)
