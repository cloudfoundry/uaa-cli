//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("tokens", func() {
	BeforeEach(func() {
		runOK("get-password-token", "cf", "-s", "", "-u", "marissa", "-p", "koala")
	})

	AfterEach(func() {
		// Other Describes assume an admin token is active; if this fails
		// silently, later suites would run with the wrong token and fail in
		// confusing ways, so assert it rather than best-effort it.
		runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
	})

	It("decodes the current access token", func() {
		token := string(runOK("context", "--access_token").Out.Contents())
		runOK("decode-token", token)
	})

	It("gets the token signing key(s)", func() {
		runOK("get-token-key")
		runOK("get-token-keys")
	})

	It("refreshes the token", func() {
		runOK("refresh-token", "-s", "")
	})
})
