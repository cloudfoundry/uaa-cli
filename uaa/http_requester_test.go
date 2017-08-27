package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("HttpGetter", func() {
	var (
		server *ghttp.Server
		client *http.Client
		config Config
		responseJson string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = &http.Client{}
		config = NewConfigWithServerURL(server.URL())
		responseJson = `{"foo": "bar"}`
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("UnauthenticatedRequester", func() {
		Describe("GetBytes", func() {
			It("calls an endpoint with Accept application/json header", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				UnauthenticatedRequester{}.GetBytes(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns helpful error when GET request fails", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				_, err := UnauthenticatedRequester{}.GetBytes(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})
		})

		Describe("PostBytes", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{
				  "access_token" : "bc4885d950854fed9a938e96b13ca519",
				  "token_type" : "bearer",
				  "expires_in" : 43199,
				  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`

				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyBody([]byte("hello=world")),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				))

				body := map[string]string{"hello": "world",}
				returnedBytes, _ := UnauthenticatedRequester{}.PostBytes(client, config, "/oauth/token", "", body)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("expires_in"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
				))

				_, err := UnauthenticatedRequester{}.PostBytes(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})
		})
	})

	Describe("AuthenticatedRequester", func() {
		Describe("GetBytes", func() {
			It("calls an endpoint with Accept and Authorization headers", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				config.AddContext(UaaContext{AccessToken: "access_token"})
				AuthenticatedRequester{}.GetBytes(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("returns a helpful error when GET request fails", func() {
				server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
					ghttp.RespondWith(500, ""),
					ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
				))

				config.AddContext(UaaContext{AccessToken: "access_token"})
				_, err := AuthenticatedRequester{}.GetBytes(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(UaaContext{AccessToken: ""})
				_, err := AuthenticatedRequester{}.GetBytes(client, config, "/testPath", "someQueryParam=true")

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})

		Describe("PostBytes", func() {
			It("calls an endpoint with correct body and headers", func() {
				responseJson = `{
				  "access_token" : "bc4885d950854fed9a938e96b13ca519",
				  "token_type" : "bearer",
				  "expires_in" : 43199,
				  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`

				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(200, responseJson),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
					ghttp.VerifyBody([]byte("hello=world")),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				))

				body := map[string]string{"hello": "world", }
				config.AddContext(UaaContext{AccessToken: "access_token"})

				returnedBytes, _ := AuthenticatedRequester{}.PostBytes(client, config, "/oauth/token", "", body)
				parsedResponse := string(returnedBytes)

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(parsedResponse).To(ContainSubstring("expires_in"))
			})

			It("returns an error when request fails", func() {
				server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
					ghttp.RespondWith(500, "garbage"),
					ghttp.VerifyRequest("POST", "/oauth/token", ""),
				))

				config.AddContext(UaaContext{AccessToken: "access_token"})
				_, err := AuthenticatedRequester{}.PostBytes(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(1))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			})

			It("returns a helpful error when no token in context", func() {
				config.AddContext(UaaContext{AccessToken: ""})
				_, err := AuthenticatedRequester{}.PostBytes(client, config, "/oauth/token", "", map[string]string{})

				Expect(server.ReceivedRequests()).To(HaveLen(0))
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
			})
		})
	})
})
