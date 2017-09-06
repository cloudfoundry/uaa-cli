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

var _ = Describe("ListClients", func() {
	const ClientsListResponseJsonPage1 = `{
		  "resources" : [{
			"client_id" : "client1"
		  },
		  {
			"client_id" : "client2"
		  }],
		  "startIndex" : 1,
		  "itemsPerPage" : 2,
		  "totalResults" : 2,
		  "schemas" : [ "http://cloudfoundry.org/schema/scim/oauth-clients-1.0" ]
		}`

	Describe("--trace flag support", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			ctx := uaa.UaaContext{AccessToken: "access_token"}
			c.AddContext(ctx)
			config.WriteConfig(c)
		})

		It("shows extra output about the request on success", func() {
			server.RouteToHandler("GET", "/oauth/clients",
				RespondWith(http.StatusOK, ClientsListResponseJsonPage1),
			)

			session := runCommand("list-clients", "--trace")

			Expect(session.Out).To(Say("GET " + server.URL() + "/oauth/clients"))
			Expect(session.Out).To(Say("Accept: application/json"))
			Expect(session.Out).To(Say("200 OK"))
			Eventually(session).Should(Exit(0))
		})

		It("shows extra output about the request on error", func() {
			server.RouteToHandler("GET", "/oauth/clients",
				RespondWith(http.StatusBadRequest, "garbage response"),
			)

			session := runCommand("list-clients", "--trace")

			Eventually(session).Should(Exit(1))
			Expect(session.Out).To(Say("GET " + server.URL() + "/oauth/clients"))
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
			server.RouteToHandler("GET", "/oauth/clients",
				CombineHandlers(
					VerifyRequest("GET", "/oauth/clients"),
					RespondWith(http.StatusOK, ClientsListResponseJsonPage1),
					VerifyHeaderKV("Authorization", "bearer access_token"),
					VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
				),
			)

			session := runCommand("list-clients", "--trace", "--zone", "twilight-zone")
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
			server.RouteToHandler("GET", "/oauth/clients",
				CombineHandlers(
					RespondWith(http.StatusOK, ClientsListResponseJsonPage1),
					VerifyHeaderKV("Authorization", "bearer access_token"),
				),
			)

			session := runCommand("list-clients")

			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(`[{ "client_id" : "client1" }, { "client_id" : "client2" }]`))
			Eventually(session).Should(Exit(0))
		})

		It("handles request errors", func() {
			server.RouteToHandler("GET", "/oauth/clients",
				RespondWith(http.StatusNotFound, ""),
			)

			session := runCommand("list-clients")

			Expect(session.Err).To(Say("An unknown error occurred while calling " + server.URL() + "/oauth/clients"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("list-clients")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must set a target in order to use this command."))
		})
	})

	Describe("when no context was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("tells the user to set a context", func() {
			session := runCommand("list-clients")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
		})
	})
})
