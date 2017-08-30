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

var _ = Describe("CreateClient", func() {
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

		Describe("cloning another client configuration", func() {
			var shinyClient string

			BeforeEach(func() {
				shinyClient = `{
			  "scope" : [ "shiny.write" ],
			  "client_id" : "shiny",
			  "resource_ids" : [ ],
			  "authorized_grant_types" : [ "client_credentials", "authorization_code" ],
			  "redirect_uri" : [ "http://localhost:8080/*" ],
			  "authorities" : [ "shiny.write", "shiny.read" ],
			  "token_salt" : "",
			  "autoapprove" : ["true"],
			  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
			  "name" : "The Shiniest Client"
			}`
			})

			It("gets the specified client and creates a copy", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusOK, shinyClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["shiny.write"],"authorized_grant_types":["client_credentials","authorization_code"],"redirect_uri":["http://localhost:8080/*"],"authorities":["shiny.write","shiny.read"],"autoapprove":["true"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"The Shiniest Client"}`
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyRequest("POST", "/oauth/clients"),
					RespondWith(http.StatusOK, shinyCopy),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(shinyCopy),
				))

				session := runCommand("create-client",
					"shinycopy",
					"--clone", "shiny",
					"--client_secret", "secretsecret")

				Expect(session.Out).To(Say("The client shinycopy has been successfully created."))
				Expect(session).Should(Exit(0))
			})

			It("overrides other properties if specified", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusOK, shinyClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["foo.read"],"authorized_grant_types":["implicit"],"redirect_uri":["http://localhost:8001/*"],"authorities":["shiny.read"],"autoapprove":["true"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"foo client"}`
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyRequest("POST", "/oauth/clients"),
					RespondWith(http.StatusOK, shinyCopy),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(shinyCopy),
				))

				session := runCommand("create-client",
					"shinycopy",
					"--clone", "shiny",
					"--scope", "foo.read",
					"--authorized_grant_types", "implicit",
					"--client_secret", "secretsecret",
					"--display_name", "foo client",
					"--redirect_uri", "http://localhost:8001/*",
					"--authorities", "shiny.read",
				)

				Expect(session.Out).To(Say("The client shinycopy has been successfully created."))
				Expect(session).Should(Exit(0))
			})

			It("displays an error when the client cannot be found", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusNotFound, ""),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				session := runCommand("create-client",
					"shinycopy",
					"--clone", "shiny",
					"--client_secret", "secretsecret")

				Expect(session.Out).To(Say("The client shiny could not be found."))
				Expect(session).Should(Exit(1))
			})

			It("displays an error when the create fails", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusOK, shinyClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["shiny.write"],"authorized_grant_types":["client_credentials","authorization_code"],"redirect_uri":["http://localhost:8080/*"],"authorities":["shiny.write","shiny.read"],"autoapprove":["true"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"The Shiniest Client"}`
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyRequest("POST", "/oauth/clients"),
					RespondWith(http.StatusBadRequest, shinyCopy),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(shinyCopy),
				))

				session := runCommand("create-client",
					"shinycopy",
					"--clone", "shiny",
					"--client_secret", "secretsecret")

				Expect(session.Out).To(Say("An error occurred while creating the client."))
				Expect(session).Should(Exit(1))
			})

			It("still insists on a client_secret", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusOK, shinyClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["shiny.write"],"authorized_grant_types":["client_credentials","authorization_code"],"redirect_uri":["http://localhost:8080/*"],"authorities":["shiny.write","shiny.read"],"autoapprove":["true"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"The Shiniest Client"}`
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyRequest("POST", "/oauth/clients"),
					RespondWith(http.StatusOK, shinyCopy),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(shinyCopy),
				))

				session := runCommand("create-client", "shinycopy", "--clone", "shiny")

				Expect(session.Out).To(Say("Missing argument `client_secret` must be specified."))
				Expect(session).Should(Exit(1))
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
