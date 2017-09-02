!#/bin/bash

set -e

go version


echo "Installing dep"

go get -u github.com/golang/dep/cmd/dep
go install github.com/golang/dep/cmd/dep

echo "Installing dependencies with dep"
dep version
dep install

echo "Running tests"
ginkgo version
ginkgo -r -randomizeAllSpecs -randomizeSuites
