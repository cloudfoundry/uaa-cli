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

var _ = Describe("ListGroups", func() {

	var groupListResponse string

	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
		groupListResponse = fmt.Sprintf(PaginatedResponseTmpl, UaaAdminGroupResponse, CloudControllerReadGroupResponse)
	})

	It("executes SCIM queries based on flags", func() {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", "filter=verified+eq+false&attributes=id%2CdisplayName&sortBy=displayName&sortOrder=descending"),
			RespondWith(http.StatusOK, groupListResponse),
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
			VerifyRequest("GET", "/Groups", "filter=verified+eq+false&attributes=id%2CdisplayName&sortBy=displayName&sortOrder=descending"),
			VerifyHeaderKV("X-Identity-Zone-Subdomain", "twilight-zone"),
			RespondWith(http.StatusOK, groupListResponse),
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
