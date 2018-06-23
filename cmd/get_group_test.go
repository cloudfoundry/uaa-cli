package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

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

var _ = Describe("GetGroup", func() {
	BeforeEach(func() {
		c := uaa.NewConfigWithServerURL(server.URL())
		ctx := uaa.NewContextWithToken("access_token")
		c.AddContext(ctx)
		config.WriteConfig(c)
	})

	It("looks up a group with a SCIM filter", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	It("can limit results data with --attributes", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22&attributes=displayName"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin", "--attributes", "displayName")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	It("can understand the --zone flag", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{DisplayName: "admin"})),
		))

		session := runCommand("get-group", "admin", "--zone", "twilight-zone")

		Eventually(session).Should(Say(`"displayName": "admin"`))
		Eventually(session).Should(Exit(0))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			err := GetGroupValidations(uaa.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := GetGroupValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})

		It("requires a groupname", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")
			ctx := uaa.NewContextWithToken("access_token")
			cfg.AddContext(ctx)

			err := GetGroupValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional argument GROUPNAME must be specified."))

			err = GetGroupValidations(cfg, []string{"groupid"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
