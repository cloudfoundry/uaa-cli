package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("ListMFAProviders", func() {
	var mfaProvidersListResponse string

	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(cfg)
		mfaProvidersListResponse = fixtures.ListMFAProvidersResponse
	})

	It("lists mfa providers", func() {
		server.RouteToHandler("GET", "/mfa-providers", CombineHandlers(
			VerifyRequest("GET", "/mfa-providers", ""),
			RespondWith(http.StatusOK, mfaProvidersListResponse, contentTypeJson)))

		session := runCommand("list-mfa-providers")

		Eventually(session).Should(Exit(0))
	})
})
