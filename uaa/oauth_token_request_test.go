package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/url"
)

var _ = Describe("OauthTokenRequest", func() {
	var (
		server *ghttp.Server
		config Config
		client *http.Client
	)

	const tokenResponse = `{
	  "access_token" : "bc4885d950854fed9a938e96b13ca519",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = &http.Client{}
		config = NewConfigWithServerURL(server.URL())
	})

	Describe("ClientCredentialsClient", func() {
		It("makes a POST to /oauth/token", func() {
			data := url.Values{}
			data.Add("client_id", "identity")
			data.Add("client_secret", "identitysecret")
			data.Add("grant_type", "client_credentials")
			data.Add("token_format", string(OPAQUE))
			data.Add("response_type", "token")

			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, tokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyForm(data),
			))

			ccClient := ClientCredentialsClient{ClientId: "identity", ClientSecret: "identitysecret"}
			tokenResponse, _ := ccClient.RequestToken(client, config, OPAQUE)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(tokenResponse.AccessToken).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			Expect(tokenResponse.TokenType).To(Equal("bearer"))
			Expect(tokenResponse.ExpiresIn).To(Equal(int32(43199)))
			Expect(tokenResponse.Scope).To(Equal("clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write"))
			Expect(tokenResponse.JTI).To(Equal("bc4885d950854fed9a938e96b13ca519"))
		})

		It("returns error if response unparasable", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, "{garbage response}"),
				ghttp.VerifyRequest("POST", "/oauth/token"),
			))

			ccClient := ClientCredentialsClient{ClientId: "identity", ClientSecret: "identitysecret"}
			_, err := ccClient.RequestToken(client, config, OPAQUE)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
		})
	})
})
