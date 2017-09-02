!#/bin/bash

set -ex

go version


echo "Installing dep"

go get -u github.com/golang/dep/cmd/dep
go install github.com/golang/dep/cmd/dep

echo "Installing dependencies with dep"
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
dep version
dep ensure

ls vendor

echo "Running tests"
ginkgo version
ginkgo -r -randomizeAllSpecs -randomizeSuites
