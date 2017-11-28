package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("DeactivateUser", func() {
	BeforeEach(func() {
		c := uaa.NewConfigWithServerURL(server.URL())
		ctx := uaa.NewContextWithToken("access_token")
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("deactivates a user", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.ScimUser{Username: "woodstock@peanuts.com", ID: "abcdef", Meta: &uaa.ScimMetaInfo{Version: 10}})),
		))
		server.RouteToHandler("PATCH", "/Users/abcdef", CombineHandlers(
			VerifyRequest("PATCH", "/Users/abcdef", ""),
			VerifyHeaderKV("If-Match", "10"),
			VerifyJSON(`{"active": false}`),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.ScimUser{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("deactivate-user", "woodstock@peanuts.com")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Account for user woodstock@peanuts.com successfully deactivated."))
	})

	It("deactivates a user in a non-default zone", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilightzone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.ScimUser{Username: "woodstock@peanuts.com", ID: "abcdef", Meta: &uaa.ScimMetaInfo{Version: 10}})),
		))
		server.RouteToHandler("PATCH", "/Users/abcdef", CombineHandlers(
			VerifyRequest("PATCH", "/Users/abcdef", ""),
			VerifyHeaderKV("If-Match", "10"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilightzone"),
			VerifyJSON(`{"active": false}`),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.ScimUser{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("deactivate-user", "woodstock@peanuts.com", "--zone", "twilightzone")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Account for user woodstock@peanuts.com successfully deactivated."))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			config.WriteConfig(uaa.NewConfig())

			session := runCommand("deactivate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must set a target in order to use this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a context", func() {
			cfg := uaa.NewConfigWithServerURL(server.URL())
			config.WriteConfig(cfg)

			session := runCommand("deactivate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a username", func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			ctx := uaa.NewContextWithToken("access_token")
			c.AddContext(ctx)
			config.WriteConfig(c)

			session := runCommand("deactivate-user")

			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
			Expect(session).Should(Exit(1))
		})
	})
})
