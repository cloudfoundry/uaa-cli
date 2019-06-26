package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("GetClientCredentialsToken", func() {
	var opaqueTokenResponseJson = `{
	  "access_token" : "bc4885d950854fed9a938e96b13ca519",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`
	var jwtResponseJson = `{
	  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`
	var c config.Config

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c = config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		Describe("when successful", func() {
			BeforeEach(func() {
				config.WriteConfig(c)
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, opaqueTokenResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
					VerifyFormKV("grant_type", "client_credentials"),
				),
				)
			})

			It("displays a success message", func() {
				session := runCommand("get-client-credentials-token", "admin", "-s", "adminsecret")
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Access token successfully fetched."))
			})

			It("updates the saved context", func() {
				runCommand("get-client-credentials-token", "admin", "-s", "adminsecret", "--format", "jwt")
				Expect(config.ReadConfig().GetActiveContext().Token.AccessToken).To(Equal("bc4885d950854fed9a938e96b13ca519"))
				Expect(config.ReadConfig().GetActiveContext().Token.TokenType).To(Equal("bearer"))
			})
		})

		Describe("configuring token format", func() {
			It("can request jwt token", func() {
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, jwtResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
					VerifyFormKV("grant_type", "client_credentials"),
					VerifyFormKV("token_format", "jwt"),
				))

				runCommand("get-client-credentials-token", "admin", "-s", "adminsecret", "--format", "jwt")
			})

			It("can request opaque token", func() {
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, opaqueTokenResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
					VerifyFormKV("grant_type", "client_credentials"),
					VerifyFormKV("token_format", "opaque"),
				))

				runCommand("get-client-credentials-token", "admin", "-s", "adminsecret", "--format", "opaque")
			})

			It("uses jwt format by default", func() {
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, jwtResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
					VerifyFormKV("grant_type", "client_credentials"),
					VerifyFormKV("token_format", "jwt"),
				))

				runCommand("get-client-credentials-token", "admin", "-s", "adminsecret")
			})

			It("displays error when unknown format is passed", func() {
				session := runCommand("get-client-credentials-token", "admin", "-s", "adminsecret", "--format", "bogus")
				Expect(session.Err).To(Say(`The token format "bogus" is unknown.`))
				Expect(session).To(Exit(1))
			})
		})
	})

	Describe("when the token request fails", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("old-token"))
			config.WriteConfig(c)
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusUnauthorized, `{"error":"unauthorized","error_description":"Bad credentials"}`),
				VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
				VerifyFormKV("grant_type", "client_credentials"),
			),
			)
		})

		It("displays help to the user", func() {
			session := runCommand("get-client-credentials-token", "admin", "-s", "adminsecret")
			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("An error occurred while fetching token."))
		})

		It("does not update the previously saved context", func() {
			runCommand("get-client-credentials-token", "admin", "-s", "adminsecret")
			Expect(config.ReadConfig().GetActiveContext().Token.AccessToken).To(Equal("old-token"))
		})
	})

	Describe("Validations", func() {
		Describe("when called with no client id", func() {
			It("displays help and does not panic", func() {
				c := config.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-client-credentials-token")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
			})
		})

		Describe("when called with no client secret", func() {
			It("displays help and does not panic", func() {
				c := config.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-client-credentials-token", "admin")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `client_secret` must be specified."))
			})
		})

		Describe("when no target was previously set", func() {
			BeforeEach(func() {
				config.WriteConfig(config.NewConfig())
			})

			It("tells the user to set a target", func() {
				session := runCommand("get-client-credentials-token")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("You must set a target in order to use this command."))
			})
		})
	})
})
