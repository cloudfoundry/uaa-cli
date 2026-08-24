//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("cli basics", func() {
	It("shows server info", func() {
		assertValidJSON(runOK("info"))
	})

	It("shows the current context", func() {
		runOK("context")
	})

	It("lists saved contexts", func() {
		runOK("contexts")
	})

	It("prints the version", func() {
		runOK("version")
	})
})
