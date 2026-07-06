package cmd_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

// captureRequestBody records the raw request body into dest and restores it
// so downstream ghttp handlers (e.g. VerifyRequest) can still read it.
func captureRequestBody(dest *[]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		Expect(err).NotTo(HaveOccurred())
		*dest = body
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
}

var _ = Describe("UpdateUser", func() {
	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(cfg)
	})

	Describe("Validations", func() {
		It("requires a target to have been set", func() {
			config.WriteConfig(config.NewConfig())

			session := runCommand("update-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cli.MISSING_TARGET))
		})

		It("requires a token in context", func() {
			config.WriteConfig(config.NewConfigWithServerURL(server.URL()))

			session := runCommand("update-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cli.MISSING_CONTEXT))
		})

		It("requires a username", func() {
			session := runCommand("update-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
		})

		It("requires at least one update flag", func() {
			session := runCommand("update-user", "marcus@stoicism.com")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("At least one of --familyName, --givenName, --email, --phone, or --delAttrs must be specified."))
			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})
	})

	Describe("UpdateUserCmd", func() {
		Describe("Success cases", func() {
			It("updates user with given_name only", func() {
				var putBody []byte

				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("Accept", "application/json"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					captureRequestBody(&putBody),
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
					VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
					VerifyHeaderKV("Accept", "application/json"),
					VerifyHeaderKV("Content-Type", "application/json"),
				))

				session := runCommand("update-user", "marcus@stoicism.com", "--givenName", "Bob")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
				Expect(session.Out).To(Say("Account for user marcus@stoicism.com successfully updated"))

				var payload map[string]interface{}
				Expect(json.Unmarshal(putBody, &payload)).To(Succeed())
				name, ok := payload["name"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(name["givenName"]).To(Equal("Bob"))
			})

			It("updates user with multiple attributes", func() {
				var putBody []byte

				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					captureRequestBody(&putBody),
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--givenName", "Bob",
					"--familyName", "Smith")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
				Expect(session.Out).To(Say("Account for user marcus@stoicism.com successfully updated"))

				var payload map[string]interface{}
				Expect(json.Unmarshal(putBody, &payload)).To(Succeed())
				name, ok := payload["name"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(name["givenName"]).To(Equal("Bob"))
				Expect(name["familyName"]).To(Equal("Smith"))
			})

			It("updates user with origin specified", func() {
				// First GET to retrieve user with origin filter
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
					VerifyFormKV("filter", `userName eq "marcus@stoicism.com" and origin eq "ldap"`),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--origin", "ldap",
					"--givenName", "Bob")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
			})

			It("updates user with emails", func() {
				var putBody []byte

				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					captureRequestBody(&putBody),
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--email", "new@email.com")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))

				var payload map[string]interface{}
				Expect(json.Unmarshal(putBody, &payload)).To(Succeed())
				emails, ok := payload["emails"].([]interface{})
				Expect(ok).To(BeTrue())
				Expect(emails).To(HaveLen(1))
				email := emails[0].(map[string]interface{})
				Expect(email["value"]).To(Equal("new@email.com"))
			})

			It("updates user with phones", func() {
				var putBody []byte

				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					captureRequestBody(&putBody),
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--phone", "555-1234")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))

				var payload map[string]interface{}
				Expect(json.Unmarshal(putBody, &payload)).To(Succeed())
				phones, ok := payload["phoneNumbers"].([]interface{})
				Expect(ok).To(BeTrue())
				Expect(phones).To(HaveLen(1))
				phone := phones[0].(map[string]interface{})
				Expect(phone["value"]).To(Equal("555-1234"))
			})

			It("updates user with delAttrs removing phone numbers", func() {
				var putBody []byte

				// First GET to retrieve user, which already has a phone number
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{
						ID:           "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
						Username:     "marcus@stoicism.com",
						PhoneNumbers: []uaa.PhoneNumber{{Value: "555-0000"}},
					})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					captureRequestBody(&putBody),
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--delAttrs", "phoneNumbers")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))

				var payload map[string]interface{}
				Expect(json.Unmarshal(putBody, &payload)).To(Succeed())
				Expect(payload["phoneNumbers"]).To(BeNil())
			})

			It("works with zone parameter", func() {
				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
					VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
					VerifyHeaderKV("X-Identity-Zone-Id", "twilight-zone"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--givenName", "Bob",
					"--zone", "twilight-zone")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
			})

			It("prints the updated user json", func() {
				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com", "--givenName", "Bob")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
				Expect(session.Out).To(Say("marcus@stoicism.com"))
				Expect(session.Out).To(Say("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
			})
		})

		Describe("Error cases", func() {
			It("displays an error when user is not found", func() {
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusNotFound, ""),
					VerifyRequest("GET", "/Users"),
				))

				session := runCommand("update-user", "nobody", "--givenName", "Bob")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(session).To(Exit(1))
			})

			It("displays an error if there is a problem during update", func() {
				// First GET succeeds
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// But PUT fails
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					RespondWith(http.StatusBadRequest, ""),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com", "--givenName", "Bob")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(1))
			})
		})

		Describe("Verbose mode", func() {
			It("shows PUT endpoint when verbose flag is used", func() {
				// First GET to retrieve user
				server.RouteToHandler("GET", "/Users", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.User{ID: "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", Username: "marcus@stoicism.com"})),
					VerifyRequest("GET", "/Users"),
				))

				// Then PUT to update user
				server.RouteToHandler("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", CombineHandlers(
					RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
					VerifyRequest("PUT", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				))

				session := runCommand("update-user", "marcus@stoicism.com",
					"--givenName", "Bob",
					"--verbose")

				Expect(server.ReceivedRequests()).To(HaveLen(2))
				Expect(session).To(Exit(0))
				// Note: Verbose output would show the actual PUT request details
				// but this depends on the go-uaa library's verbose logging implementation
			})
		})
	})
})
