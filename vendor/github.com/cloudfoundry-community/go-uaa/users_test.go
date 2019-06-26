package uaa_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

const userResponse string = `{
	  "id" : "00000000-0000-0000-0000-000000000001",
	  "externalID" : "marcus-user",
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
		"userID" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientID" : "shinyclient",
		"scope" : "philosophy.read",
		"status" : "APPROVED",
		"lastUpdatedAt" : "2017-08-15T16:54:15.765Z",
		"expiresAt" : "2017-08-15T16:54:25.765Z"
	  }, {
		"userID" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70",
		"clientID" : "identity",
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
	  "zoneID" : "uaa",
	  "passwordLastModified" : "2017-08-15T16:54:15.000Z",
	  "previousLogonTime" : 1502816055768,
	  "lastLogonTime" : 1502816055768,
	  "schemas" : [ "urn:scim:schemas:core:1.0" ]
	}`

var userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)
var testUserValue uaa.User = uaa.User{
	ID:       "00000000-0000-0000-0000-000000000001",
	Username: "marcus@stoicism.com",
	Active:   newTrueP(),
	Name:     &uaa.UserName{GivenName: "Marcus", FamilyName: "Aurelius"},
}

var testUserJSON string = `{ "id": "00000000-0000-0000-0000-000000000001", "userName": "marcus@stoicism.com", "active": true, "name" : { "familyName" : "Aurelius", "givenName" : "Marcus" }}`

func newTrueP() *bool {
	b := true
	return &b
}

func newFalseP() *bool {
	b := false
	return &b
}

func testUsers(t *testing.T, when spec.G, it spec.S) {
	var (
		s       *httptest.Server
		handler http.Handler
		called  int
		a       *uaa.API
	)

	it.Before(func() {
		RegisterTestingT(t)
		called = 0
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			called = called + 1
			Expect(handler).NotTo(BeNil())
			handler.ServeHTTP(w, req)
		}))
		c := &http.Client{Transport: http.DefaultTransport}
		u, _ := url.Parse(s.URL)
		a = &uaa.API{
			TargetURL:             u,
			AuthenticatedClient:   c,
			UnauthenticatedClient: c,
		}
	})

	it.After(func() {
		if s != nil {
			s.Close()
		}
	})

	when("GetUserByUsername()", func() {
		when("no username is specified", func() {
			it("returns an error", func() {
				u, err := a.GetUserByUsername("", "", "")
				Expect(err).To(HaveOccurred())
				Expect(u).To(BeNil())
				Expect(err.Error()).To(Equal("username cannot be blank"))
			})
		})

		when("an origin is specified", func() {
			it("looks up a user with SCIM filter", func() {
				user := uaa.User{Username: "marcus", Origin: "uaa"}
				response := PaginatedResponse(user)
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus" and origin eq "uaa"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				})

				u, err := a.GetUserByUsername("marcus", "uaa", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(u.Username).To(Equal("marcus"))
			})

			it("returns an error when request fails", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus" and origin eq "uaa"`))
					w.WriteHeader(http.StatusInternalServerError)
				})

				_, err := a.GetUserByUsername("marcus", "uaa", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			it("returns an error if no results are found", func() {
				response := PaginatedResponse()
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus" and origin eq "uaa"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				})
				_, err := a.GetUserByUsername("marcus", "uaa", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`user marcus not found in origin uaa`))
			})

			when("attributes are specified", func() {
				it("adds them to the GET request", func() {
					handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						Expect(req.Header.Get("Accept")).To(Equal("application/json"))
						Expect(req.URL.Path).To(Equal("/Users"))
						Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus" and origin eq "uaa"`))
						Expect(req.URL.Query().Get("attributes")).To(Equal(`userName,emails`))
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(PaginatedResponse(uaa.User{Username: "marcus", Origin: "uaa"})))
					})
					_, err := a.GetUserByUsername("marcus", "uaa", "userName,emails")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		when("no origin is specified", func() {
			it("looks up a user with a SCIM filter", func() {
				user := uaa.User{Username: "marcus", Origin: "uaa"}
				response := PaginatedResponse(user)
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				})
				u, err := a.GetUserByUsername("marcus", "", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(u.Username).To(Equal("marcus"))
			})

			it("returns an error when request fails", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus"`))
					w.WriteHeader(http.StatusInternalServerError)
				})
				_, err := a.GetUserByUsername("marcus", "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			it("returns an error when no users are found", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(PaginatedResponse()))
				})
				_, err := a.GetUserByUsername("marcus", "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`user marcus not found`))
			})

			it("returns an error when username found in multiple origins", func() {
				user1 := uaa.User{Username: "marcus", Origin: "uaa"}
				user2 := uaa.User{Username: "marcus", Origin: "ldap"}
				user3 := uaa.User{Username: "marcus", Origin: "okta"}
				response := PaginatedResponse(user1, user2, user3)
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal("/Users"))
					Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				})

				_, err := a.GetUserByUsername("marcus", "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Found users with username marcus in multiple origins [uaa, ldap, okta].`))
			})

			when("attributes are specified", func() {
				it("adds them to the GET request", func() {
					handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						Expect(req.Header.Get("Accept")).To(Equal("application/json"))
						Expect(req.URL.Path).To(Equal("/Users"))
						Expect(req.URL.Query().Get("filter")).To(Equal(`userName eq "marcus"`))
						Expect(req.URL.Query().Get("attributes")).To(Equal(`userName,emails`))
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(PaginatedResponse(uaa.User{Username: "marcus", Origin: "uaa"})))
					})
					_, err := a.GetUserByUsername("marcus", "", "userName,emails")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})

	when("ListAllUsers()", func() {
		it("can return multiple pages", func() {
			page1 := MultiPaginatedResponse(1, 1, 2, uaa.User{Username: "marcus", Origin: "uaa"})
			page2 := MultiPaginatedResponse(2, 1, 2, uaa.User{Username: "drseuss", Origin: "uaa"})
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				w.WriteHeader(http.StatusOK)
				if called == 1 {
					Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
					Expect(req.URL.Query().Get("count")).To(Equal("100"))
					w.Write([]byte(page1))
				} else {
					Expect(req.URL.Query().Get("startIndex")).To(Equal("2"))
					Expect(req.URL.Query().Get("count")).To(Equal("1"))
					w.Write([]byte(page2))
				}
			})

			users, err := a.ListAllUsers("", "", "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(users[0].Username).To(Equal("marcus"))
			Expect(users[1].Username).To(Equal("drseuss"))
			Expect(called).To(Equal(2))
		})
	})

	when("ListUsers()", func() {
		var userListResponse = fmt.Sprintf(PaginatedResponseTmpl, MarcusUserResponse, DrSeussUserResponse)

		it("can accept a filter query to limit results", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 1, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("does not include the filter param if no filter exists", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(""))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers("", "", "", "", 1, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("can accept an attributes list", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`))
				Expect(req.URL.Query().Get("attributes")).To(Equal(`userName,emails`))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "userName,emails", "", 1, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("can accept sortBy", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(""))
				Expect(req.URL.Query().Get("attributes")).To(Equal(""))
				Expect(req.URL.Query().Get("sortBy")).To(Equal("userName"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers("", "userName", "", "", 1, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("can accept sort order ascending/descending", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(""))
				Expect(req.URL.Query().Get("attributes")).To(Equal(""))
				Expect(req.URL.Query().Get("sortBy")).To(Equal(""))
				Expect(req.URL.Query().Get("sortOrder")).To(Equal("ascending"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers("", "", "", uaa.SortAscending, 1, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("uses a startIndex of 1 if 0 is supplied", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userListResponse))
			})
			userList, _, err := a.ListUsers("", "", "", "", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(userList[0].Username).To(Equal("marcus@stoicism.com"))
			Expect(userList[1].Username).To(Equal("drseuss@whoville.com"))
		})

		it("returns an error when /Users doesn't respond", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`))
				w.WriteHeader(http.StatusInternalServerError)
			})

			userList, _, err := a.ListUsers(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 1, 100)
			Expect(err).To(HaveOccurred())
			Expect(userList).To(BeNil())
		})

		it("returns an error when response is unparseable", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal("/Users"))
				Expect(req.URL.Query().Get("count")).To(Equal("100"))
				Expect(req.URL.Query().Get("startIndex")).To(Equal("1"))
				Expect(req.URL.Query().Get("filter")).To(Equal(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{unparsable}"))
			})
			userList, _, err := a.ListUsers(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 1, 100)
			Expect(err).To(HaveOccurred())
			Expect(userList).To(BeNil())
		})
	})

	when("ActivateUser()", func() {
		it("returns an error when the userID is empty", func() {
			err := a.ActivateUser("", 10)
			Expect(err).To(HaveOccurred())
			Expect(called).To(Equal(0))
		})

		it("activates the user using the userID", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Header.Get("If-Match")).To(Equal("10"))
				Expect(req.Method).To(Equal(http.MethodPatch))
				Expect(req.URL.Path).To(Equal("/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{ "active": true }`))
				w.WriteHeader(http.StatusOK)
			})
			err := a.ActivateUser("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", 10)
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(Equal(1))
		})

		it("returns a helpful error the request fails", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Header.Get("If-Match")).To(Equal("0"))
				Expect(req.Method).To(Equal(http.MethodPatch))
				Expect(req.URL.Path).To(Equal("/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{ "active": true }`))
				w.WriteHeader(http.StatusInternalServerError)
			})
			err := a.ActivateUser("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			Expect(called).To(Equal(1))
		})
	})

	when("DeactivateUser()", func() {
		it("returns an error when the userID is empty", func() {
			err := a.DeactivateUser("", 10)
			Expect(err).To(HaveOccurred())
			Expect(called).To(Equal(0))
		})

		it("deactivates the user using the userID", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Header.Get("If-Match")).To(Equal("10"))
				Expect(req.Method).To(Equal(http.MethodPatch))
				Expect(req.URL.Path).To(Equal("/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{ "active": false }`))
				w.WriteHeader(http.StatusOK)
			})
			err := a.DeactivateUser("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", 10)
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(Equal(1))
		})

		it("returns a helpful error the request fails", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Header.Get("If-Match")).To(Equal("0"))
				Expect(req.Method).To(Equal(http.MethodPatch))
				Expect(req.URL.Path).To(Equal("/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{ "active": false }`))
				w.WriteHeader(http.StatusInternalServerError)
			})
			err := a.DeactivateUser("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
			Expect(called).To(Equal(1))
		})
	})

	when("using user structs", func() {
		when("verified", func() {
			it("correctly shows false boolean values", func() {
				user := uaa.User{Verified: newFalseP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"verified": false}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Verified).To(BeFalse())
			})

			it("correctly shows true values", func() {
				user := uaa.User{Verified: newTrueP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"verified": true}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Verified).To(BeTrue())
			})

			it("correctly hides unset values", func() {
				user := uaa.User{}
				json.Unmarshal([]byte("{}"), &user)
				Expect(user.Verified).To(BeNil())

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{}`))
			})
		})

		when("emails", func() {
			it("correctly shows false boolean values", func() {
				user := uaa.User{}
				email := uaa.Email{Value: "foo@bar.com", Primary: newFalseP()}
				user.Emails = []uaa.Email{email}

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"emails": [ { "value": "foo@bar.com", "primary": false } ]}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Emails[0].Primary).To(BeFalse())
			})

			it("correctly shows true values", func() {
				user := uaa.User{}
				email := uaa.Email{Value: "foo@bar.com", Primary: newTrueP()}
				user.Emails = []uaa.Email{email}

				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"emails": [ { "value": "foo@bar.com", "primary": true } ]}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Emails[0].Primary).To(BeTrue())
			})
		})

		when("active", func() {
			it("correctly shows false boolean values", func() {
				user := uaa.User{Active: newFalseP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"active": false}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Active).To(BeFalse())
			})

			it("correctly shows true values", func() {
				user := uaa.User{Active: newTrueP()}
				userBytes, _ := json.Marshal(&user)
				Expect(string(userBytes)).To(MatchJSON(`{"active": true}`))

				newUser := uaa.User{}
				json.Unmarshal([]byte(userBytes), &newUser)
				Expect(*newUser.Active).To(BeTrue())
			})
		})
	})
}
