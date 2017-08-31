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

var _ = Describe("UpdateClient", func() {
	notifierClient := `{
	  "scope" : [ "notifications.write" ],
	  "client_id" : "notifier",
	  "client_secret" : "secret",
	  "resource_ids" : [ ],
	  "authorized_grant_types" : [ "client_credentials", "authorization_code" ],
	  "redirect_uri" : [ "http://localhost:8080/*" ],
	  "authorities" : [ "notifications.write", "notifications.read" ],
	  "token_salt" : "",
	  "autoapprove" : ["true"],
	  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
	  "name" : "Notifier Client"
	}`

	var c uaa.Config
	var ctx uaa.UaaContext

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c = uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.UaaContext{AccessToken: "access_token"})
			config.WriteConfig(c)
			ctx = c.GetActiveContext()
		})

		Describe("when the --trace option is used", func() {
			It("shows extra output about the request on success", func() {
				server.RouteToHandler("PUT", "/oauth/clients/notifier",
					RespondWith(http.StatusOK, notifierClient),
				)

				session := runCommand("update-client",
					"notifier",
					"--authorized_grant_types", "client_credentials,authorization_code",
					"--authorities", "notifications.write",
					"--trace")

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("PUT " + server.URL() + "/oauth/clients/notifier"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("200 OK"))
			})

			It("shows extra output about the request on error", func() {
				server.RouteToHandler("PUT", "/oauth/clients/notifier",
					RespondWith(http.StatusBadRequest, "garbage response"),
				)

				session := runCommand("update-client",
					"notifier",
					"--authorized_grant_types", "client_credentials,authorization_code",
					"--authorities", "notifications.write",
					"--trace")

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("PUT " + server.URL() + "/oauth/clients/notifier"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("400 Bad Request"))
				Expect(session.Out).To(Say("garbage response"))
			})
		})

		Describe("using the --zone flag", func() {
			BeforeEach(func() {
				c := uaa.NewConfigWithServerURL(server.URL())
				ctx := uaa.UaaContext{AccessToken: "access_token"}
				c.AddContext(ctx)
				config.WriteConfig(c)
			})

			It("adds a zone-switching header to the request", func() {
				server.RouteToHandler("PUT", "/oauth/clients/notifier", CombineHandlers(
					VerifyRequest("PUT", "/oauth/clients/notifier"),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				runCommand("update-client", "notifier", "--zone", "twilight-zone")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})

		Describe("when successful", func() {
			BeforeEach(func() {
				config.WriteConfig(c)
			})

			It("displays a success message and prints the updated configuration", func() {
				server.RouteToHandler("PUT", "/oauth/clients/notifier", CombineHandlers(
					RespondWith(http.StatusOK, notifierClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(`{"client_id":"notifier","authorized_grant_types":["client_credentials"],"authorities":["notifications.write","notifications.read"]}`),
				))

				session := runCommand("update-client",
					"notifier",
					"--authorized_grant_types", "client_credentials",
					"--authorities", "notifications.write,notifications.read")

				Eventually(session).Should(Say("The client notifier has been successfully updated."))
				Eventually(session).Should(Say(`"authorized_grant_types"`))
				Eventually(session).Should(Say(`"client_credentials"`))
				Eventually(session).Should(Exit(0))
			})

			It("knows about many flags", func() {
				server.RouteToHandler("PUT", "/oauth/clients/notifier", CombineHandlers(
					RespondWith(http.StatusOK, notifierClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(`{ "scope" : [ "notifications.write" ], "client_id" : "notifier", "authorized_grant_types" : [ "client_credentials" ], "redirect_uri" : [ "http://localhost:8080/*" ], "authorities" : [ "notifications.write", "notifications.read" ], "autoapprove" : ["scim.write", "scim.read"], "name" : "Display name" }`),
				))

				session := runCommand("update-client",
					"notifier",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
					"--display_name", "Display name",
					"--access_token_validity", "3600",
					"--refresh_token_validity", "4500",
					"--autoapprove", "scim.write,scim.read",
				)

				Eventually(session).Should(Say("The client notifier has been successfully updated."))
				Eventually(session).Should(Say(`"authorized_grant_types"`))
				Eventually(session).Should(Say(`"client_credentials"`))
				Eventually(session).Should(Exit(0))
			})

			It("does not handle changing the client secret", func() {
				session := runCommand("update-client",
					"notifier",
					"--client_secret", "newsecret",
				)

				Eventually(session).Should(Say(`Client not updated. Please see "uaa set-client-secret -h" to learn more about changing client secrets.`))
				Eventually(session).Should(Exit(1))
			})
		})
	})

	Describe("when the client update fails", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			c.AddContext(uaa.UaaContext{AccessToken: "old-token"})
			config.WriteConfig(c)
			server.RouteToHandler("PUT", "/oauth/clients/notifier", CombineHandlers(
				RespondWith(http.StatusUnauthorized, `{"error":"unauthorized","error_description":"Bad credentials"}`),
			))
		})

		It("displays help to the user", func() {
			session := runCommand("update-client",
				"notifier",
				"--authorized_grant_types", "client_credentials",
				"--scope", "notifications.write",
				"--redirect_uri", "http://localhost:8080/*",
				"--authorities", "notifications.write,notifications.read",
			)

			Eventually(session).Should(Say("An error occurred while updating the client."))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("Validations", func() {
		Describe("when called with no client id", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("update-client",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("Missing argument `client_id` must be specified."))
			})
		})

		Describe("when no target was previously set", func() {
			BeforeEach(func() {
				config.WriteConfig(uaa.NewConfig())
			})

			It("tells the user to set a target", func() {
				session := runCommand("update-client",
					"notifier",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)
				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("You must set a target in order to use this command."))
			})
		})
	})
})
