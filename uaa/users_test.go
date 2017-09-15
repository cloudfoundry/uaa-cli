package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "code.cloudfoundry.org/uaa-cli/fixtures"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Users", func() {
	var (
		um        UserManager
		uaaServer *ghttp.Server
	)

	BeforeEach(func() {
		uaaServer = ghttp.NewServer()
		config := NewConfigWithServerURL(uaaServer.URL())
		config.AddContext(NewContextWithToken("access_token"))
		um = UserManager{&http.Client{}, config}
	})

	var userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)

	Describe("UserManager#Get", func() {
		It("gets a user from the UAA by id", func() {
			uaaServer.RouteToHandler("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			user, _ := um.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(user.Id).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
			Expect(user.ExternalId).To(Equal("marcus-user"))
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

	Describe("UserManager#GetUserByUsername", func() {
		Context("when no username is specified", func() {
			It("returns an error", func() {
				_, err := um.GetUserByUsername("", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Username may not be blank."))
			})
		})

		Context("when an origin is specified", func() {
			It("looks up a user with SCIM filter", func() {
				user := ScimUser{Username: "marcus", Origin: "uaa"}
				response := PaginatedResponse(user)

				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, response),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22+and+origin+eq+%22uaa%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				user, err := um.GetUserByUsername("marcus", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Username).To(Equal("marcus"))
			})

			It("returns an error when request fails", func() {
				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, ""),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22+and+origin+eq+%22uaa%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetUserByUsername("marcus", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			It("returns an error if no results are found", func() {
				response := PaginatedResponse()

				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, response),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22+and+origin+eq+%22uaa%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetUserByUsername("marcus", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`User marcus not found in origin uaa`))
			})
		})

		Context("when no origin is specified", func() {
			It("looks up a user with a SCIM filter", func() {
				user := ScimUser{Username: "marcus", Origin: "uaa"}
				response := PaginatedResponse(user)

				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, response),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				user, err := um.GetUserByUsername("marcus", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Username).To(Equal("marcus"))
			})

			It("returns an error when request fails", func() {
				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, ""),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetUserByUsername("marcus", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			It("returns an error when no users are found", func() {
				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, PaginatedResponse()),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetUserByUsername("marcus", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`User marcus not found.`))
			})

			It("returns an error when username found in multiple origins", func() {
				user1 := ScimUser{Username: "marcus", Origin: "uaa"}
				user2 := ScimUser{Username: "marcus", Origin: "ldap"}
				user3 := ScimUser{Username: "marcus", Origin: "okta"}
				response := PaginatedResponse(user1, user2, user3)

				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, response),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetUserByUsername("marcus", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Found users with username marcus in multiple origins [uaa, ldap, okta].`))
			})
		})
	})

	Describe("UserManager#List", func() {
		It("can accept a filter query to limit results", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			resp, err := um.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(resp.Resources[1].Username).To(Equal("drseuss@whoville.com"))
		})

		It("gets all users when no filter is passed", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", ""),
			))

			resp, err := um.List("", "", "", "", 0, 0)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(resp.Resources[1].Username).To(Equal("drseuss@whoville.com"))
		})

		It("can accept an attributes list", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22&attributes=userName%2Cemails"),
			))

			resp, err := um.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "userName,emails", "", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(resp.Resources[1].Username).To(Equal("drseuss@whoville.com"))
		})

		It("can accept sortBy", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "sortBy=userName"),
			))

			_, err := um.List("", "userName", "", "", 0, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept count", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "count=10"),
			))

			_, err := um.List("", "", "", "", 0, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept sort order ascending/descending", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "sortOrder=ascending"),
			))

			_, err := um.List("", "", "", SORT_ASCENDING, 0, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept startIndex", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, userListResponse),
				ghttp.VerifyRequest("GET", "/Users", "startIndex=10"),
			))

			_, err := um.List("", "", "", "", 10, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error when /Users doesn't respond", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusInternalServerError, ""),
				ghttp.VerifyRequest("GET", "/Users", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := um.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when response is unparseable", func() {
			uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, "{unparsable}"),
				ghttp.VerifyRequest("GET", "/Users", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := um.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).To(HaveOccurred())
		})
	})

})
