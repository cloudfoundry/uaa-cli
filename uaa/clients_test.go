package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Clients", func() {
	Describe("Get", func() {
		var (
			server *ghttp.Server
			config Config
			httpClient *http.Client
		)

		const GetClientResponseJson string = `{
		  "scope" : [ "clients.read", "clients.write" ],
		  "client_id" : "clientid",
		  "resource_ids" : [ "none" ],
		  "authorized_grant_types" : [ "client_credentials" ],
		  "redirect_uri" : [ "http://ant.path.wildcard/**/passback/*", "http://test1.com" ],
		  "autoapprove" : [ "true" ],
		  "authorities" : [ "clients.read", "clients.write" ],
		  "token_salt" : "1SztLL",
		  "allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
		  "name" : "My Client Name",
		  "lastModified" : 1502816030525,
		  "required_user_groups" : [ ]
		}`

		BeforeEach(func() {
			server = ghttp.NewServer()
			httpClient = &http.Client{}
			config = NewConfigWithServerURL(server.URL())
			ctx := UaaContext{AccessToken: "access_token" }
			config.AddContext(ctx)
		})

		AfterEach(func() {
			server.Close()
		})

		It("calls the /oauth/clients endpoint", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid", ghttp.CombineHandlers(
				ghttp.RespondWith(200, GetClientResponseJson),
				ghttp.VerifyRequest("GET", "/oauth/clients/clientid"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			cm := &ClientManager{httpClient, config}
			clientResponse, _ := cm.Get("clientid")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(clientResponse.Scope[0]).To(Equal("clients.read"))
			Expect(clientResponse.Scope[1]).To(Equal("clients.write"))
			Expect(clientResponse.ClientId).To(Equal("clientid"))
			Expect(clientResponse.ResourceIds[0]).To(Equal("none"))
			Expect(clientResponse.AuthorizedGrantTypes[0]).To(Equal("client_credentials"))
			Expect(clientResponse.RedirectUri[0]).To(Equal("http://ant.path.wildcard/**/passback/*"))
			Expect(clientResponse.RedirectUri[1]).To(Equal("http://test1.com"))
			Expect(clientResponse.Autoapprove[0]).To(Equal("true")) // TODO wtf is autoapprove
			Expect(clientResponse.TokenSalt).To(Equal("1SztLL"))
			Expect(clientResponse.AllowedProviders[0]).To(Equal("uaa"))
			Expect(clientResponse.AllowedProviders[1]).To(Equal("ldap"))
			Expect(clientResponse.AllowedProviders[2]).To(Equal("my-saml-provider"))
			Expect(clientResponse.Name).To(Equal("My Client Name"))
			Expect(clientResponse.LastModified).To(Equal(int64(1502816030525)))
		})

		It("returns helpful error when /oauth/clients request fails", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid", ghttp.CombineHandlers(
				ghttp.RespondWith(500, ""),
				ghttp.VerifyRequest("GET", "/oauth/clients/clientid"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			cm := &ClientManager{httpClient, config}
			_, err := cm.Get("clientid")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})

		It("returns helpful error when /oauth/clients/clientid response can't be parsed", func() {
			server.RouteToHandler("GET", "/oauth/clients/clientid", ghttp.CombineHandlers(
				ghttp.RespondWith(200, "{unparsable-json-response}"),
				ghttp.VerifyRequest("GET", "/oauth/clients/clientid"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			cm := &ClientManager{httpClient, config}
			_, err := cm.Get("clientid")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
			Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
		})
	})


})
