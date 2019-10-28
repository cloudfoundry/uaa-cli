package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"
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

	Describe("validations", func() {
		Describe("without a target and context", func() {
			It("requires a target", func() {
				err := UnmapGroupValidations(config.Config{}, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
			})

			It("requires a context", func() {
				cfg := config.NewConfigWithServerURL("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
			})
		})

		Describe("without required params", func() {
			It("requires a external_group_name", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
			})

			It("requires a group_name", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{"external_group"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
			})
		})

		Describe("with totally valid data", func() {
			It("does not complain", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{"external_groupname", "groupname"})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
