EXECUTABLE_NAME = uaa
INSTALL_DEST = $(GOPATH)/bin/uaa

.PHONY: all
all: format test build  ## Runs: Format, Test, and Build

${EXECUTABLE_NAME}: **/*.go go.mod go.sum   ## Build
	goreleaser build --snapshot --clean --single-target --output .

**/*.go:
	@# noop

.PHONY: build
build: ${EXECUTABLE_NAME}  ## Build

.PHONY: build_all
build_all:  **/*.go go.mod go.sum ## Build all releases, in dist, using Goreleaser (requires goreleaser to be installed)
		goreleaser build --snapshot --clean

.PHONY: clean
clean: ## Clean up artifacts
		rm -rf build/ dist/ ${EXECUTABLE_NAME}

.PHONY: format
format:  ## Format Go Code
		go fmt ./...

.PHONY: test
test:  ## Run Ginkgo tests
		go run github.com/onsi/ginkgo/v2/ginkgo -v -r --randomize-suites --randomize-all -race

.PHONY: goreleaser-check
goreleaser-check:  ## Test goreleaser configuration
		goreleaser check

.PHONY: install
install: ${EXECUTABLE_NAME} ## Install executable
		rm -rf $(INSTALL_DEST)
		cp ${EXECUTABLE_NAME} $(INSTALL_DEST)

.PHONY: cve
cve: ${EXECUTABLE_NAME}  ## Scan for CVEs
		grype file:${EXECUTABLE_NAME}

.PHONY: setup
setup:  ## Setup packages needed for release
	brew install caarlos0/tap/svu
	brew install goreleaser/tap/goreleaser

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)