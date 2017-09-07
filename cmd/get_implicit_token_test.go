package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	var c uaa.Config
	var ctx uaa.UaaContext
	var logger utils.Logger

	BeforeEach(func() {
		c = uaa.NewConfigWithServerURL(server.URL())
		config.WriteConfig(c)
		ctx = c.GetActiveContext()
		logger := GetLogger()
		logger.Mute()
	})

	AfterEach(func() {
		logger.Unmute()
	})

	It("launches a browser for the authorize page and gets the callback params", func() {
		launcher := TestLauncher{}
		doneRunning := make(chan bool)
		go ImplicitTokenCommandRun(doneRunning, launcher.Run, "openid", "shinyclient", 8080)

		httpClient := &http.Client{}
		// UAA sends the user to this redirect_uri after they auth and grant approvals
		httpClient.Get("http://localhost:8080/?access_token=foo")

		<-doneRunning
		Expect(launcher.Target).To(Equal(server.URL() + "/oauth/authorize?client_id=shinyclient&redirect_uri=http%3A%2F%2Flocalhost%3A8080&response_type=token&scope=openid&token_format=jwt"))
		Expect(GetSavedConfig().GetActiveContext().AccessToken).To(Equal("foo"))
		Expect(GetSavedConfig().GetActiveContext().ClientId).To(Equal("shinyclient"))
		Expect(GetSavedConfig().GetActiveContext().GrantType).To(Equal(uaa.GrantType("implicit")))
		Expect(GetSavedConfig().GetActiveContext().TokenType).To(Equal(uaa.TokenFormat("jwt")))
		Expect(GetSavedConfig().GetActiveContext().Scope).To(Equal("openid"))
	})

	It("handles multiple scopes", func() {
		launcher := TestLauncher{}
		doneRunning := make(chan bool)
		go ImplicitTokenCommandRun(doneRunning, launcher.Run, "openid,user_attributes", "shinyclient", 8081)

		httpClient := &http.Client{}
		// UAA sends the user to this redirect_uri after they auth and grant approvals
		httpClient.Get("http://localhost:8081/?access_token=foo")

		<-doneRunning
		Expect(launcher.Target).To(Equal(server.URL() + "/oauth/authorize?client_id=shinyclient&redirect_uri=http%3A%2F%2Flocalhost%3A8081&response_type=token&scope=openid%2Cuser_attributes&token_format=jwt"))
		Expect(GetSavedConfig().GetActiveContext().AccessToken).To(Equal("foo"))
		Expect(GetSavedConfig().GetActiveContext().ClientId).To(Equal("shinyclient"))
		Expect(GetSavedConfig().GetActiveContext().GrantType).To(Equal(uaa.GrantType("implicit")))
	})
})
