package uaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/jhamon/uaa-cli/uaa"
	"net/http"
	"net/url"
	"strconv"
	"io/ioutil"
)

var _ = Describe("HttpRequestFactory", func() {
	var (
		factory HttpRequestFactory
		context UaaContext
		req *http.Request
	)

	ItBuildsUrlsFromUaaContext := func() {
		Describe("Get", func() {
			It("builds a GET request", func() {
				context.BaseUrl = "http://localhost.com"

				req, _ = factory.Get(context, "foo", "")

				Expect(req.Method).To(Equal("GET"))
			})

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
		})

		Describe("Post", func() {
			It("builds a POST request", func() {
				context.BaseUrl = "http://localhost.com"

				req, _ = factory.Post(context, "foo", "", &url.Values{})

				Expect(req.Method).To(Equal("POST"))
			})

			It("sets an Accept header", func() {
				context.BaseUrl = "http://localhost.com"
				req, _ = factory.Post(context, "foo", "", &url.Values{})
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			})

			It("sets Content-Type header", func() {
				context.BaseUrl = "http://localhost.com"
				req, _ = factory.Post(context, "foo", "", &url.Values{})
				Expect(req.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))
			})

			It("sets a url-encoded body and Content-Length header", func() {
				context.BaseUrl = "http://localhost.com"
				data := url.Values{}
				data.Add("client_id", "login")
				data.Add("client_secret", "loginsecret")
				data.Add("grant_type", "client_credentials")
				data.Add("token_format", "opaque")
				data.Add("response_type", "token")

				req, _ = factory.Post(context, "foo", "", &data)
				Expect(req.Header.Get("Content-Length")).To(Equal(strconv.Itoa(len(data.Encode()))))
				reqBody, _ := ioutil.ReadAll(req.Body)
				Expect(string(reqBody)).To(ContainSubstring("client_id=login"))
				Expect(string(reqBody)).To(ContainSubstring("client_secret=loginsecret"))
				Expect(string(reqBody)).To(ContainSubstring("grant_type=client_credentials"))
				Expect(string(reqBody)).To(ContainSubstring("token_format=opaque"))
				Expect(string(reqBody)).To(ContainSubstring("response_type=token"))
				Expect(string(reqBody)).To(HaveLen(len("client_id=login&client_secret=loginsecret&grant_type=client_credentials&token_format=opaque&response_type=token")))
			})

			It("builds requests from UaaContext", func() {
				context.BaseUrl = "http://localhost.com"

				req, _ = factory.Post(context, "foo", "", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://localhost.com/foo"))

				req, _ = factory.Post(context, "/foo", "scheme=openid", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://localhost.com/foo?scheme=openid"))
			})

			It("accepts a query string", func() {
				context.BaseUrl = "http://localhost.com"
				req, _ = factory.Post(context, "/foo", "scheme=openid&foo=bar", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://localhost.com/foo?scheme=openid&foo=bar"))
			})
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

		It("adds an Authorization header when GET", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Get(context, "foo", "")
			Expect(req.Header.Get("Authorization")).To(Equal("bearer access_token"))
		})

		It("adds an Authorization header when POST", func() {
			context.BaseUrl = "http://localhost.com"
			req, _ = factory.Post(context, "foo", "", &url.Values{})
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