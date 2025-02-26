package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"code.cloudfoundry.org/uaa-cli/config"
	"golang.org/x/oauth2"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = BeforeEach(func() {
	config.RemoveConfig()
})

var _ = AfterEach(func() {
	config.RemoveConfig()
})

func NewContextWithToken(accessToken string) config.UaaContext {
	ctx := config.UaaContext{
		Token: oauth2.Token{
			AccessToken: accessToken,
		},
	}
	return ctx
}
