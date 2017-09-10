package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
	"net/url"
)

var _ = Describe("AuthcodeClientImpersonator", func() {
	var (
		impersonator AuthcodeClientImpersonator
		logger       utils.Logger
		httpClient   *http.Client
		config       uaa.Config
		launcher     TestLauncher
		uaaServer    *Server
	)

	BeforeEach(func() {
		httpClient = &http.Client{}
		launcher = TestLauncher{}
		uaaServer = NewServer()
		config = uaa.NewConfigWithServerURL(uaaServer.URL())
		logger = utils.NewLogger(GinkgoWriter, GinkgoWriter, GinkgoWriter, GinkgoWriter)
	})

	Describe("NewAuthcodeClientImpersonator", func() {
		BeforeEach(func() {
			launcher := TestLauncher{}
			impersonator = NewAuthcodeClientImpersonator(httpClient, config, "authcodeClientId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)
		})

		Describe("configures an AuthCallbackListener", func() {
			It("with appropriate static content", func() {
				Expect(impersonator.AuthCallbackServer.CSS()).To(ContainSubstring("Source Sans Pro"))
				Expect(impersonator.AuthCallbackServer.Html()).To(ContainSubstring("Authorization Code Grant: Success"))
			})

			It("with the desired port", func() {
				Expect(impersonator.AuthCallbackServer.Port()).To(Equal(8080))
			})

			It("with its logger", func() {
				Expect(impersonator.AuthCallbackServer.Log()).NotTo(Equal(utils.Logger{}))
				Expect(impersonator.AuthCallbackServer.Log()).To(Equal(logger))
			})

			It("with hangup func that looks for code in query params", func() {
				done := make(chan url.Values)

				urlParams := url.Values{}
				urlParams.Add("code", "56575db17b164e568668c0085ed14ae1")
				go impersonator.AuthCallbackServer.Hangup(done, urlParams)

				Expect(<-done).To(Equal(urlParams))
			})
		})
	})

	Describe("#Start", func() {
		It("starts the AuthCallbackServer", func() {
			uaaServer.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				VerifyRequest("POST", "/oauth/token"),
				RespondWith(http.StatusOK, `{
				  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
				  "token_type" : "bearer",
				  "expires_in" : 3000,
				  "scope" : "openid",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`),
				VerifyFormKV("client_id", "authcodeId"),
				VerifyFormKV("client_secret", "authcodesecret"),
				VerifyFormKV("grant_type", "authorization_code"),
				VerifyFormKV("token_format", "jwt"),
				VerifyFormKV("response_type", "token"),
				VerifyFormKV("code", "secretcode"),
				VerifyFormKV("redirect_uri", "http://localhost:8080")),
			)
			impersonator = NewAuthcodeClientImpersonator(httpClient, config, "authcodeId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)

			// Start the callback server
			go impersonator.Start()

			// Hit the callback server with an authcode
			httpClient.Get("http://localhost:8080/?code=secretcode")

			// The callback server should have exchanged the code for a token
			tokenResponse := <-impersonator.Done()
			Expect(tokenResponse.AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			Expect(tokenResponse.TokenType).To(Equal("bearer"))
			Expect(tokenResponse.Scope).To(Equal("openid"))
			Expect(tokenResponse.JTI).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			Expect(tokenResponse.ExpiresIn).To(Equal(int32(3000)))
		})
	})

	Describe("#Authorize", func() {
		It("launches a browser to the authorize page", func() {
			impersonator = NewAuthcodeClientImpersonator(httpClient, config, "authcodeId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)

			impersonator.Authorize()

			Expect(launcher.TargetUrl).To(Equal(uaaServer.URL() + "/oauth/authorize?client_id=authcodeId&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=code"))
		})
	})
})
