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

	It("unlocks the user", func() {
		runOK("unlock-user", username)
	})
})
