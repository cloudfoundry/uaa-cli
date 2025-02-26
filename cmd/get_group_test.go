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

var _ = Describe("GetGroup", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("looks up a group with a SCIM filter", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22&startIndex=1&count=100"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	It("can limit results data with --attributes", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22&attributes=displayName&startIndex=1&count=100"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin", "--attributes", "displayName")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	It("can understand the --zone flag", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22&startIndex=1&count=100"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin", "--zone", "twilight-zone")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			err := GetGroupValidations(config.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := GetGroupValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})

		It("requires a groupname", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			cfg.AddContext(ctx)

			err := GetGroupValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional argument GROUPNAME must be specified."))

			err = GetGroupValidations(cfg, []string{"groupid"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
