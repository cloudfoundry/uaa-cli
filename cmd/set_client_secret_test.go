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

var _ = Describe("SetClientSecret", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(c)
	})

	It("can update client secrets", func() {
		server.RouteToHandler("PUT", "/oauth/clients/shinyclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/shinyclient/secret"),
			VerifyJSON(`{"clientId":"shinyclient","secret":"shinysecret"}`),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			RespondWith(http.StatusOK, "{}"),
		),
		)

		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Out).To(Say("The secret for client shinyclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("displays error when request fails", func() {
		server.RouteToHandler("PUT", "/oauth/clients/shinyclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/shinyclient/secret"),
			VerifyJSON(`{"clientId":"shinyclient","secret":"shinysecret"}`),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			RespondWith(http.StatusUnauthorized, "{}"),
		),
		)

		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Err).To(Say("The secret for client shinyclient was not updated."))
		Expect(session.Out).To(Say("Retry with --verbose for more information."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active target", func() {
		config.WriteConfig(config.NewConfig())
		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Err).To(Say("You must set a target in order to use this command."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active context", func() {
		c := config.NewConfig()
		t := config.NewTarget()
		c.AddTarget(t)
		config.WriteConfig(c)
		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
		Eventually(session).Should(Exit(1))
	})

	It("supports zone switching", func() {
		server.RouteToHandler("PUT", "/oauth/clients/shinyclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/shinyclient/secret"),
			VerifyJSON(`{"clientId":"shinyclient","secret":"shinysecret"}`),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			RespondWith(http.StatusOK, "{}"),
		),
		)

		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret", "--zone", "twilight-zone")

		Expect(session.Out).To(Say("The secret for client shinyclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("complains when no secret is provided", func() {
		session := runCommand("set-client-secret", "shinyclient")

		Expect(session.Err).To(Say("Missing argument `client_secret` must be specified."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when no clientid is provided", func() {
		session := runCommand("set-client-secret", "-s", "shinysecret")

		Expect(session.Err).To(Say("Missing argument `client_id` must be specified."))
		Eventually(session).Should(Exit(1))
	})
})
