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

type TestData struct {
	Field1 string
	Field2 string
}

var _ = Describe("HttpRequestFactory", func() {
	var (
		factory HttpRequestFactory
		context UaaContext
		req *http.Request
		config Config
	)

	ItBuildsUrlsFromUaaContext := func() {
		Describe("Get", func() {
			It("builds a GET request", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "foo", "")

				Expect(req.Method).To(Equal("GET"))
			})

			It("builds requests from UaaContext", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "foo", "")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				req, _ = factory.Get(config.GetActiveTarget(), "/foo", "scheme=openid")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid"))
			})

			It("sets an Accept header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "foo", "")
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			})

			It("handles path slashes", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "/foo", "")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				config = NewConfigWithServerURL("http://www.localhost.com/")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "foo", "")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				config = NewConfigWithServerURL("http://www.localhost.com/")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "/foo", "")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "foo", "")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))
			})

			It("accepts a query string", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.Get(config.GetActiveTarget(), "/foo", "scheme=openid&foo=bar")
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid&foo=bar"))
			})
		})

		Describe("PostForm", func() {
			It("builds a POST request", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &url.Values{})

				Expect(req.Method).To(Equal("POST"))
			})

			It("sets an Accept header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &url.Values{})
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			})

			It("sets Content-Type header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &url.Values{})
				Expect(req.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))
			})

			It("sets a url-encoded body and Content-Length header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				data := url.Values{}
				data.Add("client_id", "login")
				data.Add("client_secret", "loginsecret")
				data.Add("grant_type", "client_credentials")
				data.Add("token_format", "opaque")
				data.Add("response_type", "token")

				req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &data)
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
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				req, _ = factory.PostForm(config.GetActiveTarget(), "/foo", "scheme=openid", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid"))
			})

			It("accepts a query string", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostForm(config.GetActiveTarget(), "/foo", "scheme=openid&foo=bar", &url.Values{})
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid&foo=bar"))
			})
		})

		Describe("PostJson", func() {
			var dataToPost TestData
			BeforeEach(func() {
				dataToPost = TestData{Field1:"foo", Field2: "bar"}
			})

			It("builds a POST request", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.PostJson(config.GetActiveTarget(), "/foo", "", dataToPost)

				Expect(req.Method).To(Equal("POST"))
			})

			It("sets an Accept header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostJson(config.GetActiveTarget(), "foo", "", dataToPost)
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			})

			It("sets Content-Type header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostJson(config.GetActiveTarget(), "foo", "", dataToPost)
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("sets a json body and Content-Length header", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				expectedBody := `{"Field1":"foo","Field2":"bar"}`

				req, _ = factory.PostJson(config.GetActiveTarget(), "foo", "", dataToPost)
				Expect(req.Header.Get("Content-Length")).To(Equal(strconv.Itoa(len(expectedBody))))
				reqBody, _ := ioutil.ReadAll(req.Body)
				Expect(string(reqBody)).To(MatchJSON(expectedBody))
			})

			It("builds requests from UaaContext", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})

				req, _ = factory.PostJson(config.GetActiveTarget(), "foo", "", dataToPost)
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo"))

				req, _ = factory.PostJson(config.GetActiveTarget(), "/foo", "scheme=openid", dataToPost)
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid"))
			})

			It("accepts a query string", func() {
				config = NewConfigWithServerURL("http://www.localhost.com")
				config.AddContext(UaaContext{AccessToken: "access_token"})
				req, _ = factory.PostJson(config.GetActiveTarget(), "/foo", "scheme=openid&foo=bar", dataToPost)
				Expect(req.URL.String()).To(Equal("http://www.localhost.com/foo?scheme=openid&foo=bar"))
			})
		})
	}

	Describe("UnauthenticatedRequestFactory", func() {
		BeforeEach(func() {
			factory = UnauthenticatedRequestFactory{}
			config = NewConfigWithServerURL("http://www.localhost.com")
			context = UaaContext{}
			config.AddContext(context)

		})

		ItBuildsUrlsFromUaaContext()
	})

	Describe("AuthenticatedRequestFactory", func() {
		BeforeEach(func() {
			factory = AuthenticatedRequestFactory{}
			config = NewConfigWithServerURL("http://www.localhost.com")
			context = UaaContext{AccessToken: "access_token"}
			config.AddContext(context)
		})

		ItBuildsUrlsFromUaaContext()

		It("adds an Authorization header when GET", func() {
			req, _ = factory.Get(config.GetActiveTarget(), "foo", "")
			Expect(req.Header.Get("Authorization")).To(Equal("bearer access_token"))
		})

		It("adds an Authorization header when POST", func() {
			req, _ = factory.PostForm(config.GetActiveTarget(), "foo", "", &url.Values{})
			Expect(req.Header.Get("Authorization")).To(Equal("bearer access_token"))
		})


		It("returns an error when context has no token", func() {
			config = NewConfigWithServerURL("http://www.localhost.com")
			context.AccessToken = ""
			config.AddContext(context)
			_, err := factory.Get(config.GetActiveTarget(), "foo", "")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("An access token is required to call http://www.localhost.com/foo"))
		})

	})
})