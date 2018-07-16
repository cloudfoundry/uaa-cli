package cmd_test

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("ListUsers", func() {

	var userListResponse string

	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(cfg)
		userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)
	})

	It("executes SCIM queries based on flags", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending&count=100&startIndex=1"),
			RespondWith(http.StatusOK, userListResponse, contentTypeJson),
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
			VerifyRequest("GET", "/Users", "filter=verified+eq+false&attributes=id%2CuserName&sortBy=userName&sortOrder=descending&count=100&startIndex=1"),
			VerifyHeaderKV("X-Identity-Zone-Id", "foozone"),
			RespondWith(http.StatusOK, userListResponse, contentTypeJson),
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
