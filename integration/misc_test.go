//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("misc", func() {
	It("curls an arbitrary endpoint", func() {
		assertValidJSON(runOK("curl", "/info"))
	})

	It("gets userinfo for a user-context token", func() {
		runOK("get-password-token", "cf", "-s", "", "-u", "marissa", "-p", "koala")
		assertValidJSON(runOK("userinfo"))
		runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
	})
})
