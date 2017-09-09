package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"io/ioutil"
)

type TestLauncher struct {
	TargetUrl string
}

func (tl *TestLauncher) Run(target string) error {
	tl.TargetUrl = target
	return nil
}

var _ = Describe("ImplicitClientImpersonator", func() {
	var (
		impersonator ImplicitClientImpersonator
		logger       utils.Logger
	)

	BeforeEach(func() {
		logger = utils.NewLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	})

	Describe("NewImplicitClientImpersonator", func() {
		BeforeEach(func() {
			launcher := TestLauncher{}
			impersonator = NewImplicitClientImpersonator("implicitId", "http://uaa.com", "jwt", "openid", 8080, logger, launcher.Run)
		})

		Describe("configures an AuthCallbackListener", func() {
			It("with appropriate static content", func() {
				Expect(impersonator.AuthCallbackServer.Javascript()).To(ContainSubstring("XMLHttpRequest"))
				Expect(impersonator.AuthCallbackServer.CSS()).To(ContainSubstring("Source Sans Pro"))
				Expect(impersonator.AuthCallbackServer.Html()).To(ContainSubstring("Implicit Grant: Success"))
			})

			It("with the desired port", func() {
				Expect(impersonator.AuthCallbackServer.Port()).To(Equal(8080))
			})

			It("with its logger", func() {
				Expect(impersonator.AuthCallbackServer.Log()).NotTo(Equal(utils.Logger{}))
				Expect(impersonator.AuthCallbackServer.Log()).To(Equal(logger))
			})

			It("with hangup func that looks for access_token in query params", func() {
				done := make(chan string)

				urlParams := url.Values{}
				urlParams.Add("access_token", "56575db17b164e568668c0085ed14ae1")
				go impersonator.AuthCallbackServer.Hangup(done, urlParams)

				Expect(<-done).To(Equal("56575db17b164e568668c0085ed14ae1"))
			})
		})
	})

	Describe("#Start", func() {
		BeforeEach(func() {
			launcher := TestLauncher{}
			impersonator = NewImplicitClientImpersonator("implicitId", "http://uaa.com", "jwt", "openid", 8080, logger, launcher.Run)
			impersonator.AuthCallbackServer = FakeCallbackServer{}
		})

		It("starts the AuthCallbackServer", func() {
			go impersonator.Start()
			Expect(<-impersonator.Done()).To(Equal("server was started"))
		})
	})

	Describe("#Authorize", func() {
		It("launches a browser to the authorize page", func() {
			launcher := TestLauncher{}
			impersonator = NewImplicitClientImpersonator("implicitId", "http://uaa.com", "jwt", "openid", 8080, logger, launcher.Run)

			impersonator.Authorize()

			Expect(launcher.TargetUrl).To(Equal("http://uaa.com/oauth/authorize?client_id=implicitId&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=token&scope=openid&token_format=jwt"))
		})
	})
})
