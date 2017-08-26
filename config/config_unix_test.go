// +build !windows

package config_test

import (
	"os"

	"path"

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

	It("set appropriate permissions for persisted files", func() {
		config.WriteConfig(cfg)

		parentStat, _ := os.Stat(path.Dir(config.ConfigPath()))
		Expect(parentStat.Mode().String()).To(Equal("drwxr-xr-x"))

		fileStat, _ := os.Stat(config.ConfigPath())
		Expect(fileStat.Mode().String()).To(Equal("-rw-------"))
	})
})
