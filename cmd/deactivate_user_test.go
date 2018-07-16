package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("DeactivateUser", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("deactivates a user", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22&startIndex=1&count=100"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PATCH", "/Users/abcdef", CombineHandlers(
			VerifyRequest("PATCH", "/Users/abcdef", ""),
			VerifyHeaderKV("If-Match", "10"),
			VerifyJSON(`{"active": false}`),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("deactivate-user", "woodstock@peanuts.com")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Account for user woodstock@peanuts.com successfully deactivated."))
	})

	It("deactivates a user in a non-default zone", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22&startIndex=1&count=100"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilightzone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PATCH", "/Users/abcdef", CombineHandlers(
			VerifyRequest("PATCH", "/Users/abcdef", ""),
			VerifyHeaderKV("If-Match", "10"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilightzone"),
			VerifyJSON(`{"active": false}`),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("deactivate-user", "woodstock@peanuts.com", "--zone", "twilightzone")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Account for user woodstock@peanuts.com successfully deactivated."))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			config.WriteConfig(config.NewConfig())

			session := runCommand("deactivate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must set a target in order to use this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(cfg)

			session := runCommand("deactivate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a username", func() {
			c := config.NewConfigWithServerURL(server.URL())
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			c.AddContext(ctx)
			config.WriteConfig(c)

			session := runCommand("deactivate-user")

			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
			Expect(session).Should(Exit(1))
		})
	})
})
