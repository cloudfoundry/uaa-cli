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

var _ = Describe("AddMember", func() {
	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("creates a membership in a group", func() {
			membershipJson := `{"origin":"uaa","type":"USER","value":"fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"}`

			server.RouteToHandler("POST", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members", CombineHandlers(
				VerifyRequest("POST", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members"),
				VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
				VerifyHeaderKV("Accept", "application/json"),
				VerifyJSON(membershipJson),
				RespondWith(http.StatusOK, membershipJson, contentTypeJson),
			))
			server.RouteToHandler("GET", "/Groups", CombineHandlers(
				VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22uaa.admin%22&count=100&startIndex=1"),
				RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{ID: "05a0c169-3592-4a45-b109-a16d9246e0ab", DisplayName: "uaa.admin"})),
			))
			server.RouteToHandler("GET", "/Users", CombineHandlers(
				VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22&count=100&startIndex=1"),
				RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "woodstock@peanuts.com"})),
			))

			session := runCommand("add-member", "uaa.admin", "woodstock@peanuts.com")

			Eventually(session).Should(Exit(0))
			Expect(session).To(Say("User woodstock@peanuts.com successfully added to group uaa.admin"))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("add-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(MISSING_TARGET))
		})
	})

	Describe("when no token in context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("tells the user to get a token", func() {
			session := runCommand("add-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(MISSING_CONTEXT))
		})
	})

	Describe("validations", func() {
		It("only accepts groupname and username", func() {
			session := runCommand("add-member", "first-arg", "second-arg", "third-arg")
			Eventually(session).Should(Exit(1))

			session = runCommand("add-member", "woodstock")
			Eventually(session).Should(Exit(1))
		})
	})
})
