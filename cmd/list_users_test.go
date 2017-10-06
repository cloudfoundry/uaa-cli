package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("ListUsers", func() {

	var userListResponse string

	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
		userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)
	})

	It("executes SCIM queries based on flags", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending&count=50&startIndex=100"),
			RespondWith(http.StatusOK, userListResponse),
		))

		session := runCommand("list-users",
			"--filter", "verified eq false",
			"--attributes", "id,userName",
			"--sortBy", "userName",
			"--sortOrder", "descending",
			"--count", "50",
			"--startIndex", "100")

		Eventually(session).Should(Exit(0))
	})

	It("understands the --zone flag", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending&count=50&startIndex=100"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "foozone"),
			RespondWith(http.StatusOK, userListResponse),
		))

		session := runCommand("list-users",
			"--filter", "verified eq false",
			"--attributes", "id,userName",
			"--sortBy", "userName",
			"--sortOrder", "descending",
			"--count", "50",
			"--startIndex", "100",
			"--zone", "foozone")

		Eventually(session).Should(Exit(0))
	})
})
