package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Me", func() {
	var (
		server       *ghttp.Server
		client       *http.Client
		config       Config
		userinfoJson string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = &http.Client{}
		config = NewConfigWithServerURL(server.URL())
		userinfoJson = `{
		  "user_id": "d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87",
		  "sub": "d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87",
		  "user_name": "charlieb",
		  "given_name": "Charlie",
		  "family_name": "Brown",
		  "email": "charlieb@peanuts.com",
		  "phone_number": null,
		  "previous_logon_time": 1503123277743,
		  "name": "Charlie Brown"
		}`
	})

	AfterEach(func() {
		server.Close()
	})

	It("calls the /userinfo endpoint", func() {
		server.RouteToHandler("GET", "/userinfo", ghttp.CombineHandlers(
			ghttp.RespondWith(200, userinfoJson),
			ghttp.VerifyRequest("GET", "/userinfo", "scheme=openid"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
			ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
		))

		config.AddContext(UaaContext{AccessToken: "access_token"})
		userinfo, _ := Me(client, config)

		Expect(server.ReceivedRequests()).To(HaveLen(1))
		Expect(userinfo.UserId).To(Equal("d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87"))
		Expect(userinfo.Sub).To(Equal("d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87"))
		Expect(userinfo.Username).To(Equal("charlieb"))
		Expect(userinfo.GivenName).To(Equal("Charlie"))
		Expect(userinfo.FamilyName).To(Equal("Brown"))

		Expect(userinfo.Email).To(Equal("charlieb@peanuts.com"))
	})

	It("returns helpful error when /userinfo request fails", func() {
		server.RouteToHandler("GET", "/userinfo", ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/userinfo", "scheme=openid"),
			ghttp.RespondWith(500, "error response"),
			ghttp.VerifyRequest("GET", "/userinfo"),
		))

		config.AddContext(UaaContext{AccessToken: "access_token"})
		_, err := Me(client, config)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
	})

	It("returns helpful error when /userinfo response can't be parsed", func() {
		server.RouteToHandler("GET", "/userinfo", ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/userinfo", "scheme=openid"),
			ghttp.RespondWith(200, "{unparsable-json-response}"),
			ghttp.VerifyRequest("GET", "/userinfo"),
		))

		config.AddContext(UaaContext{AccessToken: "access_token"})
		_, err := Me(client, config)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
		Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
	})
})
