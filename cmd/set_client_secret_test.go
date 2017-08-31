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

var _ = Describe("SetClientSecret", func() {
	BeforeEach(func() {
		c := uaa.NewConfigWithServerURL(server.URL())
		ctx := uaa.UaaContext{AccessToken: "access_token"}
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("can update client secrets", func() {
		server.RouteToHandler("PUT", "/oauth/clients/shinyclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/shinyclient/secret"),
			VerifyJSON(`{"clientId":"shinyclient","secret":"shinysecret"}`),
			VerifyHeaderKV("Authorization", "bearer access_token"),
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
			VerifyHeaderKV("Authorization", "bearer access_token"),
			RespondWith(http.StatusUnauthorized, "{}"),
		),
		)

		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Out).To(Say("The secret for client shinyclient was not updated."))
		Expect(session.Out).To(Say("Retry with --trace for more information."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active target", func() {
		config.WriteConfig(uaa.NewConfig())
		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Out).To(Say("You must set a target in order to use this command."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when there is no active context", func() {
		c := uaa.NewConfig()
		t := uaa.NewTarget()
		c.AddTarget(t)
		config.WriteConfig(c)
		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret")

		Expect(session.Out).To(Say("You must set have a token in your context to perform this command."))
		Eventually(session).Should(Exit(1))
	})

	It("supports zone switching", func() {
		server.RouteToHandler("PUT", "/oauth/clients/shinyclient/secret", CombineHandlers(
			VerifyRequest("PUT", "/oauth/clients/shinyclient/secret"),
			VerifyJSON(`{"clientId":"shinyclient","secret":"shinysecret"}`),
			VerifyHeaderKV("Authorization", "bearer access_token"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
			RespondWith(http.StatusOK, "{}"),
		),
		)

		session := runCommand("set-client-secret", "shinyclient", "-s", "shinysecret", "--zone", "twilight-zone")

		Expect(session.Out).To(Say("The secret for client shinyclient has been successfully updated."))
		Eventually(session).Should(Exit(0))
	})

	It("complains when no secret is provided", func() {
		session := runCommand("set-client-secret", "shinyclient")

		Expect(session.Out).To(Say("Missing argument `client_secret` must be specified."))
		Eventually(session).Should(Exit(1))
	})

	It("complains when no clientid is provided", func() {
		session := runCommand("set-client-secret", "-s", "shinysecret")

		Expect(session.Out).To(Say("Missing argument `client_id` must be specified."))
		Eventually(session).Should(Exit(1))
	})
})
