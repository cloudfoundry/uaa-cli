package uaa_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

const groupResponse string = `{
	"id" : "00000000-0000-0000-0000-000000000001",
	"meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	},
	"displayName" : "cloud_controller.read",
	"description" : "View details of your applications and services",
	"members" : [ {
		"origin" : "uaa",
		"type" : "USER",
		"value" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"
	} ],
	"zoneID" : "uaa",
	"schemas" : [ "urn:scim:schemas:core:1.0" ]
}`

var groupListResponse = fmt.Sprintf(PaginatedResponseTmpl, UaaAdminGroupResponse, CloudControllerReadGroupResponse)

const CloudControllerReadGroupResponse string = `{
	"id" : "00000000-0000-0000-0000-000000000002",
	"meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	},
	"displayName" : "cloud_controller.read",
	"description" : "View details of your applications and services",
	"members" : [ {
		"origin" : "uaa",
		"type" : "USER",
		"value" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"
	} ],
	"zoneID" : "uaa",
	"schemas" : [ "urn:scim:schemas:core:1.0" ]
}`

const UaaAdminGroupResponse string = `{
	"id" : "00000000-0000-0000-0000-000000000001",
	"meta" : {
		"version" : 1,
		"created" : "2017-01-15T16:54:15.677Z",
		"lastModified" : "2017-08-15T16:54:15.677Z"
	},
	"displayName" : "uaa.admin",
	"description" : "Act as an administrator throughout the UAA",
	"members" : [ {
		"origin" : "uaa",
		"type" : "USER",
		"value" : "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"
	} ],
	"zoneID" : "uaa",
	"schemas" : [ "urn:scim:schemas:core:1.0" ]
}`

var testGroupValue uaa.Group = uaa.Group{
	ID:          "00000000-0000-0000-0000-000000000001",
	DisplayName: "uaa.admin",
}

const testGroupJSON string = `{ "id" : "00000000-0000-0000-0000-000000000001", "displayName": "uaa.admin" }`

func testGroupsExtra(t *testing.T, when spec.G, it spec.S) {
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

	when("GetGroupByName()", func() {
		when("when no group name is specified", func() {
			it("returns an error", func() {
				_, err := a.GetGroupByName("", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("group name may not be blank"))
			})
		})

		when("when no origin is specified", func() {
			it("looks up a group with a SCIM filter", func() {
				group := uaa.Group{DisplayName: "uaa.admin"}
				response := PaginatedResponse(group)
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal(uaa.GroupsEndpoint))
					Expect(req.URL.Query().Get("filter")).To(Equal(`displayName eq "uaa.admin"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(response))
				})

				g, err := a.GetGroupByName("uaa.admin", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(g.DisplayName).To(Equal("uaa.admin"))
			})

			it("returns an error when request fails", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal(uaa.GroupsEndpoint))
					Expect(req.URL.Query().Get("filter")).To(Equal(`displayName eq "uaa.admin"`))
					w.WriteHeader(http.StatusInternalServerError)
				})

				_, err := a.GetGroupByName("uaa.admin", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			it("returns an error when no groups are found", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal(uaa.GroupsEndpoint))
					Expect(req.URL.Query().Get("filter")).To(Equal(`displayName eq "uaa.admin"`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(PaginatedResponse()))
				})

				_, err := a.GetGroupByName("uaa.admin", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`group uaa.admin not found`))
			})
		})

		when("when attributes are specified", func() {
			it("adds them to the GET request", func() {
				handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					Expect(req.Header.Get("Accept")).To(Equal("application/json"))
					Expect(req.URL.Path).To(Equal(uaa.GroupsEndpoint))
					Expect(req.URL.Query().Get("filter")).To(Equal(`displayName eq "uaa.admin"`))
					Expect(req.URL.Query().Get("attributes")).To(Equal(`displayName`))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(PaginatedResponse(uaa.Group{DisplayName: "uaa.admin"})))
				})
				_, err := a.GetGroupByName("uaa.admin", "displayName")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	when("AddGroupMember()", func() {
		it("adds a membership", func() {
			membershipJSON := `{"origin":"uaa","type":"USER","value":"user-id-1"}`
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal(fmt.Sprintf("%s/%s/members", uaa.GroupsEndpoint, "group-id-1")))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(membershipJSON))
			})
			err := a.AddGroupMember("group-id-1", "user-id-1", "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(Equal(1))
		})
	})

	when("RemoveGroupMember", func() {
		it("removes a membership", func() {
			membershipJSON := `{"origin":"uaa","type":"USER","value":"user-id-1"}`
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal(fmt.Sprintf("%s/%s/members", uaa.GroupsEndpoint, "group-id-1")))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(membershipJSON))
			})
			err := a.AddGroupMember("group-id-1", "user-id-1", "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(Equal(1))

			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.URL.Path).To(Equal(fmt.Sprintf("%s/%s/members/%s", uaa.GroupsEndpoint, "group-id-1", "user-id-1")))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(membershipJSON))
			})
			err = a.RemoveGroupMember("group-id-1", "user-id-1", "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(Equal(2))
		})
	})
}
