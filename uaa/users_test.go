package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "code.cloudfoundry.org/uaa-cli/utils"
	"encoding/json"
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

	Describe("ScimUser Json", func() {
		// All this dance is necessary because the --attributes option means I need
		// to be able to hide values that aren't sent in the server response. If I
		// just used omitempty, I wouldn't be able to distinguish between empty values
		// (false, empty string) and ones that were never sent by the server.

		Describe("Verified", func() {
			It("correctly shows false boolean values", func() {
				user := ScimUser{Verified: NewFalseP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"verified": false}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Verified).To(BeFalse())
			})

			It("correctly shows true values", func() {
				user := ScimUser{Verified: NewTrueP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"verified": true}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Verified).To(BeTrue())
			})

			It("correctly hides unset values", func() {
				user := ScimUser{}
				json.Unmarshal([]byte("{}"), &user)
				Expect(user.Verified).To(BeNil())

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{}`))
			})
		})

		Describe("Active", func() {
			It("correctly shows false boolean values", func() {
				user := ScimUser{Active: NewFalseP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"active": false}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Active).To(BeFalse())
			})

			It("correctly shows true values", func() {
				user := ScimUser{Active: NewTrueP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"active": true}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Active).To(BeTrue())
			})
		})

		Describe("Emails", func() {
			It("correctly shows false boolean values", func() {
				user := ScimUser{}
				email := ScimUserEmail{"foo@bar.com", NewFalseP()}
				user.Emails = []ScimUserEmail{email}

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"emails": [ { "value": "foo@bar.com", "primary": false } ]}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Emails[0].Primary).To(BeFalse())
			})

			It("correctly shows true values", func() {
				user := ScimUser{}
				email := ScimUserEmail{"foo@bar.com", NewTrueP()}
				user.Emails = []ScimUserEmail{email}

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"emails": [ { "value": "foo@bar.com", "primary": true } ]}`))

				newUser := ScimUser{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Emails[0].Primary).To(BeTrue())
			})
		})
	})

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
			Expect(user.Meta.Created).To(Equal("2017-01-15T16:54:15.677Z"))
			Expect(user.Meta.LastModified).To(Equal("2017-08-15T16:54:15.677Z"))
			Expect(user.Meta.Version).To(Equal(1))
			Expect(user.Username).To(Equal("marcus@stoicism.com"))
			Expect(user.Name.GivenName).To(Equal("Marcus"))
			Expect(user.Name.FamilyName).To(Equal("Aurelius"))
			Expect(*user.Emails[0].Primary).To(Equal(false))
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
			Expect(*user.Active).To(Equal(true))
			Expect(*user.Verified).To(Equal(true))
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

	Describe("UserManager#GetByUsername", func() {
		Context("when no username is specified", func() {
			It("returns an error", func() {
				_, err := um.GetByUsername("", "", "")
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

				user, err := um.GetByUsername("marcus", "uaa", "")
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

				_, err := um.GetByUsername("marcus", "uaa", "")
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

				_, err := um.GetByUsername("marcus", "uaa", "")
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

				user, err := um.GetByUsername("marcus", "", "")
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

				_, err := um.GetByUsername("marcus", "", "")
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

				_, err := um.GetByUsername("marcus", "", "")
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

				_, err := um.GetByUsername("marcus", "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Found users with username marcus in multiple origins [uaa, ldap, okta].`))
			})
		})

		Context("when attributes are specified", func() {
			It("adds them to the GET request", func() {
				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, PaginatedResponse(ScimUser{Username: "marcus", Origin: "uaa"})),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22&attributes=userName%2Cemails"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetByUsername("marcus", "", "userName,emails")
				Expect(err).NotTo(HaveOccurred())
			})

			It("adds them to the GET request", func() {
				uaaServer.RouteToHandler("GET", "/Users", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, PaginatedResponse(ScimUser{Username: "marcus", Origin: "uaa"})),
					ghttp.VerifyRequest("GET", "/Users", "filter=userName+eq+%22marcus%22+and+origin+eq+%22uaa%22&attributes=userName%2Cemails"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := um.GetByUsername("marcus", "uaa", "userName,emails")
				Expect(err).NotTo(HaveOccurred())
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

	Describe("UserManager#Create", func() {
		var user ScimUser

		BeforeEach(func() {
			user = ScimUser{
				Username: "marcus@stoicism.com",
				Active: NewTrueP(),
			}
			user.Name = &ScimUserName{GivenName: "Marcus", FamilyName: "Aurelius"}
		})

		It("performs POST with user data and bearer token", func() {
			uaaServer.RouteToHandler("POST", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Users"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSON(`{ "userName": "marcus@stoicism.com", "active": true, "name" : { "familyName" : "Aurelius", "givenName" : "Marcus" }}`),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			um.Create(user)

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the created user", func() {
			uaaServer.RouteToHandler("POST", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Users"),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			user, _ := um.Create(user)

			Expect(user.Username).To(Equal("marcus@stoicism.com"))
			Expect(user.ExternalId).To(Equal("marcus-user"))
		})

		It("returns error when response cannot be parsed", func() {
			uaaServer.RouteToHandler("POST", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Users"),
				ghttp.RespondWith(http.StatusOK, "{unparseable}"),
			))

			_, err:= um.Create(user)

			Expect(err).To(HaveOccurred())
		})

		It("returns error when response is not 200 OK", func() {
			uaaServer.RouteToHandler("POST", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Users"),
				ghttp.RespondWith(http.StatusBadRequest, ""),
			))

			_, err:= um.Create(user)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("UserManager#Update", func() {
		var user ScimUser

		BeforeEach(func() {
			user = ScimUser{
				Username: "marcus@stoicism.com",
				Active: NewTrueP(),
			}
			user.Name = &ScimUserName{GivenName: "Marcus", FamilyName: "Aurelius"}
		})

		It("performs PUT with user data and bearer token", func() {
			uaaServer.RouteToHandler("PUT", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Users"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSON(`{ "userName": "marcus@stoicism.com", "active": true, "name" : { "familyName" : "Aurelius", "givenName" : "Marcus" }}`),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			um.Update(user)

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the updated user", func() {
			uaaServer.RouteToHandler("PUT", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Users"),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			user, _ := um.Update(user)

			Expect(user.Username).To(Equal("marcus@stoicism.com"))
			Expect(user.ExternalId).To(Equal("marcus-user"))
		})

		It("returns error when response cannot be parsed", func() {
			uaaServer.RouteToHandler("PUT", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Users"),
				ghttp.RespondWith(http.StatusOK, "{unparseable}"),
			))

			_, err:= um.Update(user)

			Expect(err).To(HaveOccurred())
		})

		It("returns error when response is not 200 OK", func() {
			uaaServer.RouteToHandler("PUT", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Users"),
				ghttp.RespondWith(http.StatusBadRequest, ""),
			))

			_, err:= um.Update(user)

			Expect(err).To(HaveOccurred())
		})
	})

})
