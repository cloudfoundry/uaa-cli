package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("GetPasswordToken", func() {

	var opaqueTokenResponseJson = `{
	  "access_token" : "bc4885d950854fed9a938e96b13ca519",
	  "refresh_token" : "abcd5d950854fed9a938e96b13ca519",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`
	var jwtTokenResponseJson = `{
	  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	  "refresh_token" : "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiJlMTQ0NTE3N2YyYmU0YzQ3Yjk4MmIzNzI1MzllN2NkNy1yIiwic3ViIjoiODkwZmY4MWItMzMyNC00NDRiLTgwNTAtNDRmNWVjOGQ3MDUzIiwic2NvcGUiOlsib3BlbmlkIiwidXNlcl9hdHRyaWJ1dGVzIiwic2NpbS53cml0ZSIsInNjaW0ucmVhZCJdLCJpYXQiOjE1MDUwNzk4MjMsImV4cCI6MTUwNzY3MTgyMywiY2lkIjoiamF1dGhjb2RlIiwiY2xpZW50X2lkIjoiamF1dGhjb2RlIiwiaXNzIjoiaHR0cHM6Ly91YWEudWFhLWFjY2VwdGFuY2UuY2YtYXBwLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImdyYW50X3R5cGUiOiJhdXRob3JpemF0aW9uX2NvZGUiLCJ1c2VyX25hbWUiOiJqaGFtb25AZ21haWwuY29tIiwib3JpZ2luIjoidWFhIiwidXNlcl9pZCI6Ijg5MGZmODFiLTMzMjQtNDQ0Yi04MDUwLTQ0ZjVlYzhkNzA1MyIsInJldl9zaWciOiI1NjFiNGRjMCIsImF1ZCI6WyJzY2ltIiwiamF1dGhjb2RlIiwib3BlbmlkIl19.hxTIL6pbybnpXwioYepdAEWHHwBB6hqJJjWW4atZJ4jeg1ZZCe6KKPM0xEo43mwLfuqcPim7Y7GAJFiJfcM9iqilzCLWAYvQi4aeliOgsYRrWpExYXSQ76bnJ584co7a4xSbxk6W_uXFGbcgBqJaOMlJ_TbIqtFqrvsf3CzGcDy7Mnir8caQru2tEr8Zlz4zuZImj6-FJ4AQkYW1RwXD2m94I2ZoCrv2eP-AVQjgbCDHgoN2jv9-Y1eyLagVqOXBgcd9KOQFqvm4D6ker3_grbq5VmZ-8QxwbsFZ5Sl6Q-Bk7y00nhQccLIKmNqECoAb520Zwm5OhcJERbq9jgTz9Q",
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

		Context("not successful", func() {
			It("shows extra output about the request on error", func() {
				server.RouteToHandler("POST", "/oauth/token",
					RespondWith(http.StatusBadRequest, "garbage response"),
				)

				session := runCommand("get-password-token",
					"admin",
					"-s", "adminsecret",
					"-u", "woodstock",
					"-p", "secret",
					"--verbose")

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("Unable to retrieve token"))
			})
		})

		Describe("when successful", func() {
			BeforeEach(func() {
				config.WriteConfig(c)
				server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
					RespondWith(http.StatusOK, jwtTokenResponseJson, http.Header{
						"Content-Type": []string{"application/json"},
					}),
					VerifyFormKV("client_id", "admin"),
					//base64 <<< 'admin:adminsecret'
					VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
					VerifyFormKV("grant_type", "password"),
				))
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

				Expect(config.ReadConfig().GetActiveContext().Token.AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
				Expect(config.ReadConfig().GetActiveContext().Token.RefreshToken).To(Equal("eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiJlMTQ0NTE3N2YyYmU0YzQ3Yjk4MmIzNzI1MzllN2NkNy1yIiwic3ViIjoiODkwZmY4MWItMzMyNC00NDRiLTgwNTAtNDRmNWVjOGQ3MDUzIiwic2NvcGUiOlsib3BlbmlkIiwidXNlcl9hdHRyaWJ1dGVzIiwic2NpbS53cml0ZSIsInNjaW0ucmVhZCJdLCJpYXQiOjE1MDUwNzk4MjMsImV4cCI6MTUwNzY3MTgyMywiY2lkIjoiamF1dGhjb2RlIiwiY2xpZW50X2lkIjoiamF1dGhjb2RlIiwiaXNzIjoiaHR0cHM6Ly91YWEudWFhLWFjY2VwdGFuY2UuY2YtYXBwLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImdyYW50X3R5cGUiOiJhdXRob3JpemF0aW9uX2NvZGUiLCJ1c2VyX25hbWUiOiJqaGFtb25AZ21haWwuY29tIiwib3JpZ2luIjoidWFhIiwidXNlcl9pZCI6Ijg5MGZmODFiLTMzMjQtNDQ0Yi04MDUwLTQ0ZjVlYzhkNzA1MyIsInJldl9zaWciOiI1NjFiNGRjMCIsImF1ZCI6WyJzY2ltIiwiamF1dGhjb2RlIiwib3BlbmlkIl19.hxTIL6pbybnpXwioYepdAEWHHwBB6hqJJjWW4atZJ4jeg1ZZCe6KKPM0xEo43mwLfuqcPim7Y7GAJFiJfcM9iqilzCLWAYvQi4aeliOgsYRrWpExYXSQ76bnJ584co7a4xSbxk6W_uXFGbcgBqJaOMlJ_TbIqtFqrvsf3CzGcDy7Mnir8caQru2tEr8Zlz4zuZImj6-FJ4AQkYW1RwXD2m94I2ZoCrv2eP-AVQjgbCDHgoN2jv9-Y1eyLagVqOXBgcd9KOQFqvm4D6ker3_grbq5VmZ-8QxwbsFZ5Sl6Q-Bk7y00nhQccLIKmNqECoAb520Zwm5OhcJERbq9jgTz9Q"))
				Expect(config.ReadConfig().GetActiveContext().Token.TokenType).To(Equal("bearer"))
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
				VerifyFormKV("client_id", "admin"),
				//base64 <<< 'admin:adminsecret'
				VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
				VerifyFormKV("grant_type", "password"),
			))
		})

		It("displays help to the user", func() {
			session := runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret")

			Eventually(session).Should(Exit(1))
			expectedOutput := fmt.Sprintf("An error occurred while calling %s/oauth/token", server.URL())
			Eventually(session.Err).Should(Say(expectedOutput))

			var unauthorizedErrorMsg bytes.Buffer
			_ = json.Indent(&unauthorizedErrorMsg, []byte(`{"error":"unauthorized","error_description":"Bad credentials"}`), "", "  ")
			Eventually(session.Err).Should(Say(unauthorizedErrorMsg.String()))
		})

		It("does not update the previously saved context", func() {
			runCommand("get-password-token", "admin",
				"-s", "adminsecret",
				"-u", "woodstock",
				"-p", "secret")
			Expect(config.ReadConfig().GetActiveContext().Token.AccessToken).To(Equal("old-token"))
		})
	})

	Describe("configuring token format", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("can request jwt token", func() {
			server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusOK, jwtTokenResponseJson, contentTypeJson),
				VerifyFormKV("client_id", "admin"),
				//base64 <<< 'admin:adminsecret'
				VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
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
				RespondWith(http.StatusOK, opaqueTokenResponseJson, contentTypeJson),
				VerifyFormKV("client_id", "admin"),
				//base64 <<< 'admin:adminsecret'
				VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
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
				RespondWith(http.StatusOK, jwtTokenResponseJson, contentTypeJson),
				//base64 <<< 'admin:adminsecret'
				VerifyHeaderKV("Authorization", "Basic YWRtaW46YWRtaW5zZWNyZXQ="),
				VerifyFormKV("client_id", "admin"),
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
				c := config.NewConfigWithServerURL("http://localhost")
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
			It("succeeds", func() {
				c := config.NewConfigWithServerURL(server.URL())
				config.WriteConfig(c)
				server.RouteToHandler("POST", "/oauth/token",
					RespondWith(http.StatusOK, jwtTokenResponseJson, contentTypeJson),
				)

				session := runCommand("get-password-token", "admin",
					"-u", "woodstock",
					"-p", "secret",
					"--verbose")

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("Access token successfully fetched and added to context."))
			})
		})

		Describe("when called with no username", func() {
			It("displays help and does not panic", func() {
				c := config.NewConfigWithServerURL("http://localhost")
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
				c := config.NewConfigWithServerURL("http://localhost")
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
				config.WriteConfig(config.NewConfig())
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
