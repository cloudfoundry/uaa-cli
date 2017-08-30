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

var _ = Describe("GetClient", func() {
	const GetClientResponseJson string = `{
		  "scope" : [ "clients.read", "clients.write" ],
		  "client_id" : "clientid",
		  "resource_ids" : [ "none" ],
		  "authorized_grant_types" : [ "client_credentials" ],
		  "redirect_uri" : [ "http://ant.path.wildcard/**/passback/*", "http://test1.com" ],
		  "autoapprove" : [ "true" ],
		  "authorities" : [ "clients.read", "clients.write" ],
		  "token_salt" : "1SztLL",
		  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
		  "name" : "My Client Name",
		  "lastModified" : 1502816030525,
		  "required_user_groups" : [ "cloud_controller.admin" ]
		}`

	Describe("--trace flag support", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			ctx := uaa.UaaContext{AccessToken: "access_token"}
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("shows extra output about the request on success", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				RespondWith(http.StatusOK, GetClientResponseJson),
			)

			session := runCommand("get-client", "clientid", "--trace")

			Expect(session.Out).To(Say("GET " + server.URL() + "/oauth/clients/clientid"))
			Expect(session.Out).To(Say("Accept: application/json"))
			Expect(session.Out).To(Say("200 OK"))
			Eventually(session).Should(Exit(0))
		})

		It("shows extra output about the request on error", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				RespondWith(http.StatusBadRequest, "garbage response"),
			)

			session := runCommand("get-client", "clientid", "--trace")

			Eventually(session).Should(Exit(1))
			Expect(session.Out).To(Say("GET " + server.URL() + "/oauth/clients/clientid"))
			Expect(session.Out).To(Say("Accept: application/json"))
			Expect(session.Out).To(Say("400 Bad Request"))
			Expect(session.Out).To(Say("garbage response"))
		})
	})

	Describe("zone switching support", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			ctx := uaa.UaaContext{AccessToken: "access_token"}
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("adds the zone switching header", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/clientid"),
					RespondWith(http.StatusOK, GetClientResponseJson),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				),
			)

			session := runCommand("get-client", "clientid", "--zone", "twilight-zone")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			ctx := uaa.UaaContext{AccessToken: "access_token"}
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("shows the client configuration response", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				CombineHandlers(
					RespondWith(http.StatusOK, GetClientResponseJson),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				),
			)

			session := runCommand("get-client", "clientid")

			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(GetClientResponseJson))
			Eventually(session).Should(Exit(0))
		})

		It("handles request errors", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				RespondWith(http.StatusNotFound, ""),
			)

			session := runCommand("get-client", "clientid")

			Expect(session.Out).To(Say("An unknown error occurred while calling " + server.URL() + "/oauth/clients/clientid"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no client_id is supplied", func() {
		It("displays and error message to the user", func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
			session := runCommand("get-client")

			Expect(session.Out).To(Say("Missing argument `client_id` must be specified."))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("get-client", "clientid")

			Eventually(session).Should(Exit(1))
			Expect(session.Out).To(Say("You must set a target in order to use this command."))
		})
	})
})
