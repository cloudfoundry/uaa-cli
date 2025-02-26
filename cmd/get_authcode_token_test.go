package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("GetAuthcodeToken", func() {
	var (
		c          config.Config
		logger     cli.Logger
		launcher   TestLauncher
		httpClient *http.Client
	)

	BeforeEach(func() {
		c = config.NewConfigWithServerURL(server.URL())
		config.WriteConfig(c)
		launcher = TestLauncher{}
		logger = cli.NewLogger(GinkgoWriter, GinkgoWriter, GinkgoWriter, GinkgoWriter)
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
				}`, contentTypeJson),
			VerifyFormKV("code", "ASDFGHJKL"),
			VerifyHeaderKV("Authorization", "Basic c2hpbnljbGllbnQ6c2hpbnlzZWNyZXQ="),
			VerifyFormKV("grant_type", "authorization_code"),
			VerifyFormKV("token_format", "jwt")),
		)

		doneRunning := make(chan bool)

		imp := cli.NewAuthcodeClientImpersonator(c, "shinyclient", "shinysecret", "jwt", "openid", 9090, logger, launcher.Run)
		go AuthcodeTokenCommandRun(doneRunning, "shinyclient", imp, &logger)

		// UAA sends the user to this redirect_uri after they auth and grant approvals
		Eventually(func() (*http.Response, error) {
			return httpClient.Get("http://localhost:9090/?code=ASDFGHJKL")
		}, AuthCallbackTimeout, AuthCallbackPollInterval).Should(gstruct.PointTo(gstruct.MatchFields(
			gstruct.IgnoreExtras, gstruct.Fields{
				"StatusCode": Equal(200),
				"Body":       Not(BeNil()),
			},
		)))

		Eventually(doneRunning, AuthCallbackTimeout, AuthCallbackPollInterval).Should(Receive())

		Expect(launcher.Target).To(Equal(server.URL() + "/oauth/authorize?client_id=shinyclient&redirect_uri=http%3A%2F%2Flocalhost%3A9090&response_type=code&scope=openid"))
		Expect(GetSavedConfig().GetActiveContext().Token.AccessToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		Expect(GetSavedConfig().GetActiveContext().Token.RefreshToken).To(Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gSADFJSKADJFLsdfandydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
		Expect(GetSavedConfig().GetActiveContext().Token.TokenType).To(Equal("bearer"))
	})

	Describe("Validations", func() {
		It("requires a client id", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := AuthcodeTokenArgumentValidation(cfg, []string{}, "secret", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_id` must be specified."))
		})

		It("requires a client secret", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `client_secret` must be specified."))
		})

		It("requires a port", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "secret", "jwt", 0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing argument `port` must be specified."))
		})

		It("rejects invalid token formats", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := AuthcodeTokenArgumentValidation(cfg, []string{"clientid"}, "secret", "bogus-format", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`The token format "bogus-format" is unknown. Available formats: [jwt, opaque]`))
		})

		It("requires a target to have been set", func() {
			err := AuthcodeTokenArgumentValidation(config.NewConfig(), []string{"clientid"}, "secret", "jwt", 8001)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(cli.MISSING_TARGET))
		})
	})
})
