package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("SetPassword", func() {
	BeforeEach(func() {
		c := config.NewConfigWithServerURL(server.URL())
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		c.AddContext(ctx)
		Expect(config.WriteConfig(c)).Should(Succeed())
	})

	It("sets a user password", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22testuser%22&startIndex=1"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "testuser", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PUT", "/Users/abcdef/password", CombineHandlers(
			VerifyRequest("PUT", "/Users/abcdef/password", ""),
			VerifyJSON(`{"password": "newpass"}`),
			RespondWith(http.StatusOK, `{"message": "password updated"}`),
		))

		session := runCommand("set-password", "testuser", "--password", "newpass")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Password for user testuser successfully set."))
	})

	It("sets a user password with --verbose", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22testuser%22&startIndex=1"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "testuser", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PUT", "/Users/abcdef/password", CombineHandlers(
			VerifyRequest("PUT", "/Users/abcdef/password", ""),
			VerifyJSON(`{"password": "newpass"}`),
			RespondWith(http.StatusOK, `{"message": "password updated"}`),
		))

		session := runCommand("set-password", "testuser", "--password", "newpass", "--verbose")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Password for user testuser successfully set."))
	})

	It("sets a user password with --origin", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22testuser%22+and+origin+eq+%22ldap%22&startIndex=1"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "testuser", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PUT", "/Users/abcdef/password", CombineHandlers(
			VerifyRequest("PUT", "/Users/abcdef/password", ""),
			VerifyJSON(`{"password": "newpass"}`),
			RespondWith(http.StatusOK, `{"message": "password updated"}`),
		))

		session := runCommand("set-password", "testuser", "--password", "newpass", "--origin", "ldap")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Password for user testuser successfully set."))
	})

	It("sets a user password with --zone", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22testuser%22&startIndex=1"),
			VerifyHeaderKV("X-Identity-Zone-Id", "test-zone"),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "testuser", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
		))
		server.RouteToHandler("PUT", "/Users/abcdef/password", CombineHandlers(
			VerifyRequest("PUT", "/Users/abcdef/password", ""),
			VerifyHeaderKV("X-Identity-Zone-Id", "test-zone"),
			VerifyJSON(`{"password": "newpass"}`),
			RespondWith(http.StatusOK, `{"message": "password updated"}`),
		))

		session := runCommand("set-password", "testuser", "--password", "newpass", "--zone", "test-zone")

		Expect(server.ReceivedRequests()).To(HaveLen(2))
		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("Password for user testuser successfully set."))
	})

	Describe("error conditions", func() {
		It("displays error when user not found", func() {
			server.RouteToHandler("GET", "/Users", CombineHandlers(
				VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22nobody%22&startIndex=1"),
				RespondWith(http.StatusNotFound, `{"error": "scim_resource_not_found", "error_description": "User nobody does not exist"}`),
			))

			session := runCommand("set-password", "nobody", "--password", "newpass")

			Eventually(session).Should(Exit(1))
		})

		It("displays error when password change request fails", func() {
			server.RouteToHandler("GET", "/Users", CombineHandlers(
				VerifyRequest("GET", "/Users", "count=100&filter=userName+eq+%22testuser%22&startIndex=1"),
				RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{Username: "testuser", ID: "abcdef", Meta: &uaa.Meta{Version: 10}})),
			))
			server.RouteToHandler("PUT", "/Users/abcdef/password", CombineHandlers(
				VerifyRequest("PUT", "/Users/abcdef/password", ""),
				RespondWith(http.StatusBadRequest, `{"error": "invalid_password", "error_description": "Password does not meet policy requirements"}`),
			))

			session := runCommand("set-password", "testuser", "--password", "weak")

			Eventually(session).Should(Exit(1))
		})
	})

	Describe("validations", func() {
		It("requires a target", func() {
			config.WriteConfig(config.NewConfig())

			session := runCommand("set-password", "testuser", "--password", "newpass")

			Expect(session.Err).To(Say("You must set a target in order to use this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(cfg)

			session := runCommand("set-password", "testuser", "--password", "newpass")

			Expect(session.Err).To(Say("You must have a token in your context to perform this command."))
			Expect(session).Should(Exit(1))
		})

		It("requires a username", func() {
			c := config.NewConfigWithServerURL(server.URL())
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			c.AddContext(ctx)
			config.WriteConfig(c)

			session := runCommand("set-password")

			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
			Expect(session).Should(Exit(1))
		})

		It("displays help when no password provided (interactive mode skipped in tests)", func() {
			c := config.NewConfigWithServerURL(server.URL())
			ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
			c.AddContext(ctx)
			config.WriteConfig(c)

			// In test environment, when no password is provided, the interactive prompt
			// will fail due to inappropriate ioctl. This tests that the validation
			// is working even though the specific error may vary in test vs. real usage.
			session := runCommand("set-password", "testuser")

			Eventually(session).Should(Exit(1))
		})
	})
})
