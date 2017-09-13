package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("GetAuthcodeToken", func() {
	var (
		c          uaa.Config
		ctx        uaa.UaaContext
		logger     utils.Logger
		launcher   TestLauncher
		httpClient *http.Client
	)

	BeforeEach(func() {
		c = uaa.NewConfigWithServerURL(server.URL())
		config.WriteConfig(c)
		ctx = c.GetActiveContext()
		launcher = TestLauncher{}
		logger = utils.NewLogger(GinkgoWriter, GinkgoWriter, GinkgoWriter, GinkgoWriter)
		httpClient = &http.Client{}
	})

	It("updates the saved context with the user's access token", func() {
		server.RouteToHandler("POST", "/oauth/token", CombineHandlers(
			VerifyRequest("POST", "/oauth/token"),
			RespondWith(http.StatusOK, `{
				  "access_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
				  "refresh_token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gSADFJSKADJFLsdfandydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
				  "token_type" : "bearer",
				  "expires_in" : 3000,
				  "scope" : "openid",
				  "jti" : "bc4885d950854fed9a938e96b13ca519"
				}`),
			VerifyFormKV("code", "ASDFGHJKL"),
			VerifyFormKV("client_id", "shinyclient"),
			VerifyFormKV("client_secret", "shinysecret"),
			VerifyFormKV("grant_type", "authorization_code"),
			VerifyFormKV("token_format", "jwt"),
			VerifyFormKV("response_type", "token"),
			VerifyFormKV("redirect_uri", "http://localhost:8080")),
		)

		doneRunning := make(chan bool)

		imp := cli.NewAuthcodeClientImpersonator(httpClient, c, "shinyclient", "shinysecret", "jwt", "openid", 8080, logger, launcher.Run)
		go AuthcodeTokenCommandRun(doneRunning, "shinyclient", imp, &logger)

		// UAA sends the user to this redirect_uri after they auth and grant approvals
		httpClient.Get("http://localhost:8080/?code=ASDFGHJKL")

		<-doneRunning
		Expect(launcher.Target).To(Equal(server.URL() + "/oauth/authorize?client_id=shinyclient&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=code"))
		Expect(GetSavedConfig().GetActiveContext().AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		Expect(GetSavedConfig().GetActiveContext().RefreshToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gSADFJSKADJFLsdfandydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		Expect(GetSavedConfig().GetActiveContext().ClientId).To(Equal("shinyclient"))
		Expect(GetSavedConfig().GetActiveContext().GrantType).To(Equal(uaa.GrantType("authorization_code")))
		Expect(GetSavedConfig().GetActiveContext().TokenType).To(Equal("bearer"))
		Expect(GetSavedConfig().GetActiveContext().Scope).To(Equal("openid"))
	})

	Describe("Validations", func() {
		It("requires a client id", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := AuthcodeTokenArgumentValidation(cfg, []string{}, "secret", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_id` must be specified."))
		})

		It("requires a client secret", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_secret` must be specified."))
		})

		It("requires a port", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "secret", "jwt", 0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `port` must be specified."))
		})

		It("rejects invalid token formats", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "secret", "bogus-format", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`The token format "bogus-format" is unknown. Available formats: [jwt, opaque]`))
		})

		It("requires a target to have been set", func() {
			err := AuthcodeTokenArgumentValidation(uaa.NewConfig(), []string{"clientid"}, "secret", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(MISSING_TARGET))
		})
	})
})
