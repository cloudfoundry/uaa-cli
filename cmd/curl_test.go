package cmd_test

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Curl", func() {

	var userListResponse string

	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
		userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)
	})

	It("sends request", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", ""),
			RespondWith(http.StatusOK, userListResponse),
		))

		session := runCommand("curl",
			"/Users",
			"-X", "GET",
			"-H", "Accept: application/json")

		Eventually(session).Should(Exit(0))
	})
})
