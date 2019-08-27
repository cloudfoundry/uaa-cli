package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("RemoveMember", func() {
	Describe("when a target is set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		Describe("when the user belongs to the group", func() {
			It("removes a membership in a group", func() {
				membershipJson := `{"origin":"uaa","type":"USER","value":"fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"}`

				server.RouteToHandler("DELETE", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					VerifyRequest("DELETE", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("Accept", "application/json"),
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

				session := runCommand("remove-member", "uaa.admin", "woodstock@peanuts.com")

				Eventually(session).Should(Exit(0))
				Expect(session).To(Say("User woodstock@peanuts.com successfully removed from group uaa.admin"))
			})
		})

		Describe("when the user belongs to the group", func() {
			It("removes a membership in a group", func() {
				errorJson := `{"error_description":"Member 88553b0f-af97-489a-a884-d05f54e3b9ce does not exist in group 9c7c01bc-d234-4b43-b93d-35bb3ec5b906","error":"member_not_found","message":"Member 88553b0f-af97-489a-a884-d05f54e3b9ce does not exist in group 9c7c01bc-d234-4b43-b93d-35bb3ec5b906"}`

				server.RouteToHandler("DELETE", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					VerifyRequest("DELETE", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab/members/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("Accept", "application/json"),
					RespondWith(http.StatusNotFound, errorJson, contentTypeJson),
				))

				server.RouteToHandler("GET", "/Groups", CombineHandlers(
					VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22uaa.admin%22&count=100&startIndex=1"),
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{ID: "05a0c169-3592-4a45-b109-a16d9246e0ab", DisplayName: "uaa.admin"})),
				))

				server.RouteToHandler("GET", "/Users", CombineHandlers(
					VerifyRequest("GET", "/Users", "filter=userName+eq+%22woodstock@peanuts.com%22&count=100&startIndex=1"),
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "woodstock@peanuts.com"})),
				))

				session := runCommand("remove-member", "uaa.admin", "woodstock@peanuts.com")

				Eventually(session).Should(Exit(1))
				Expect(session).To(Say("Retry with --verbose for more information."))
			})
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := config.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("remove-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cli.MISSING_TARGET))
		})
	})

	Describe("when no token in context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("tells the user to get a token", func() {
			session := runCommand("remove-member", "uaa.admin", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cli.MISSING_CONTEXT))
		})
	})

	Describe("validations", func() {
		It("only accepts groupname and username", func() {
			session := runCommand("remove-member", "first-arg", "second-arg", "third-arg")
			Eventually(session).Should(Exit(1))

			session = runCommand("remove-member", "woodstock")
			Eventually(session).Should(Exit(1))
		})
	})
})
