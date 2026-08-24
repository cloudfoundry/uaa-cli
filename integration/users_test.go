//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("users", Ordered, func() {
	username := "uaa-cli-integration-test-user"

	AfterAll(func() {
		run("delete-user", username)
	})

	It("creates a user", func() {
		runOK("create-user", username,
			"--givenName", "Integration",
			"--familyName", "Test",
			"--email", username+"@example.com",
			"-p", "S0meSecur3Pass!",
		)
	})

	It("gets the user", func() {
		assertValidJSON(runOK("get-user", username))
	})

	It("lists users with a filter", func() {
		assertValidJSON(runOK("list-users", "--filter", `userName eq "`+username+`"`))
	})

	It("updates the user", func() {
		runOK("update-user", username, "--givenName", "Updated")
	})

	It("deactivates and reactivates the user", func() {
		runOK("deactivate-user", username)
		runOK("activate-user", username)
	})

	PIt("unlocks the user", func() {
		// PATCH /Users/{userId}/status always 500s on current UAA: the
		// updateAccountStatus handler in ScimUserEndpoints.java is missing
		// @ResponseBody (accidentally dropped in cloudfoundry/uaa commit
		// 6159e2f4, 2016), so Spring tries to resolve an HTML view instead
		// of returning JSON. Reproduced against a freshly-booted UAA, not
		// suite/state-dependent. Un-pend once fixed upstream.
		runOK("unlock-user", username)
	})
})
