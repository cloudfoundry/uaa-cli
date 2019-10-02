package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"net/http"
	"net/url"

	"time"

	cliConfig "code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("AuthcodeClientImpersonator", func() {
	var (
		impersonator AuthcodeClientImpersonator
		logger       Logger
		httpClient   *http.Client
		config       cliConfig.Config
		launcher     TestLauncher
		uaaServer    *Server
	)

	BeforeEach(func() {
		httpClient = &http.Client{}
		launcher = TestLauncher{}
		uaaServer = NewServer()
		config = cliConfig.NewConfigWithServerURL(uaaServer.URL())
		logger = NewLogger(GinkgoWriter, GinkgoWriter, GinkgoWriter, GinkgoWriter)
	})

	Describe("NewAuthcodeClientImpersonator", func() {
		BeforeEach(func() {
			launcher := TestLauncher{}
			impersonator = NewAuthcodeClientImpersonator(config, "authcodeClientId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)
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
				Expect(impersonator.AuthCallbackServer.Log()).NotTo(Equal(Logger{}))
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
			contentTypeJson := http.Header{}
			contentTypeJson.Add("Content-Type", "application/json")

			uaaServer.RouteToHandler("POST", "/oauth/token", CombineHandlers(
				RespondWith(http.StatusOK, `{
				  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
				  "token_type" : "bearer",
				  "expires_in" : 3000,
				  "scope" : "openid",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`, contentTypeJson),
				VerifyRequest("POST", "/oauth/token"),
				VerifyHeader(http.Header{"Authorization": []string{"Basic YXV0aGNvZGVJZDphdXRoY29kZXNlY3JldA=="}}),
				VerifyBody([]byte(`code=secretcode&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=token&token_format=jwt`)),
			),
			)

			impersonator = NewAuthcodeClientImpersonator(config, "authcodeId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)

			// Start the callback server
			go impersonator.Start()

			// Hit the callback server with an authcode
			Eventually(func() (*http.Response, error) {
				return httpClient.Get("http://localhost:8080/?code=secretcode")
			}, AuthCallbackTimeout, AuthCallbackPollInterval).Should(gstruct.PointTo(gstruct.MatchFields(
				gstruct.IgnoreExtras, gstruct.Fields{
					"StatusCode": Equal(200),
					"Body":       Not(BeNil()),
				},
			)))

			// The callback server should have exchanged the code for a token
			tokenResponse := <-impersonator.Done()
			Expect(tokenResponse.AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			Expect(tokenResponse.TokenType).To(Equal("bearer"))
			Expect(tokenResponse.Extra("scope")).To(Equal("openid"))
			Expect(tokenResponse.Extra("jti")).To(Equal("bc4885d950854fed9a938e96b13ca519"))
			Expect(tokenResponse.Expiry).Should(BeTemporally("~", time.Now(), 3000*time.Second))
		})
	})

	Describe("#Authorize", func() {
		It("launches a browser to the authorize page", func() {
			impersonator = NewAuthcodeClientImpersonator(config, "authcodeId", "authcodesecret", "jwt", "openid", 8080, logger, launcher.Run)

			impersonator.Authorize()

			Expect(launcher.TargetUrl).To(Equal(uaaServer.URL() + "/oauth/authorize?client_id=authcodeId&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=code&scope=openid"))
		})
	})
})
