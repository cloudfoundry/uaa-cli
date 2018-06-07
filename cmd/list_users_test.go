package cmd_test

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
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
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending"),
			RespondWith(http.StatusOK, userListResponse),
		))

		session := runCommand(
			"list-users",
			"--filter", "verified eq false",
			"--attributes", "id,userName",
			"--sortBy", "userName",
			"--sortOrder", "descending",
		)

		Eventually(session).Should(Exit(0))
	})

	It("understands the --zone flag", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "foozone"),
			RespondWith(http.StatusOK, userListResponse),
		))

		session := runCommand(
			"list-users",
			"--filter", "verified eq false",
			"--attributes", "id,userName",
			"--sortBy", "userName",
			"--sortOrder", "descending",
			"--zone", "foozone",
		)

		Eventually(session).Should(Exit(0))
	})
})
