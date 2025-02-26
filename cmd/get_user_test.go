package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("GetUser", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("looks up a user with a SCIM filter", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22+and+origin+eq+%22uaa%22&startIndex=1&count=100"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("get-user", "woodstock@peanuts.com", "--origin", "uaa")

		Eventually(session).Should(Say(`"userName": "woodstock@peanuts.com"`))
		Eventually(session).Should(Exit(0))
	})

	It("can understand the --zone flag", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22&startIndex=1&count=100"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("get-user", "woodstock@peanuts.com", "--zone", "twilight-zone")

		Eventually(session).Should(Say(`"userName": "woodstock@peanuts.com"`))
		Eventually(session).Should(Exit(0))
	})

	It("can limit results data with --attributes", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22+and+origin+eq+%22uaa%22&attributes=userName&startIndex=1&count=100"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "woodstock@peanuts.com"})),
		))

		session := runCommand("get-user", "woodstock@peanuts.com", "--origin", "uaa", "--attributes", "userName")

		Eventually(session).Should(Say(`"userName": "woodstock@peanuts.com"`))
		Eventually(session).Should(Exit(0))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			err := GetUserValidations(config.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})

		It("requires a username", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			cfg.AddContext(ctx)

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional argument USERNAME must be specified."))

			err = GetUserValidations(cfg, []string{"userid"})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
