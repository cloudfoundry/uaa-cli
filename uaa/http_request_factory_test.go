package uaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/jhamon/uaa-cli/uaa"
	"net/http"
)

var _ = Describe("HttpRequestFactory", func() {
	var (
		factory HttpRequestFactory
		context UaaContext
		req *http.Request
	)

	ItBuildsUrlsFromUaaContext := func() {
		It("builds requests from UaaContext", func() {
			context.BaseUrl = "http://localhost.com"

			req, _ = factory.Get(context, "foo", "")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))

			req, _ = factory.Get(context, "/foo", "scheme=openid")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo?scheme=openid"))
		})

		It("sets an Accept header", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "foo", "")
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
		})

		It("handles path slashes", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "/foo", "")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))

			context.BaseUrl = "http://localhost.com/"
			req, _ = factory.Get(context, "foo", "")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))

			context.BaseUrl = "http://localhost.com/"
			req, _ = factory.Get(context, "/foo", "")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))

			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "foo", "")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))
		})

		It("accepts a query string", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "/foo", "scheme=openid&foo=bar")
			Expect(req.URL.String()).To(Equal("http://localhost.com/foo?scheme=openid&foo=bar"))
		})
	}

	Describe("UnauthenticatedRequestFactory", func() {
		BeforeEach(func() {
			factory = UnauthenticatedRequestFactory{}
			context = UaaContext{}
		})

		ItBuildsUrlsFromUaaContext()
	})

	Describe("AuthenticatedRequestFactory", func() {
		BeforeEach(func() {
			factory = AuthenticatedRequestFactory{}
			context = UaaContext{}
			context.AccessToken = "access_token"
		})

		ItBuildsUrlsFromUaaContext()

		It("adds an Authorization header", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "foo", "")
			Expect(req.Header.Get("Authorization")).To(Equal("bearer access_token"))
		})

		It("returns an error when context has no token", func() {
			context.BaseUrl = "http://localhost.com"
			context.AccessToken = ""
			_, err := factory.Get(context, "foo", "")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("An access token is required to call http://localhost.com/foo"))
		})

	})
})