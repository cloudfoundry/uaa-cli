package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
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
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, tokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyFormKV("client_id", "identity"),
				ghttp.VerifyFormKV("client_secret", "identitysecret"),
				ghttp.VerifyFormKV("grant_type", "client_credentials"),
				ghttp.VerifyFormKV("token_format", string(OPAQUE)),
				ghttp.VerifyFormKV("response_type", "token"),
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

	Describe("ResourceOwnerPasswordClient", func() {
		It("makes a POST to /oauth/token", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, tokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyFormKV("client_id", "identity"),
				ghttp.VerifyFormKV("client_secret", "identitysecret"),
				ghttp.VerifyFormKV("grant_type", "password"),
				ghttp.VerifyFormKV("token_format", string(OPAQUE)),
				ghttp.VerifyFormKV("response_type", "token"),
				ghttp.VerifyFormKV("username", "woodstock"),
				ghttp.VerifyFormKV("password", "birdsrule"),
			))

			ropClient := ResourceOwnerPasswordClient{
				ClientId: "identity",
				ClientSecret: "identitysecret",
				Username: "woodstock",
				Password: "birdsrule",
			}
			tokenResponse, _ := ropClient.RequestToken(client, config, OPAQUE)

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

			ropClient := ResourceOwnerPasswordClient{
				ClientId: "identity",
				ClientSecret: "identitysecret",
				Username: "woodstock",
				Password: "birdsrule",
			}
			_, err := ropClient.RequestToken(client, config, OPAQUE)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
		})
	})
})
