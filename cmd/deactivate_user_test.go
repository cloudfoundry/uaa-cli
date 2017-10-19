package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

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

	Describe("validations", func() {
		It("requires a target", func() {
			err := GetUserValidations(uaa.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})

		It("requires a username", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")
			ctx := uaa.NewContextWithToken("access_token")
			cfg.AddContext(ctx)

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional argument USERNAME must be specified."))

			err = GetUserValidations(cfg, []string{"userid"})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
