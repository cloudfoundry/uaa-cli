package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
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

	Describe("zone switching support", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("adds the zone switching header", func() {
			server.RouteToHandler("GET", "/oauth/clients",
				CombineHandlers(
					VerifyRequest("GET", "/oauth/clients"),
					RespondWith(http.StatusOK, ClientsListResponseJsonPage1, contentTypeJson),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
				),
			)

			session := runCommand("list-clients", "--verbose", "--zone", "twilight-zone")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		Context("when the request to /oauth/clients succeeds", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/oauth/clients",
					CombineHandlers(
						RespondWith(http.StatusOK, ClientsListResponseJsonPage1, contentTypeJson),
						VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					),
				)
			})

			It("shows the client configuration response", func() {
				session := runCommand("list-clients")

				outputBytes := session.Out.Contents()
				Expect(outputBytes).To(MatchJSON(`[{ "client_id" : "client1" }, { "client_id" : "client2" }]`))
				Eventually(session).Should(Exit(0))
			})

			It("shows verbose output when passed --verbose", func() {
				session := runCommand("list-clients", "--verbose")

				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("GET " + "/oauth/clients"))
				Expect(session.Out).To(Say("Accept: application/json"))
				Expect(session.Out).To(Say("200 OK"))
			})

		})

		Context("when the request to /oauth/clients fails", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/oauth/clients",
					RespondWith(http.StatusNotFound, ""),
				)
			})

			It("handles request errors", func() {
				session := runCommand("list-clients")

				Expect(session.Err).To(Say("An error occurred while calling " + server.URL() + "/oauth/clients"))
				Eventually(session).Should(Exit(1))
			})
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.Config{}
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
			c := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("tells the user to set a context", func() {
			session := runCommand("list-clients")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
		})
	})
})
