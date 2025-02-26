package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("ActivateUser", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		c.AddContext(ctx)
		Expect(config.WriteConfig(c)).Should(Succeed())
	})

	It("activates a user", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22woodstock%40peanuts.com%22&startIndex=1"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PATCH", "/Users/abcdef", CombineHandlers(
			VerifyRequest("PATCH", "/Users/abcdef", ""),
			VerifyHeaderKV("If-Match", "10"),
			VerifyJSON(`{"active": true}`),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("activate-user", "woodstock@peanuts.com")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Account for user woodstock@peanuts.com successfully activated."))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			config.WriteConfig(config.NewConfig())

			session := runCommand("activate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must set a target in order to use this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(cfg)

			session := runCommand("activate-user", "woodstock@peanuts.com")

			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a username", func() {
			c := config.NewConfigWithServerURL(server.URL())
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			c.AddContext(ctx)
			config.WriteConfig(c)

			session := runCommand("activate-user")

			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
			Expect(session).Should(Exit(1))
		})
	})
})
