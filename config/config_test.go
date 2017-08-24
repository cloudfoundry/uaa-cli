// +build !windows

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jhamon/uaa-cli/config"
)

var _ = Describe("Config", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			TargetUrl:  "https://login.example.com",
		}
	})

	It("places the config file in .uaa in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.uaa/config.json`))
	})
})
