package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("RemoveMember", func() {
	Describe("validations", func() {
		It("only accepts groupname and username", func() {
			session := runCommand("remove-member", "first-arg", "second-arg", "third-arg")
			Eventually(session).Should(Exit(1))

			session = runCommand("remove-member", "woodstock")
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("add-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(MISSING_TARGET))
		})
	})

	Describe("when no token in context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("tells the user to get a token", func() {
			session := runCommand("add-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(MISSING_CONTEXT))
		})
	})

	Describe("validations", func() {
		It("only accepts groupname and username", func() {
			session := runCommand("add-member", "first-arg", "second-arg", "third-arg")
			Eventually(session).Should(Exit(1))

			session = runCommand("add-member", "woodstock")
			Eventually(session).Should(Exit(1))
		})
	})
})
