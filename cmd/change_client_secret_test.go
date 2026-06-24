package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("ChangeClientSecret", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		// Create a client context with client_credentials grant type
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		ctx.GrantType = config.CLIENT_CREDENTIALS
		ctx.ClientId = "myclient"
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("successfully changes client secret", func() {
		server.RouteToHandler("PUT", "/oauth/clients/myclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/myclient/secret"),
			VerifyJSON(`{"oldSecret":"oldsecret","secret":"newsecret"}`),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			VerifyHeaderKV("Content-Type", "application/json"),
			RespondWith(http.StatusOK, `{"status":"ok","message":"Secret is updated"}`),
		))

		session := runCommand("change-client-secret", "--old_secret", "oldsecret", "--secret", "newsecret")

		Expect(session.Out).To(Say("The secret for client myclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("displays error when API request fails", func() {
		server.RouteToHandler("PUT", "/oauth/clients/myclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/myclient/secret"),
			VerifyJSON(`{"oldSecret":"wrongsecret","secret":"newsecret"}`),
			RespondWith(http.StatusBadRequest, `{"error":"invalid_secret","error_description":"The old secret is incorrect"}`),
		))

		session := runCommand("change-client-secret", "--old_secret", "wrongsecret", "--secret", "newsecret")

		Expect(session.Err).To(Say("The secret for client myclient was not updated."))
		Expect(session.Out).To(Say("Retry with --verbose for more information."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active target", func() {
		config.WriteConfig(config.NewConfig())
		session := runCommand("change-client-secret", "--old_secret", "old", "--secret", "new")

		Expect(session.Err).To(Say("You must set a target in order to use this command."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active context", func() {
		c := config.NewConfig()
		t := config.NewTarget()
		c.AddTarget(t)
		config.WriteConfig(c)
		session := runCommand("change-client-secret", "--old_secret", "old", "--secret", "new")

		Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when context is not client_credentials", func() {
		c := config.NewConfigWithServerURL(server.URL())
		// Create a password grant context (user context)
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		ctx.GrantType = config.PASSWORD
		ctx.Username = "testuser"
		c.AddContext(ctx)
		config.WriteConfig(c)

		session := runCommand("change-client-secret", "--old_secret", "old", "--secret", "new")

		Expect(session.Err).To(Say("You must have a client_credentials token in your context to perform this command."))
		Eventually(session).Should(Exit(1))
	})

	It("supports zone switching", func() {
		server.RouteToHandler("PUT", "/oauth/clients/myclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/myclient/secret"),
			VerifyJSON(`{"oldSecret":"oldsecret","secret":"newsecret"}`),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			RespondWith(http.StatusOK, `{"status":"ok","message":"Secret is updated"}`),
		))

		session := runCommand("change-client-secret", "--old_secret", "oldsecret", "--secret", "newsecret", "--zone", "twilight-zone")

		Expect(session.Out).To(Say("The secret for client myclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("shows verbose output when requested", func() {
		server.RouteToHandler("PUT", "/oauth/clients/myclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/myclient/secret"),
			RespondWith(http.StatusOK, `{"status":"ok","message":"Secret is updated"}`),
		))

		session := runCommand("change-client-secret", "--old_secret", "oldsecret", "--secret", "newsecret", "--verbose")

		Expect(session.Out).To(Say("The secret for client myclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("complains when no old secret is provided", func() {
		session := runCommand("change-client-secret", "--secret", "newsecret")

		Expect(session.Err).To(Say("Missing argument `old_secret` must be specified."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when no new secret is provided", func() {
		session := runCommand("change-client-secret", "--old_secret", "oldsecret")

		Expect(session.Err).To(Say("Missing argument `secret` must be specified."))
		Eventually(session).Should(Exit(1))
	})
})
