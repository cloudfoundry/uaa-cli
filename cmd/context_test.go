package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Context", func() {
	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfig()
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("context")

			Expect(session.Err).To(Say("No context is currently set."))
			Expect(session.Err).To(Say(`To get started, target a UAA and fetch a token. See "uaa target -h" for details`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when a target was previously set but there is no active context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL("http://login.somewhere.com")
			config.WriteConfig(c)
		})

		It("tells the user to set a context", func() {
			session := runCommand("context")

			Expect(session.Err).To(Say("No context is currently set."))
			Expect(session.Err).To(Say(`Use a token command such as "uaa get-password-token" or "uaa get-client-credentials-token" to fetch a token.`))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when there is an active context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL("http://login.somewhere.com")
			ctx := config.UaaContext{ClientId: "admin", Username: "woodstock"}
			ctx.Token.TokenType = "bearer"
			ctx.Token.AccessToken = "token"
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("displays the context", func() {
			activeContextJson := `{
			  "client_id": "admin",
			  "grant_type": "",
			  "username": "woodstock",
				"Token": {
					"access_token": "token",
					"token_type": "bearer",
					"expiry": "0001-01-01T00:00:00Z"
				}
			}`
			session := runCommand("context")

			Expect(session.Out.Contents()).To(MatchJSON(activeContextJson))
			Eventually(session).Should(Exit(0))
		})

		It("displays the context access_token", func() {
			session := runCommand("context", "--access_token")

			Expect(session.Out).To(Say(`token`))
			Eventually(session).Should(Exit(0))
		})

		It("displays the context authentication header", func() {
			session := runCommand("context", "--auth_header")

			Expect(session.Out).To(Say(`bearer token`))
			Eventually(session).Should(Exit(0))
		})
	})
})
