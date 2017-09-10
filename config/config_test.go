// +build !windows

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
)

var _ = Describe("Config", func() {
	var cfg uaa.Config

	It("can read saved data", func() {
		cfg = uaa.NewConfig()
		target := uaa.NewTarget()
		target.BaseUrl = "http://nowhere.com"
		target.SkipSSLValidation = true

		ctx := uaa.NewContextWithToken("foo-token")
		ctx.ClientId = "cid"
		ctx.Username = "woodstock"
		ctx.GrantType = "client_credentials"

		cfg.AddTarget(target)
		cfg.AddContext(ctx)

		config.WriteConfig(cfg)

		Expect(cfg.ActiveTargetName).To(Equal("url:http://nowhere.com"))
		Expect(cfg.GetActiveContext().AccessToken).To(Equal("foo-token"))
		cfg2 := config.ReadConfig()
		Expect(cfg2.ActiveTargetName).To(Equal("url:http://nowhere.com"))
		Expect(cfg2.GetActiveContext().AccessToken).To(Equal("foo-token"))
	})

	It("can accept a context without previously setting target", func() {
		cfg = uaa.NewConfig()
		ctx := uaa.NewContextWithToken("foo-token")
		ctx.ClientId = "cid"
		ctx.Username = "woodstock"
		ctx.GrantType = "client_credentials"
		cfg.AddContext(ctx)

		config.WriteConfig(cfg)

		Expect(cfg.GetActiveContext().TokenResponse.AccessToken).To(Equal("foo-token"))

		cfg2 := config.ReadConfig()
		Expect(cfg2.ActiveTargetName).To(Equal("url:"))
		Expect(cfg2.GetActiveContext().TokenResponse.AccessToken).To(Equal("foo-token"))
	})

	It("places the config file in .uaa in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.uaa/config.json`))
	})
})