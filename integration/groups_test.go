//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("groups", Ordered, func() {
	groupName := "uaa.cli.integration.test.group"
	memberUsername := "uaa-cli-integration-test-group-member"

	BeforeAll(func() {
		runOK("create-user", memberUsername,
			"--givenName", "Group", "--familyName", "Member",
			"--email", memberUsername+"@example.com", "-p", "S0meSecur3Pass!",
		)
	})

	AfterAll(func() {
		run("delete-user", memberUsername)
		// uaa-cli has no delete-group command yet, so clean up via the raw
		// SCIM endpoint to avoid a name collision if this suite runs again
		// against the same (non-freshly-booted) UAA instance.
		if id := groupID(groupName); id != "" {
			run("curl", "-X", "DELETE", "/Groups/"+id)
		}
	})

	It("creates a group", func() {
		runOK("create-group", groupName, "-d", "uaa-cli integration test group")
	})

	It("gets the group", func() {
		assertValidJSON(runOK("get-group", groupName))
	})

	It("lists groups", func() {
		assertValidJSON(runOK("list-groups"))
	})

	It("adds and removes a member", func() {
		runOK("add-member", groupName, memberUsername)
		runOK("remove-member", groupName, memberUsername)
	})

	It("maps and unmaps an external group", func() {
		runOK("map-group", "cn=integration-test,ou=groups,dc=example,dc=com", groupName)
		assertValidJSON(runOK("list-group-mappings"))
		runOK("unmap-group", "cn=integration-test,ou=groups,dc=example,dc=com", groupName)
	})
})
