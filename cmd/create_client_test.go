package cmd_test

import (
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("CreateClient", func() {
	//	clientJson := `{
	//  "scope" : [ "clients.read", "clients.write" ],
	//  "client_id" : "myclient",
	//  "client_secret" : "secret",
	//  "resource_ids" : [ ],
	//  "authorized_grant_types" : [ "client_credentials", "authorization_code" ],
	//  "redirect_uri" : [ "http://test1.com", "http://ant.path.wildcard/**/passback/*" ],
	//  "authorities" : [ "clients.read", "clients.write" ],
	//  "token_salt" : "zCSAYx",
	//  "autoapprove" : true,
	//  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
	//  "name" : "My Client Name"
	//}`

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
				server.RouteToHandler("POST", "/oauth/clients",
					RespondWith(http.StatusOK, notifierClient),
				)

				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials,authorization_code",
					"--authorities", "notifications.write",
					"--trace")

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("POST " + server.URL() + "/oauth/clients"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("200 OK"))
			})

			It("shows extra output about the request on error", func() {
				server.RouteToHandler("POST", "/oauth/clients",
					RespondWith(http.StatusBadRequest, "garbage response"),
				)

				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials,authorization_code",
					"--authorities", "notifications.write",
					"--trace")

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("POST " + server.URL() + "/oauth/clients"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("400 Bad Request"))
				Expect(session.Out).To(Say("garbage response"))
			})
		})

		Describe("when successful", func() {
			BeforeEach(func() {
				config.WriteConfig(c)
			})

			It("displays a success message and prints the created configuration", func() {
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					RespondWith(http.StatusOK, notifierClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(`{"client_id":"notifier","client_secret":"secret","authorized_grant_types":["client_credentials"],"authorities":["notifications.write","notifications.read"]}`),
				))

				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--authorities", "notifications.write,notifications.read")

				Eventually(session).Should(Say("The client notifier has been successfully created."))
				Eventually(session).Should(Say(`"authorized_grant_types"`))
				Eventually(session).Should(Say(`"client_credentials"`))
				Eventually(session).Should(Exit(0))
			})

			It("knows about many flags", func() {
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					RespondWith(http.StatusOK, notifierClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(`{ "scope" : [ "notifications.write" ], "client_id" : "notifier", "client_secret" : "secret", "authorized_grant_types" : [ "client_credentials" ], "redirect_uri" : [ "http://localhost:8080/*" ], "authorities" : [ "notifications.write", "notifications.read" ], "autoapprove" : ["scim.write", "scim.read"], "name" : "Display name" }`),
				))

				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
					"--display_name", "Display name",
					"--access_token_validity", "3600",
					"--refresh_token_validity", "4500",
					"--autoapprove", "scim.write,scim.read",
				)

				Eventually(session).Should(Say("The client notifier has been successfully created."))
				Eventually(session).Should(Say(`"authorized_grant_types"`))
				Eventually(session).Should(Say(`"client_credentials"`))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Describe("when the client creation fails", func() {
		BeforeEach(func() {
			c := uaa.NewConfig()
			c.AddContext(uaa.UaaContext{AccessToken: "old-token"})
			config.WriteConfig(c)
			server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
				RespondWith(http.StatusUnauthorized, `{"error":"unauthorized","error_description":"Bad credentials"}`),
			))
		})

		It("displays help to the user", func() {
			session := runCommand("create-client",
				"notifier",
				"--client_secret", "secret",
				"--authorized_grant_types", "client_credentials",
				"--scope", "notifications.write",
				"--redirect_uri", "http://localhost:8080/*",
				"--authorities", "notifications.write,notifications.read",
			)

			Eventually(session).Should(Say("An error occurred while creating the client."))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("Validations", func() {
		Describe("when called with no client id", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("create-client",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("Missing argument `client_id` must be specified."))
			})
		})

		Describe("when called with no authorized_grant_type", func() {
			It("displays help and does not panic", func() {
				c := uaa.NewConfigWithServerURL("http://localhost")
				config.WriteConfig(c)
				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)

				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("Missing argument `authorized_grant_types` must be specified."))
			})
		})

		Describe("when no target was previously set", func() {
			BeforeEach(func() {
				config.WriteConfig(uaa.NewConfig())
			})

			It("tells the user to set a target", func() {
				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
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
