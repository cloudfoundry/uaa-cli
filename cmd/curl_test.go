package cmd_test

import (
	"fmt"
	"github.com/onsi/gomega/gbytes"
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Curl", func() {

	var userListResponse string

	BeforeEach(func() {
		cfg := config.NewConfigWithServerURL(server.URL())
		cfg.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		config.WriteConfig(cfg)
		userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)
	})

	It("appends the access token from saved context", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", ""),
			VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
			RespondWith(http.StatusOK, userListResponse, contentTypeJson),
		))

		session := runCommand("curl", "/Users")

		Eventually(session).Should(Exit(0))
	})

	It("sends GET request by default", func() {
		server.RouteToHandler("GET", "/Users", CombineHandlers(
			VerifyRequest("GET", "/Users", ""),
			RespondWith(http.StatusOK, userListResponse, contentTypeJson),
		))

		session := runCommand("curl", "/Users")

		Eventually(session).Should(Exit(0))
	})

	It("can send POST request", func() {
		server.RouteToHandler("POST", "/Users", CombineHandlers(
			VerifyRequest("POST", "/Users", ""),
			RespondWith(http.StatusCreated, userListResponse),
		))

		session := runCommand("curl",
			"/Users",
			"-X", "POST",
			"-H", "Accept: application/json")

		Eventually(session).Should(Exit(0))
	})

	When("when the server returns a 4xx or 5xx", func() {
		Context("status code 400", func() {
			itShouldFailWhenServerReturns(http.StatusBadRequest)
		})

		Context("status code 404", func() {
			itShouldFailWhenServerReturns(http.StatusNotFound)
		})

		Context("status code 500", func() {
			itShouldFailWhenServerReturns(http.StatusInternalServerError)
		})

		Context("status code 503", func() {
			itShouldFailWhenServerReturns(http.StatusServiceUnavailable)
		})
	})

	It("can send DELETE request", func() {
		server.RouteToHandler("DELETE", "/Users/userguid", CombineHandlers(
			VerifyRequest("DELETE", "/Users/userguid", ""),
			RespondWith(http.StatusOK, MarcusUserResponse, contentTypeJson),
		))

		session := runCommand("curl",
			"/Users/userguid",
			"-X", "DELETE",
			"-H", "Accept: application/json")

		Eventually(session).Should(Exit(0))
	})

	It("can send PUT request with body", func() {
		server.RouteToHandler("PUT", "/Users/userguid", CombineHandlers(
			VerifyRequest("PUT", "/Users/userguid", ""),
			VerifyBody([]byte(`{ "active" : false }`)),
			VerifyHeaderKV("Content-Type", "application/json"),
			RespondWith(http.StatusOK, MarcusUserResponse, contentTypeJson),
		))

		session := runCommand("curl",
			"/Users/userguid",
			"-X", "PUT",
			"-d", `{ "active" : false }`,
			"-H", "Content-Type: application/json")

		Eventually(session).Should(Exit(0))
	})

	It("can send PATCH request with body", func() {
		server.RouteToHandler("PATCH", "/Users/userguid", CombineHandlers(
			VerifyRequest("PATCH", "/Users/userguid", ""),
			VerifyBody([]byte(`{ "active" : false }`)),
			VerifyHeaderKV("Content-Type", "application/json"),
			RespondWith(http.StatusOK, MarcusUserResponse, contentTypeJson),
		))

		session := runCommand("curl",
			"/Users/userguid",
			"-X", "PATCH",
			"-d", `{ "active" : false }`,
			"-H", "Content-Type: application/json")

		Eventually(session).Should(Exit(0))
	})

	It("handles parses multiple header flags correctly", func() {
		server.RouteToHandler("POST", "/Users", CombineHandlers(
			VerifyRequest("POST", "/Users", ""),
			VerifyHeaderKV("Accept", "application/json"),
			VerifyHeaderKV("Content-Type", "application/json"),
			VerifyHeaderKV("Pragma", "no-cache"),
			RespondWith(http.StatusCreated, userListResponse),
		))

		session := runCommand("curl",
			"/Users",
			"-X", "POST",
			"-H", "Accept: application/json",
			"-H", "Content-Type: application/json",
			"-H", "Pragma: no-cache",
		)

		Eventually(session).Should(Exit(0))
	})
})

func itShouldFailWhenServerReturns(statusCode int) {
	BeforeEach(func() {
		server.RouteToHandler("POST", "/Users", CombineHandlers(
			VerifyRequest("POST", "/Users", ""),
			RespondWith(statusCode, "test-response"),
		))
	})

	It("returns a nonzero status code", func() {
		session := runCommand("curl", "/Users", "-X", "POST", "-H", "Accept: application/json")
		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(gbytes.Say("test-response"))
	})
}
