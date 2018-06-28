#!/bin/bash

mkdir -p go/src/code.cloudfoundry.org/uaa-cli
cp -R uaa-cli/* go/src/code.cloudfoundry.org/uaa-cli

cd go/src/code.cloudfoundry.org/uaa-cli

go get github.com/onsi/ginkgo/ginkgo
ginkgo -v -r -randomizeAllSpecs -randomizeSuites

