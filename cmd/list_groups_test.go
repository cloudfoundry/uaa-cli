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

var _ = Describe("ListGroups", func() {

	var groupListResponse string

	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(cfg)
		groupListResponse = fmt.Sprintf(PaginatedResponseTmpl, UaaAdminGroupResponse, CloudControllerReadGroupResponse)
	})

	It("executes SCIM queries based on flags", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=verified+eq+false&attributes=id%2CdisplayName&sortBy=displayName&sortOrder=descending&startIndex=1&count=100"),
			RespondWith(http.StatusOK, groupListResponse, contentTypeJson),
		))

		session := runCommand("list-groups",
			"--filter", "verified eq false",
			"--attributes", "id,displayName",
			"--sortBy", "displayName",
			"--sortOrder", "descending",
		)

		Eventually(session).Should(Exit(0))
	})

	It("understands the --zone flag", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=verified+eq+false&attributes=id%2CdisplayName&sortBy=displayName&sortOrder=descending&startIndex=1&count=100"),
			VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
			RespondWith(http.StatusOK, groupListResponse, contentTypeJson),
		))

		session := runCommand("list-groups",
			"--filter", "verified eq false",
			"--attributes", "id,displayName",
			"--sortBy", "displayName",
			"--sortOrder", "descending",
			"--zone", "twilight-zone",
		)

		Eventually(session).Should(Exit(0))
	})
})
