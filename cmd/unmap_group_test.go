package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("UnmapGroup", func() {
	Describe("by default", func() {
		BeforeEach(func() {
			config.WriteConfig(buildConfig(server.URL()))
		})

		It("Resolves the group name and performs the unmapping", func() {
			mockGroupLookup("internal-group-id", "internal-group")
			mockExternalGroupUnmapping("external-group", "internal-group-id", "internal-group", "ldap")

			session := runCommand("unmap-group", "external-group", "internal-group")
			Eventually(session).Should(Say(`Successfully unmapped internal-group from external-group for origin ldap`))
			Eventually(session).Should(Exit(0))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Describe("with origin", func() {
		BeforeEach(func() {
			config.WriteConfig(buildConfig(server.URL()))
		})

		It("Resolves the group name and performs the unmapping", func() {
			mockGroupLookup("internal-group-id", "internal-group")
			mockExternalGroupUnmapping("external-group", "internal-group-id", "internal-group", "saml")

			session := runCommand("unmap-group", "external-group", "internal-group", "--origin", "saml")
			Eventually(session).Should(Say(`Successfully unmapped internal-group from external-group for origin saml`))
			Eventually(session).Should(Exit(0))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Describe("with invalid input", func() {
		BeforeEach(func() {
			config.WriteConfig(buildConfig(server.URL()))
		})

		It("fails", func() {
			session := runCommand("unmap-group", "external-group")
			Eventually(session.Err).Should(Say(`The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified.`))
			Eventually(session).Should(Exit(1))
			Expect(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})
})
