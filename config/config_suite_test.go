package config_test

import (
	"github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"code.cloudfoundry.org/uaa-cli/config"
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

func NewContextWithToken(accessToken string) uaa.AuthContext {
	ctx := uaa.AuthContext{}
	ctx.AccessToken = accessToken
	return ctx
}
