// +build !windows

package config_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Config", func() {
	var cfg config.Config

	It("can read saved data", func() {
		cfg = config.NewConfig()

		target := config.NewTarget()
		target.BaseUrl = "http://nowhere.com"
		target.SkipSSLValidation = true

		ctx := NewContextWithToken("foo-token")
		ctx.ClientId = "cid"
		ctx.Username = "woodstock"
		ctx.GrantType = "client_credentials"

		cfg.AddTarget(target)
		cfg.AddContext(ctx)

		config.WriteConfig(cfg)

		Expect(cfg.ActiveTargetName).To(Equal("url:http://nowhere.com"))
		Expect(cfg.GetActiveContext().Token.AccessToken).To(Equal("foo-token"))
		cfg2 := config.ReadConfig()
		Expect(cfg2.ActiveTargetName).To(Equal("url:http://nowhere.com"))
		Expect(cfg2.GetActiveContext().Token.AccessToken).To(Equal("foo-token"))
	})

	It("can accept a context without previously setting target", func() {
		cfg = config.NewConfig()
		ctx := NewContextWithToken("foo-token")
		ctx.ClientId = "cid"
		ctx.Username = "woodstock"
		ctx.GrantType = "client_credentials"
		cfg.AddContext(ctx)

		config.WriteConfig(cfg)

		Expect(cfg.GetActiveContext().Token.AccessToken).To(Equal("foo-token"))

		cfg2 := config.ReadConfig()
		Expect(cfg2.ActiveTargetName).To(Equal("url:"))
		Expect(cfg2.GetActiveContext().Token.AccessToken).To(Equal("foo-token"))
	})

	Context("when UAA_HOME env var is not set", func() {
		It("places the config file in .uaa in the home directory", func() {
			homeDir := os.Getenv("HOME")
			Expect(config.ConfigPath()).To(HavePrefix(homeDir))
			Expect(config.ConfigPath()).To(HaveSuffix(`/.uaa/config.json`))
		})
	})

	Context("when UAA_HOME env var is set", func() {
		var uaaHome string

		BeforeEach(func() {
			var err error
			uaaHome, err = ioutil.TempDir(os.TempDir(), "uaa-home")
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("UAA_HOME", uaaHome)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			if uaaHome != "" {
				if _, err := os.Stat(uaaHome); !os.IsNotExist(err) {
					err := os.RemoveAll(uaaHome)
					Expect(err).NotTo(HaveOccurred())
				}
			}

			err := os.Unsetenv("UAA_HOME")
			Expect(err).ToNot(HaveOccurred())
		})

		It("places the config file in the directory pointed to by UAA_HOME", func() {
			Expect(config.ConfigPath()).To(HavePrefix(uaaHome))
			Expect(config.ConfigPath()).To(HaveSuffix(`config.json`))
		})
	})
})
