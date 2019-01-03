package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"net/http"
)

type TestLauncher struct {
	Target string
}

func (tl *TestLauncher) Run(target string) error {
	tl.Target = target
	return nil
}

var _ = Describe("GetImplicitToken", func() {
	var c config.Config
	var logger cli.Logger

	BeforeEach(func() {
		c = config.NewConfigWithServerURL(server.URL())
		config.WriteConfig(c)
		logger = cli.NewLogger(GinkgoWriter, GinkgoWriter, GinkgoWriter, GinkgoWriter)
	})

	It("launches a browser for the authorize page and gets the callback params", func() {
		launcher := TestLauncher{}
		doneRunning := make(chan bool)

		imp := cli.NewImplicitClientImpersonator("shinyclient", server.URL(), "jwt", "openid", 8080, logger, launcher.Run)
		go ImplicitTokenCommandRun(doneRunning, "shinyclient", imp, &logger)

		httpClient := &http.Client{}
		// UAA sends the user to this redirect_uri after they auth and grant approvals
		Eventually(func() (*http.Response, error) {
			return httpClient.Get("http://localhost:8080/?access_token=foo&scope=openid&token_type=bearer")
		}, AuthCallbackTimeout, AuthCallbackPollInterval).Should(gstruct.PointTo(gstruct.MatchFields(
			gstruct.IgnoreExtras, gstruct.Fields{
				"StatusCode": Equal(200),
				"Body":       Not(BeNil()),
			},
		)))

		Eventually(doneRunning, AuthCallbackTimeout, AuthCallbackPollInterval).Should(Receive())

		Expect(launcher.Target).To(Equal(server.URL() + "/oauth/authorize?client_id=shinyclient&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=token&scope=openid&token_format=jwt"))
		Expect(GetSavedConfig().GetActiveContext().Token.AccessToken).To(Equal("foo"))
		Expect(GetSavedConfig().GetActiveContext().Token.TokenType).To(Equal("bearer"))
	})
})
