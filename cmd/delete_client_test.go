package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("DeleteClient", func() {
	const DeleteClientResponseJson string = `{
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

	Describe("using the --zone flag", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("adds a zone-switching header to the delete request", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/notifier", CombineHandlers(
				VerifyRequest("DELETE", "/oauth/clients/notifier"),
				VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
				VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			))

			runCommand("delete-client", "notifier", "--zone", "twilight-zone")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("shows the client configuration response", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/clientid",
				CombineHandlers(
					RespondWith(http.StatusOK, DeleteClientResponseJson, contentTypeJson),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
				),
			)

			session := runCommand("delete-client", "clientid")

			Expect(session.Out).To(Say("Successfully deleted client clientid."))
			Eventually(session).Should(Exit(0))
		})

		It("handles request errors", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/clientid",
				RespondWith(http.StatusNotFound, ""),
			)

			session := runCommand("delete-client", "clientid")

			Expect(session.Err).To(Say("An error occurred while calling " + server.URL() + "/oauth/clients/clientid"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no client_id is supplied", func() {
		It("displays and error message to the user", func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
			session := runCommand("delete-client")

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
			session := runCommand("delete-client", "clientid")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must set a target in order to use this command."))
		})
	})

	Describe("Validations", func() {
		It("requires a client id", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")
			cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))

			err := cmd.DeleteClientValidations(cfg, []string{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_id` must be specified."))
		})

		It("requires token in context to have been set", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")
			err := cmd.DeleteClientValidations(cfg, []string{"clientid"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(cli.MISSING_CONTEXT))
		})
	})
})
