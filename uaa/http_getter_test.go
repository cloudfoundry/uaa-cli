package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("HttpGetter", func() {
	var (
		server *ghttp.Server
		context UaaContext
		responseJson string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		context = UaaContext{}
		context.BaseUrl = server.URL()
		responseJson = `{"foo": "bar"}`
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("UnauthenticatedGetter", func() {
		It("calls an endpoint with Accept application/json header", func() {
			server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
				ghttp.RespondWith(200, responseJson),
				ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
			))

			UnauthenticatedGetter{}.Get(context, "/testPath", "someQueryParam=true")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns helpful error when GET request fails", func() {
			server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
				ghttp.RespondWith(500, ""),
				ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
			))

			_, err := UnauthenticatedGetter{}.Get(context, "/testPath", "someQueryParam=true")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})
	})

	Describe("AuthenticatedGetter", func() {
		It("calls an endpoint with Accept and Authorization headers", func() {
			server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
				ghttp.RespondWith(200, responseJson),
				ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			context.AccessToken = "access_token"
			AuthenticatedGetter{}.Get(context, "/testPath", "someQueryParam=true")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns a helpful error when GET request fails", func() {
			server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
				ghttp.RespondWith(500, ""),
				ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
			))

			context.AccessToken = "access_token"
			_, err := AuthenticatedGetter{}.Get(context, "/testPath", "someQueryParam=true")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})

		It("returns a helpful error when no token in context", func() {
			server.RouteToHandler("GET", "/testPath", ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/testPath", "someQueryParam=true"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
			))

			context.AccessToken = ""
			_, err := AuthenticatedGetter{}.Get(context, "/testPath", "someQueryParam=true")

			Expect(server.ReceivedRequests()).To(HaveLen(0))
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An access token is required to call"))
		})
	})
})
