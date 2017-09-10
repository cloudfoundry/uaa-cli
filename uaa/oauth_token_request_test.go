package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

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

	const opaqueTokenResponse = `{
	  "access_token" : "bc4885d950854fed9a938e96b13ca519",
	  "refresh_token" : "0cb0e2670f7642e9b501a79252f90f02",
	  "token_type" : "bearer",
	  "expires_in" : 43199,
	  "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
	  "jti" : "bc4885d950854fed9a938e96b13ca519"
	}`
	const jwtTokenResponse = `{
	  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	  "refresh_token" : "0cb0e2670f7642e9b501a79252f90f02",
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

	Describe("ClientCredentialsClient#RequestToken", func() {
		It("makes a POST to /oauth/token", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, opaqueTokenResponse),
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

	Describe("ResourceOwnerPasswordClient#RequestToken", func() {
		It("makes a POST to /oauth/token", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, opaqueTokenResponse),
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
				ClientId:     "identity",
				ClientSecret: "identitysecret",
				Username:     "woodstock",
				Password:     "birdsrule",
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
				ClientId:     "identity",
				ClientSecret: "identitysecret",
				Username:     "woodstock",
				Password:     "birdsrule",
			}
			_, err := ropClient.RequestToken(client, config, OPAQUE)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("ClientCredentialsClient#RequestToken", func() {
		It("makes a POST to /oauth/token", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, opaqueTokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyFormKV("client_id", "my_authcode_client"),
				ghttp.VerifyFormKV("client_secret", "clientsecret"),
				ghttp.VerifyFormKV("grant_type", "authorization_code"),
				ghttp.VerifyFormKV("token_format", string(OPAQUE)),
				ghttp.VerifyFormKV("response_type", "token"),
				ghttp.VerifyFormKV("code", "abcde"),
				ghttp.VerifyFormKV("redirect_uri", "http://localhost:8080"),
			))

			authcodeClient := AuthorizationCodeClient{
				ClientId:     "my_authcode_client",
				ClientSecret: "clientsecret",
			}
			tokenResponse, _ := authcodeClient.RequestToken(client, config, OPAQUE, "abcde", "http://localhost:8080")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(tokenResponse.AccessToken).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			Expect(tokenResponse.TokenType).To(Equal("bearer"))
			Expect(tokenResponse.ExpiresIn).To(Equal(int32(43199)))
			Expect(tokenResponse.Scope).To(Equal("clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write"))
			Expect(tokenResponse.JTI).To(Equal("bc4885d950854fed9a938e96b13ca519"))
		})

		It("can request opaque or jwt format tokens", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, jwtTokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyFormKV("client_id", "my_authcode_client"),
				ghttp.VerifyFormKV("client_secret", "clientsecret"),
				ghttp.VerifyFormKV("grant_type", "authorization_code"),
				ghttp.VerifyFormKV("token_format", string(JWT)),
				ghttp.VerifyFormKV("response_type", "token"),
				ghttp.VerifyFormKV("code", "abcde"),
				ghttp.VerifyFormKV("redirect_uri", "http://localhost:8080"),
			))

			authcodeClient := AuthorizationCodeClient{
				ClientId:     "my_authcode_client",
				ClientSecret: "clientsecret",
			}
			tokenResponse, _ := authcodeClient.RequestToken(client, config, JWT, "abcde", "http://localhost:8080")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(tokenResponse.AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
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

			authcodeClient := AuthorizationCodeClient{
				ClientId:     "my_authcode_client",
				ClientSecret: "clientsecret",
			}

			_, err := authcodeClient.RequestToken(client, config, JWT, "abcde", "http://localhost:8080")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("RefreshTokenClient#RequestToken", func() {
		It("requests a new token by passing a refresh token", func() {
			server.RouteToHandler("POST", "/oauth/token", ghttp.CombineHandlers(
				ghttp.RespondWith(200, opaqueTokenResponse),
				ghttp.VerifyRequest("POST", "/oauth/token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/x-www-form-urlencoded"),
				ghttp.VerifyFormKV("client_id", "someclient"),
				ghttp.VerifyFormKV("client_secret", "somesecret"),
				ghttp.VerifyFormKV("token_format", string(OPAQUE)),
				ghttp.VerifyFormKV("grant_type", "refresh_token"),
				ghttp.VerifyFormKV("response_type", "token"),
				ghttp.VerifyFormKV("refresh_token", "the_refresh_token"),
			))

			refreshClient := RefreshTokenClient{
				ClientId:     "someclient",
				ClientSecret: "somesecret",
			}
			tokenResponse, _ := refreshClient.RequestToken(client, config, OPAQUE, "the_refresh_token")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(tokenResponse.AccessToken).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			Expect(tokenResponse.RefreshToken).To(Equal("0cb0e2670f7642e9b501a79252f90f02"))
			Expect(tokenResponse.TokenType).To(Equal("bearer"))
			Expect(tokenResponse.ExpiresIn).To(Equal(int32(43199)))
			Expect(tokenResponse.Scope).To(Equal("clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write"))
			Expect(tokenResponse.JTI).To(Equal("bc4885d950854fed9a938e96b13ca519"))
		})
	})
})
