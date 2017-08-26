// +build !windows

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
)

var _ = Describe("Config", func() {
	var cfg uaa.Config

	BeforeEach(func() {
		cfg = uaa.Config{}
		cfg.Context = uaa.UaaContext{
			BaseUrl: "https://login.example.com",
		}
	})

	It("places the config file in .uaa in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.uaa/config.json`))
	})
})
