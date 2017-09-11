package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Users", func() {
	var (
		um UserManager
		uaaServer *ghttp.Server
	)

	BeforeEach(func() {
		uaaServer = ghttp.NewServer()
		config := NewConfigWithServerURL(uaaServer.URL())
		config.AddContext(NewContextWithToken("access_token"))
		um = UserManager{&http.Client{}, config}
	})

	const userJsonResponse = `{
	  "id" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
	  "externalId" : "test-user",
	  "meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	  },
	  "userName" : "marcus@stoicism.com",
	  "name" : {
		"familyName" : "Aurelius",
		"givenName" : "Marcus"
	  },
	  "emails" : [ {
		"value" : "marcus@stoicism.com",
		"primary" : false
	  } ],
	  "groups" : [ {
		"value" : "ac2ab20e-0a2d-4b68-82e4-817ee6b258b4",
		"display" : "philosophy.read",
		"type" : "DIRECT"
	  }, {
		"value" : "110b2434-4a30-439b-b5fc-f4cf47fc04f0",
		"display" : "philosophy.write",
		"type" : "DIRECT"
	  }],
	  "approvals" : [ {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "shinyclient",
		"scope" : "philosophy.read",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:15.765Z",
		"expiresAt" : "2017-08-15T16:54:25.765Z"
	  }, {
		"userId" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientId" : "identity",
		"scope" : "uaa.user",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:45.767Z",
		"expiresAt" : "2017-08-15T16:54:45.767Z"
	  } ],
	  "phoneNumbers" : [ {
		"value" : "5555555555"
	  } ],
	  "active" : true,
	  "verified" : true,
	  "origin" : "uaa",
	  "zoneId" : "uaa",
	  "passwordLastModified" : "2017-08-15T16:54:15.000Z",
	  "previousLogonTime" : 1502816055768,
	  "lastLogonTime" : 1502816055768,
	  "schemas" : [ "urn:scim:schemas:core:1.0" ]
	}`

	Describe("UserManager#Get", func() {
		It("gets a user from the UAA", func() {
			uaaServer.RouteToHandler("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, userJsonResponse),
			))

			user, _ := um.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(user.Id).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
			Expect(user.ExternalId).To(Equal("test-user"))
			Expect(user.Active).To(Equal(true))
			Expect(user.Meta.Created).To(Equal("2017-01-15T16:54:15.677Z"))
			Expect(user.Meta.LastModified).To(Equal("2017-08-15T16:54:15.677Z"))
			Expect(user.Meta.Version).To(Equal(1))
			Expect(user.Username).To(Equal("marcus@stoicism.com"))
			Expect(user.Name.GivenName).To(Equal("Marcus"))
			Expect(user.Name.FamilyName).To(Equal("Aurelius"))
			Expect(user.Emails[0].Primary).To(Equal(false))
			Expect(user.Emails[0].Value).To(Equal("marcus@stoicism.com"))
			Expect(user.Groups[0].Display).To(Equal("philosophy.read"))
			Expect(user.Groups[0].Type).To(Equal("DIRECT"))
			Expect(user.Groups[0].Value).To(Equal("ac2ab20e-0a2d-4b68-82e4-817ee6b258b4"))
			Expect(user.Groups[1].Display).To(Equal("philosophy.write"))
			Expect(user.Groups[1].Type).To(Equal("DIRECT"))
			Expect(user.Groups[1].Value).To(Equal("110b2434-4a30-439b-b5fc-f4cf47fc04f0"))
			Expect(user.Approvals[0].UserId).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
			Expect(user.Approvals[0].ClientId).To(Equal("shinyclient"))
			Expect(user.Approvals[0].ExpiresAt).To(Equal("2017-08-15T16:54:25.765Z"))
			Expect(user.Approvals[0].LastUpdatedAt).To(Equal("2017-08-15T16:54:15.765Z"))
			Expect(user.Approvals[0].Scope).To(Equal("philosophy.read"))
			Expect(user.Approvals[0].Status).To(Equal("APPROVED"))
			Expect(user.PhoneNumbers[0].Value).To(Equal("5555555555"))
			Expect(user.Active).To(Equal(true))
			Expect(user.Verified).To(Equal(true))
			Expect(user.Origin).To(Equal("uaa"))
			Expect(user.ZoneId).To(Equal("uaa"))
			Expect(user.PasswordLastModified).To(Equal("2017-08-15T16:54:15.000Z"))
			Expect(user.PreviousLogonTime).To(Equal(1502816055768))
			Expect(user.LastLogonTime).To(Equal(1502816055768))
			Expect(user.Schemas[0]).To(Equal("urn:scim:schemas:core:1.0"))
		})

		It("returns helpful error when /Users/userid request fails", func() {
			uaaServer.RouteToHandler("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusInternalServerError, ""),
				ghttp.VerifyRequest("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := um.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})

		It("returns helpful error when /Users/userid response can't be parsed", func() {
			uaaServer.RouteToHandler("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, "{unparsable-json-response}"),
				ghttp.VerifyRequest("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := um.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
			Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
		})
	})

})
