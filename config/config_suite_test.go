package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"testing"
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

func NewContextWithToken(accessToken string) uaa.UaaContext {
	ctx := uaa.UaaContext{}
	ctx.AccessToken = accessToken
	return ctx
}
