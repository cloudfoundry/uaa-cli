package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("GetPasswordToken", func() {

	var opaqueTokenResponseJson = `{
	  "access_token" : "bc4885d950854fed9a938e96b13ca519",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`
	var jwtTokenResponseJson = `{
	  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`

	var c uaa.Config
	var ctx uaa.UaaContext

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c = uaa.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
			ctx = c.GetActiveContext()
		})

		Describe("when the --trace option is used", func() {
			It("shows extra output about the request on success", func() {
				server.RouteToHandler("POST", "/oauth/token",
					RespondWith(http.StatusOK, jwtTokenResponseJson),
				)

				session := runCommand("get-password-token",
					"admin",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret",
					"--trace")

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("POST " + server.URL() + "/oauth/token"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("200 OK"))
			})

			It("shows extra output about the request on error", func() {
				server.RouteToHandler("POST", "/oauth/token",
					RespondWith(http.StatusBadRequest, "garbage response"),
				)

				session := runCommand("get-password-token",
					"admin",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret",
					"--trace")

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("POST " + server.URL() + "/oauth/token"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("400 Bad Request"))
				Expect(session.Out).To(Say("garbage response"))
			})
		})

		Describe("when successful", func() {
			BeforeEach(func() {
				config.WriteConfig(c)
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, jwtTokenResponseJson),
					VerifyFormKV("client_id", "admin"),
					VerifyFormKV("client_secret", "adminsecret"),
					VerifyFormKV("grant_type", "password"),
				),
				)
			})

			It("displays a success message", func() {
				session := runCommand("get-password-token",
					"admin",
					"-s", "adminsecret",
					"--username", "woodstock",
					"--password", "secret")

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Access token successfully fetched."))
			})

			It("updates the saved context", func() {
				runCommand("get-password-token",
					"admin",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret")

				Expect(config.ReadConfig().GetActiveContext().AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
				Expect(config.ReadConfig().GetActiveContext().ClientId).To(Equal("admin"))
				Expect(config.ReadConfig().GetActiveContext().Username).To(Equal("woodstock"))
				Expect(config.ReadConfig().GetActiveContext().GrantType).To(Equal(uaa.PASSWORD))
				Expect(config.ReadConfig().GetActiveContext().TokenType).To(Equal(uaa.JWT))
				Expect(config.ReadConfig().GetActiveContext().ExpiresIn).To(Equal(int32(43199)))
				Expect(config.ReadConfig().GetActiveContext().Scope).To(Equal("clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write"))
				Expect(config.ReadConfig().GetActiveContext().JTI).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			})
		})
	})

	Describe("when the token request fails", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			c.AddContext(uaa.UaaContext{AccessToken: "old-token"})
			config.WriteConfig(c)
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusUnauthorized, `{"error":"unauthorized","error_description":"Bad credentials"}`),
				VerifyFormKV("client_id", "admin"),
				VerifyFormKV("client_secret", "adminsecret"),
				VerifyFormKV("grant_type", "password"),
			),
			)
		})

		It("displays help to the user", func() {
			session := runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("An error occurred while fetching token."))
		})

		It("does not update the previously saved context", func() {
			runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret")
			Expect(config.ReadConfig().GetActiveContext().AccessToken).To(Equal("old-token"))
		})
	})

	Describe("configuring token format", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			c.AddContext(uaa.UaaContext{AccessToken: "access_token"})
			config.WriteConfig(c)
		})

		It("can request jwt token", func() {
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusOK, jwtTokenResponseJson),
				VerifyFormKV("client_id", "admin"),
				VerifyFormKV("client_secret", "adminsecret"),
				VerifyFormKV("grant_type", "password"),
				VerifyFormKV("token_format", "jwt"),
			))

			runCommand("get-password-token", "admin",
											"-s", "adminsecret",
											"-u", "woodstock",
											"-p", "secret",
											"--format", "jwt")
		})

		It("can request opaque token", func() {
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusOK, opaqueTokenResponseJson),
				VerifyFormKV("client_id", "admin"),
				VerifyFormKV("client_secret", "adminsecret"),
				VerifyFormKV("grant_type", "password"),
				VerifyFormKV("token_format", "opaque"),
			))

			runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret",
				"--format", "opaque")
		})

		It("uses jwt format by default", func() {
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusOK, jwtTokenResponseJson),
				VerifyFormKV("client_id", "admin"),
				VerifyFormKV("client_secret", "adminsecret"),
				VerifyFormKV("grant_type", "password"),
				VerifyFormKV("token_format", "jwt"),
			))

			runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret")
		})

		It("displays error when unknown format is passed", func() {
			session := runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret",
				"--format", "bogus")
			Expect(session.Err).To(Say(`The token format "bogus" is unknown.`))
			Expect(session).To(Exit(1))
		})
	})

	Describe("Validations", func() {
		Describe("when called with no client id", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-password-token",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
			})
		})

		Describe("when called with no client secret", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-password-token", "admin",
					"-u", "woodstock",
					"-p", "secret")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `client_secret` must be specified."))
			})
		})

		Describe("when called with no username", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-password-token", "admin",
					"-s", "adminsecret",
					"-p", "secret")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `username` must be specified."))
			})
		})

		Describe("when called with no password", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("get-password-token", "admin",
					"-s", "adminsecret",
					"-u", "woodstock")

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `password` must be specified."))
			})
		})

		Describe("when no target was previously set", func() {
			BeforeEach(func() {
				config.WriteConfig(uaa.NewConfig())
			})

			It("tells the user to set a target", func() {
				session := runCommand("get-password-token", "admin",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret")
				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("You must set a target in order to use this command."))
			})
		})
	})
})
