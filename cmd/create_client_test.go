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

	BeforeEach(func() {
		c = uaa.NewConfigWithServerURL(server.URL())
		c.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(c)
		ctx = c.GetActiveContext()
	})

	Describe("and a target was previously set", func() {
		Describe("when the --verbose option is used", func() {
			It("shows extra output about the request on success", func() {
				server.RouteToHandler("POST", "/oauth/clients",
					RespondWith(http.StatusOK, notifierClient),
				)

				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--authorities", "notifications.write",
					"--verbose")

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
					"--authorized_grant_types", "client_credentials",
					"--authorities", "notifications.write",
					"--verbose")

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
					VerifyJSON(`{ "scope" : [ "notifications.write" ], "client_id" : "notifier", "client_secret" : "secret", "authorized_grant_types" : [ "client_credentials" ], "redirect_uri" : [ "http://localhost:8080/*" ], "authorities" : [ "notifications.write", "notifications.read" ], "name" : "Display name", "access_token_validity": 3600, "refresh_token_validity": 4500 }`),
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
				)

				Eventually(session).Should(Say("The client notifier has been successfully created."))
				Eventually(session).Should(Say(`"authorized_grant_types"`))
				Eventually(session).Should(Say(`"client_credentials"`))
				Eventually(session).Should(Exit(0))
			})
		})

		Describe("using the --zone flag", func() {
			It("adds a zone-switching header to the create request", func() {
				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				))

				runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--authorities", "notifications.write,notifications.read",
					"--zone", "twilight-zone")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
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

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["shiny.write"],"authorized_grant_types":["client_credentials","authorization_code"],"redirect_uri":["http://localhost:8080/*"],"authorities":["shiny.write","shiny.read"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"The Shiniest Client"}`
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

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "scope":["foo.read"],"authorized_grant_types":["implicit"],"redirect_uri":["http://localhost:8001/*"],"authorities":["shiny.read"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"foo client", "access_token_validity": 3600, "refresh_token_validity": 4500}`
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
					"--access_token_validity", "3600",
					"--refresh_token_validity", "4500",
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

				Expect(session.Err).To(Say("The client shiny could not be found."))
				Expect(session).Should(Exit(1))
			})

			It("displays an error when the create fails", func() {
				server.RouteToHandler("GET", "/oauth/clients/shiny", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/shiny"),
					RespondWith(http.StatusOK, shinyClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				shinyCopy := `{"client_id":"shinycopy","client_secret":"secretsecret", "authorities":["shiny.write","shiny.read"], "scope":["shiny.write"],"authorized_grant_types":["client_credentials","authorization_code"],"redirect_uri":["http://localhost:8080/*"],"allowedproviders":["uaa","ldap","my-saml-provider"],"name":"The Shiniest Client"}`
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

				Expect(session.Err).To(Say("An unknown error occurred while calling"))
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

				Expect(session.Err).To(Say("client_secret must be specified"))
				Expect(session).Should(Exit(1))
			})

			It("does not require client_secret when cloning implicit grant type", func() {
				var implicitClient = `{
				  "scope" : [ "implicit.write" ],
				  "client_id" : "myImplicitClient",
				  "resource_ids" : [ ],
				  "authorized_grant_types" : [ "implicit" ],
				  "redirect_uri" : [ "http://localhost:8080/*" ],
				  "authorities" : [ "implicit.write", "implicit.read" ],
				  "token_salt" : "",
				  "autoapprove" : ["true"],
				  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
				  "name" : "Implicit Client"
				}`

				server.RouteToHandler("GET", "/oauth/clients/myImplicitClient", CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/myImplicitClient"),
					RespondWith(http.StatusOK, implicitClient),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				implicitCopy := `{ "scope" : [ "implicit.write" ], "client_id" : "implicitcopy", "authorized_grant_types" : [ "implicit" ], "redirect_uri" : [ "http://localhost:8080/*" ], "authorities" : [ "implicit.write", "implicit.read" ], "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ], "name" : "Implicit Client" }`

				server.RouteToHandler("POST", "/oauth/clients", CombineHandlers(
					VerifyRequest("POST", "/oauth/clients"),
					RespondWith(http.StatusOK, implicitCopy),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyJSON(implicitCopy),
				))

				session := runCommand("create-client", "implicitcopy", "--clone", "myImplicitClient")

				Expect(session).Should(Exit(0))
			})
		})
	})

	Describe("when the client creation fails", func() {
		BeforeEach(func() {
			c = uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("old-token"))
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

			Eventually(session.Err).Should(Say("An unknown error occurred while calling"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("Validations", func() {
		Describe("when called with no client id", func() {
			It("displays help and does not panic", func() {
				session := runCommand("create-client",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
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
				Expect(session.Err).To(Say("You must set a target in order to use this command."))
			})
		})

		Describe("when no context with token was previously set", func() {
			BeforeEach(func() {
				c := uaa.NewConfig()
				t := uaa.NewTarget()
				c.AddTarget(t)
				config.WriteConfig(c)
			})

			It("tells the user to get a token", func() {
				session := runCommand("create-client",
					"notifier",
					"--client_secret", "secret",
					"--authorized_grant_types", "client_credentials",
					"--scope", "notifications.write",
					"--redirect_uri", "http://localhost:8080/*",
					"--authorities", "notifications.write,notifications.read",
				)
				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
			})
		})
	})
})
