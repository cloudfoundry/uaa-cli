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

var _ = Describe("GetClient", func() {
	const GetClientResponseJson string = `{
		  "scope" : [ "clients.read", "clients.write" ],
		  "client_id" : "clientid",
		  "resource_ids" : [ "none" ],
		  "authorized_grant_types" : [ "client_credentials" ],
		  "redirect_uri" : [ "http://ant.path.wildcard/**/passback/*", "http://test1.com" ],
		  "authorities" : [ "clients.read", "clients.write" ],
		  "token_salt" : "1SztLL",
		  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
		  "name" : "My Client Name",
		  "lastModified" : 1502816030525,
		  "required_user_groups" : [ "cloud_controller.admin" ]
		}`

	Describe("zone switching support", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("adds the zone switching header", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				CombineHandlers(
					VerifyRequest("GET", "/oauth/clients/clientid"),
					RespondWith(http.StatusOK, GetClientResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
				),
			)

			session := runCommand("get-client", "clientid", "--zone", "twilight-zone")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("shows the client configuration response", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid",
				CombineHandlers(
					RespondWith(http.StatusOK, GetClientResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
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

			Expect(session.Err).To(Say("An error occurred while calling " + server.URL() + "/oauth/clients/clientid"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no client_id is supplied", func() {
		It("displays an error message to the user", func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("sometoken"))
			config.WriteConfig(c)
			session := runCommand("get-client")

			Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("get-client", "clientid")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must set a target in order to use this command."))
		})
	})
})
