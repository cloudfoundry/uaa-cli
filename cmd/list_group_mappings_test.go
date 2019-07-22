package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("ListGroupMappings", func() {
	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		err := config.WriteConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	It("requests group mappings from the backend with default parameters", func() {
		server.RouteToHandler("GET", "/Groups/External", CombineHandlers(
			VerifyRequest("GET", "/Groups/External", "startIndex=1&count=100"),
			RespondWith(http.StatusOK, ExternalGroupsApiResponse, contentTypeJson),
		))

		session := runCommand("list-group-mappings")

		Eventually(session).Should(Exit(0))

		// We can't verify that the right JSON was output
		// There seems to be a gap in the tooling.
		// We can test a regex against a buffer
		// We can test JSON against a string
		// But we can't test JSON against a buffer
		Eventually(session.Out).Should(gbytes.Say("organizations.acme"))
	})

	It("prints a useful description in the help menu", func() {
		session := runCommand("list-group-mappings", "-h")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(gbytes.Say("List all the mappings between uaa scopes and external groups"))
	})
})
