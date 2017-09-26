package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"
	"fmt"

	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Groups", func() {
	var (
		gm        GroupManager
		uaaServer *ghttp.Server
	)

	BeforeEach(func() {
		uaaServer = ghttp.NewServer()
		config := NewConfigWithServerURL(uaaServer.URL())
		config.AddContext(NewContextWithToken("access_token"))
		gm = GroupManager{&http.Client{}, config}
	})

	var groupListResponse = fmt.Sprintf(PaginatedResponseTmpl, AdminGroupResponse, ReadGroupResponse)

	Describe("GroupManager#Get", func() {
		It("gets a group from the UAA by id", func() {
			uaaServer.RouteToHandler("GET", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab", ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/Groups/05a0c169-3592-4a45-b109-a16d9246e0ab"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			group, _ := gm.Get("05a0c169-3592-4a45-b109-a16d9246e0ab")

			Expect(group.ID).To(Equal("05a0c169-3592-4a45-b109-a16d9246e0ab"))
			Expect(group.Meta.Created).To(Equal("2017-01-15T16:54:15.677Z"))
			Expect(group.Meta.LastModified).To(Equal("2017-08-15T16:54:15.677Z"))
			Expect(group.Meta.Version).To(Equal(1))
			Expect(group.DisplayName).To(Equal("admin"))
			Expect(group.ZoneID).To(Equal("uaa"))
			Expect(group.Description).To(Equal("admin"))
			Expect(group.Members[0].Origin).To(Equal("uaa"))
			Expect(group.Members[0].Type).To(Equal("USER"))
			Expect(group.Members[0].Value).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
			Expect(group.Schemas[0]).To(Equal("urn:scim:schemas:core:1.0"))
		})

		It("returns helpful error when /Users/userid request fails", func() {
			uaaServer.RouteToHandler("GET", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusInternalServerError, ""),
				ghttp.VerifyRequest("GET", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := gm.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})

		It("returns helpful error when /Groups/groupid response can't be parsed", func() {
			uaaServer.RouteToHandler("GET", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, "{unparsable-json-response}"),
				ghttp.VerifyRequest("GET", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := gm.Get("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
			Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
		})
	})

	Describe("GroupManager#GetByName", func() {
		Context("when no groupname is specified", func() {
			It("returns an error", func() {
				_, err := gm.GetByName("", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Groupname may not be blank."))
			})
		})

		Context("when no origin is specified", func() {
			It("looks up a user with a SCIM filter", func() {
				group := ScimGroup{DisplayName: "admin"}
				response := PaginatedResponse(group)

				uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, response),
					ghttp.VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				group, err := gm.GetByName("admin", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(group.DisplayName).To(Equal("admin"))
			})

			It("returns an error when request fails", func() {
				uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, ""),
					ghttp.VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := gm.GetByName("admin", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("An unknown error"))
			})

			It("returns an error when no users are found", func() {
				uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, PaginatedResponse()),
					ghttp.VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := gm.GetByName("admin", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Group admin not found.`))
			})
		})

		Context("when attributes are specified", func() {
			It("adds them to the GET request", func() {
				uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, PaginatedResponse(ScimGroup{DisplayName: "admin"})),
					ghttp.VerifyRequest("GET", "/Groups", "filter=displayName+eq+%22admin%22&attributes=displayName"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				))

				_, err := gm.GetByName("admin", "displayName")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GroupManager#List", func() {
		It("can accept a filter query to limit results", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			resp, err := gm.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].DisplayName).To(Equal("admin"))
			Expect(resp.Resources[1].DisplayName).To(Equal("read"))
		})

		It("gets all users when no filter is passed", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", ""),
			))

			resp, err := gm.List("", "", "", "", 0, 0)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].DisplayName).To(Equal("admin"))
			Expect(resp.Resources[1].DisplayName).To(Equal("read"))
		})

		It("can accept an attributes list", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22&attributes=displayName"),
			))

			resp, err := gm.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "displayName", "", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Resources[0].DisplayName).To(Equal("admin"))
			Expect(resp.Resources[1].DisplayName).To(Equal("read"))
		})

		It("can accept sortBy", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "sortBy=displayName"),
			))

			_, err := gm.List("", "displayName", "", "", 0, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept count", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "count=10"),
			))

			_, err := gm.List("", "", "", "", 0, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept sort order ascending/descending", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "sortOrder=ascending"),
			))

			_, err := gm.List("", "", "", SORT_ASCENDING, 0, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can accept startIndex", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, groupListResponse),
				ghttp.VerifyRequest("GET", "/Groups", "startIndex=10"),
			))

			_, err := gm.List("", "", "", "", 10, 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error when /Users doesn't respond", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusInternalServerError, ""),
				ghttp.VerifyRequest("GET", "/Groups", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := gm.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when response is unparseable", func() {
			uaaServer.RouteToHandler("GET", "/Groups", ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, "{unparsable}"),
				ghttp.VerifyRequest("GET", "/Groups", "filter=id+eq+%22fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7%22"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			_, err := gm.List(`id eq "fb5f32e1-5cb3-49e6-93df-6df9c8c8bd7"`, "", "", "", 0, 0)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GroupManager#Create", func() {
		var group ScimGroup

		BeforeEach(func() {
			group = ScimGroup{
				DisplayName: "admin",
			}
		})

		It("performs POST with user data and bearer token", func() {
			uaaServer.RouteToHandler("POST", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Groups"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSON(`{ "displayName": "admin", "members": null }`),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			gm.Create(group)

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the created user", func() {
			uaaServer.RouteToHandler("POST", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Groups"),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			group, _ := gm.Create(group)

			Expect(group.DisplayName).To(Equal("admin"))
		})

		It("returns error when response cannot be parsed", func() {
			uaaServer.RouteToHandler("POST", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Groups"),
				ghttp.RespondWith(http.StatusOK, "{unparseable}"),
			))

			_, err := gm.Create(group)

			Expect(err).To(HaveOccurred())
		})

		It("returns error when response is not 200 OK", func() {
			uaaServer.RouteToHandler("POST", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Groups"),
				ghttp.RespondWith(http.StatusBadRequest, ""),
			))

			_, err := gm.Create(group)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GroupManager#Update", func() {
		var group ScimGroup

		BeforeEach(func() {
			group = ScimGroup{
				DisplayName: "admin",
			}
		})

		It("performs PUT with user data and bearer token", func() {
			uaaServer.RouteToHandler("PUT", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Groups"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSON(`{ "displayName": "admin", "members": null }`),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			gm.Update(group)

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the updated group", func() {
			uaaServer.RouteToHandler("PUT", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Groups"),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			group, _ := gm.Update(group)

			Expect(group.DisplayName).To(Equal("admin"))
		})

		It("returns error when response cannot be parsed", func() {
			uaaServer.RouteToHandler("PUT", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Groups"),
				ghttp.RespondWith(http.StatusOK, "{unparseable}"),
			))

			_, err := gm.Update(group)

			Expect(err).To(HaveOccurred())
		})

		It("returns error when response is not 200 OK", func() {
			uaaServer.RouteToHandler("PUT", "/Groups", ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/Groups"),
				ghttp.RespondWith(http.StatusBadRequest, ""),
			))

			_, err := gm.Update(group)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GroupManager#Delete", func() {
		It("performs DELETE with user data and bearer token", func() {
			uaaServer.RouteToHandler("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			gm.Delete("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(uaaServer.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the deleted group", func() {
			uaaServer.RouteToHandler("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.RespondWith(http.StatusOK, AdminGroupResponse),
			))

			group, _ := gm.Delete("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(group.DisplayName).To(Equal("admin"))
		})

		It("returns error when response cannot be parsed", func() {
			uaaServer.RouteToHandler("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.RespondWith(http.StatusOK, "{unparseable}"),
			))

			_, err := gm.Delete("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(err).To(HaveOccurred())
		})

		It("returns error when response is not 200 OK", func() {
			uaaServer.RouteToHandler("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/Groups/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.RespondWith(http.StatusBadRequest, ""),
			))

			_, err := gm.Delete("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70")

			Expect(err).To(HaveOccurred())
		})
	})
})
