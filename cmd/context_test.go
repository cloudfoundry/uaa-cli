package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Context", func() {
	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("context")

			Expect(session.Out).To(Say("No context is currently set."))
			Expect(session.Out).To(Say(`To get started, target a UAA and fetch a token. See "uaa target -h" for details`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when a target was previously set but there is no active context", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL("http://login.somewhere.com")
			config.WriteConfig(c)
		})

		It("tells the user to set a context", func() {
			session := runCommand("context")

			Expect(session.Out).To(Say("No context is currently set."))
			Expect(session.Out).To(Say(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token.`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when there is an active context", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL("http://login.somewhere.com")
			ctx := uaa.UaaContext{ClientId: "admin", Username: "woodstock"}
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("displays the context", func() {
			activeContextJson := `{ "ClientId": "admin", "GrantType": "", "Username": "woodstock", "AccessToken": "",
				   "TokenType": "", "ExpiresIn": 0, "Scope": "", "JTI": "" }`
			session := runCommand("context")

			Expect(session.Out.Contents()).To(MatchJSON(activeContextJson))
			Eventually(session).Should(Exit(0))
		})
	})
})
