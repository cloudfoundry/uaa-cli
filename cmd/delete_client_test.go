package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
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

	Describe("--verbose flag support", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("access_token"))
			config.WriteConfig(c)
		})

		It("shows extra output about the request on success", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/clientid",
				RespondWith(http.StatusOK, DeleteClientResponseJson),
			)

			session := runCommand("delete-client", "clientid", "--verbose")

			Expect(session.Out).To(Say("DELETE /oauth/clients/clientid"))
			Expect(session.Out).To(Say("Accept: application/json"))
			Expect(session.Out).To(Say("200 OK"))
			Eventually(session).Should(Exit(0))
		})

		It("shows extra output about the request on error", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/clientid",
				RespondWith(http.StatusBadRequest, "garbage response"),
			)

			session := runCommand("delete-client", "clientid", "--verbose")

			Eventually(session).Should(Exit(1))
			Expect(session.Out).To(Say("DELETE /oauth/clients/clientid"))
			Expect(session.Out).To(Say("Accept: application/json"))
			Expect(session.Out).To(Say("400 Bad Request"))
			Expect(session.Out).To(Say("garbage response"))
		})
	})

	Describe("using the --zone flag", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("access_token"))
			config.WriteConfig(c)
		})

		It("adds a zone-switching header to the delete request", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/notifier", CombineHandlers(
				VerifyRequest("DELETE", "/oauth/clients/notifier"),
				VerifyHeaderKV("Authorization", "bearer access_token"),
				VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
			))

			runCommand("delete-client", "notifier", "--zone", "twilight-zone")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("access_token"))
			config.WriteConfig(c)
		})

		It("shows the client configuration response", func() {
			server.RouteToHandler("DELETE", "/oauth/clients/clientid",
				CombineHandlers(
					RespondWith(http.StatusOK, DeleteClientResponseJson),
					VerifyHeaderKV("Authorization", "bearer access_token"),
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

			Expect(session.Err).To(Say("An unknown error occurred while calling " + server.URL() + "/oauth/clients/clientid"))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no client_id is supplied", func() {
		It("displays and error message to the user", func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("access_token"))
			config.WriteConfig(c)
			session := runCommand("delete-client")

			Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
			Eventually(session).Should(Exit(1))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.Config{}
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
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")
			cfg.AddContext(uaa.NewContextWithToken("access_token"))

			err := cmd.DeleteClientValidations(cfg, []string{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_id` must be specified."))
		})

		It("requires token in context to have been set", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")
			err := cmd.DeleteClientValidations(cfg, []string{"clientid"})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(cmd.MISSING_CONTEXT))
		})
	})
})
