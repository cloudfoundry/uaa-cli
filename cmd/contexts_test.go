package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Contexts", func() {
	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("contexts")

			Expect(session.Out).To(Say("No contexts are currently available."))
			Expect(session.Out).To(Say(`To get started, target a UAA and fetch a token. See "uaa target -h" for details.`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when a target was previously set but there is no active context", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL("http://login.somewhere.com")
			config.WriteConfig(c)
		})

		It("tells the user to set a context", func() {
			session := runCommand("contexts")

			Expect(session.Out).To(Say("No contexts are currently available."))
			Expect(session.Out).To(Say(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token and create a context.`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when there are contexts", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL("http://login.somewhere.com")
			ctx1 := uaa.UaaContext{ClientId: "admin", Username: "woodstock", GrantType: uaa.PASSWORD}
			c.AddContext(ctx1)
			config.WriteConfig(c)
		})

		It("prints a table of results", func() {
			session := runCommand("contexts")

			// Headings
			Expect(session.Out).Should(Say("CLIENTID"))
			Expect(session.Out).Should(Say("USERNAME"))
			Expect(session.Out).Should(Say("GRANT TYPE"))

			Expect(session.Out).Should(Say("admin"))
			Expect(session.Out).Should(Say("woodstock"))
			Expect(session.Out).Should(Say("password"))

			Eventually(session).Should(Exit(0))
		})
	})

})
