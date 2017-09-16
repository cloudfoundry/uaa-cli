.PHONY: all build ci clean dep format ginkgo test install

DEST = build/uaa
INSTALL_DEST = $(GOPATH)/bin/uaa
COMMIT_HASH=`git rev-parse --short HEAD`

ifndef VERSION
	VERSION = DEV
endif

GOFLAGS := -v -o $(DEST) -ldflags "-X code.cloudfoundry.org/uaa-cli/version.Version=${VERSION} -X code.cloudfoundry.org/uaa-cli/version.Commit=${COMMIT_HASH}"

all: dep test clean build

clean:
		rm -rf build

format:
		go fmt .

ginkgo:
		ginkgo -r -randomizeSuites -randomizeAllSpecs -race 2>&1

test: format ginkgo

dep:
		go get github.com/onsi/ginkgo/ginkgo
		go get github.com/onsi/gomega
		go get -u github.com/golang/dep/cmd/dep
		go install github.com/golang/dep/cmd/dep
		dep ensure

ci: ginkgo

build:
		mkdir -p build
		go build $(GOFLAGS)

install:
		rm -rf $(INSTALL_DEST)
		cp $(DEST) $(INSTALL_DEST)
