package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("GetTokenKeys", func() {
	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		It("shows an array of keys", func() {
			keyList := `{ "keys": [{
			  "kty" : "RSA",
			  "e" : "AQAB",
			  "use" : "sig",
			  "kid" : "testKey",
			  "alg" : "RS256",
			  "value" : "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0m59l2u9iDnMbrXHfqkO\nrn2dVQ3vfBJqcDuFUK03d+1PZGbVlNCqnkpIJ8syFppW8ljnWweP7+LiWpRoz0I7\nfYb3d8TjhV86Y997Fl4DBrxgM6KTJOuE/uxnoDhZQ14LgOU2ckXjOzOdTsnGMKQB\nLCl0vpcXBtFLMaSbpv1ozi8h7DJyVZ6EnFQZUWGdgTMhDrmqevfx95U/16c5WBDO\nkqwIn7Glry9n9Suxygbf8g5AzpWcusZgDLIIZ7JTUldBb8qU2a0Dl4mvLZOn4wPo\njfj9Cw2QICsc5+Pwf21fP+hzf+1WSRHbnYv8uanRO0gZ8ekGaghM/2H6gqJbo2nI\nJwIDAQAB\n-----END PUBLIC KEY-----",
			  "n" : "ANJufZdrvYg5zG61x36pDq59nVUN73wSanA7hVCtN3ftT2Rm1ZTQqp5KSCfLMhaaVvJY51sHj-_i4lqUaM9CO32G93fE44VfOmPfexZeAwa8YDOikyTrhP7sZ6A4WUNeC4DlNnJF4zsznU7JxjCkASwpdL6XFwbRSzGkm6b9aM4vIewyclWehJxUGVFhnYEzIQ65qnr38feVP9enOVgQzpKsCJ-xpa8vZ_UrscoG3_IOQM6VnLrGYAyyCGeyU1JXQW_KlNmtA5eJry2Tp-MD6I34_QsNkCArHOfj8H9tXz_oc3_tVkkR252L_Lmp0TtIGfHpBmoITP9h-oKiW6NpyCc"
			},
			{
			  "kty" : "RSA",
			  "e" : "AQAB",
			  "use" : "sig",
			  "kid" : "testKey2",
			  "alg" : "RS256",
			  "value" : "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0m59l2u9iDnMbrXHfqkO\nrn2dVQ3vfBJqcDuFUK03d+1PZGbVlNCqnkpIJ8syFppW8ljnWweP7+LiWpRoz0I7\nfYb3d8TjhV86Y997Fl4DBrxgM6KTJOuE/uxnoDhZQ14LgOU2ckXjOzOdTsnGMKQB\nLCl0vpcXBtFLMaSbpv1ozi8h7DJyVZ6EnFQZUWGdgTMhDrmqevfx95U/16c5WBDO\nkqwIn7Glry9n9Suxygbf8g5AzpWcusZgDLIIZ7JTUldBb8qU2a0Dl4mvLZOn4wPo\njfj9Cw2QICsc5+Pwf21fP+hzf+1WSRHbnYv8uanRO0gZ8ekGaghM/2H6gqJbo2nI\nJwIDAQAB\n-----END PUBLIC KEY-----",
			  "n" : "ANJufZdrvYg5zG61x36pDq59nVUN73wSanA7hVCtN3ftT2Rm1ZTQqp5KSCfLMhaaVvJY51sHj-_i4lqUaM9CO32G93fE44VfOmPfexZeAwa8YDOikyTrhP7sZ6A4WUNeC4DlNnJF4zsznU7JxjCkASwpdL6XFwbRSzGkm6b9aM4vIewyclWehJxUGVFhnYEzIQ65qnr38feVP9enOVgQzpKsCJ-xpa8vZ_UrscoG3_IOQM6VnLrGYAyyCGeyU1JXQW_KlNmtA5eJry2Tp-MD6I34_QsNkCArHOfj8H9tXz_oc3_tVkkR252L_Lmp0TtIGfHpBmoITP9h-oKiW6NpyCc"
			}]}`

			server.RouteToHandler("GET", "/token_keys",
				CombineHandlers(
					RespondWith(http.StatusOK, keyList, contentTypeJson),
				),
			)

			session := runCommand("get-token-keys")

			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(`[{
			  "kty" : "RSA",
			  "e" : "AQAB",
			  "use" : "sig",
			  "kid" : "testKey",
			  "alg" : "RS256",
			  "value" : "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0m59l2u9iDnMbrXHfqkO\nrn2dVQ3vfBJqcDuFUK03d+1PZGbVlNCqnkpIJ8syFppW8ljnWweP7+LiWpRoz0I7\nfYb3d8TjhV86Y997Fl4DBrxgM6KTJOuE/uxnoDhZQ14LgOU2ckXjOzOdTsnGMKQB\nLCl0vpcXBtFLMaSbpv1ozi8h7DJyVZ6EnFQZUWGdgTMhDrmqevfx95U/16c5WBDO\nkqwIn7Glry9n9Suxygbf8g5AzpWcusZgDLIIZ7JTUldBb8qU2a0Dl4mvLZOn4wPo\njfj9Cw2QICsc5+Pwf21fP+hzf+1WSRHbnYv8uanRO0gZ8ekGaghM/2H6gqJbo2nI\nJwIDAQAB\n-----END PUBLIC KEY-----",
			  "n" : "ANJufZdrvYg5zG61x36pDq59nVUN73wSanA7hVCtN3ftT2Rm1ZTQqp5KSCfLMhaaVvJY51sHj-_i4lqUaM9CO32G93fE44VfOmPfexZeAwa8YDOikyTrhP7sZ6A4WUNeC4DlNnJF4zsznU7JxjCkASwpdL6XFwbRSzGkm6b9aM4vIewyclWehJxUGVFhnYEzIQ65qnr38feVP9enOVgQzpKsCJ-xpa8vZ_UrscoG3_IOQM6VnLrGYAyyCGeyU1JXQW_KlNmtA5eJry2Tp-MD6I34_QsNkCArHOfj8H9tXz_oc3_tVkkR252L_Lmp0TtIGfHpBmoITP9h-oKiW6NpyCc"
			},
			{
			  "kty" : "RSA",
			  "e" : "AQAB",
			  "use" : "sig",
			  "kid" : "testKey2",
			  "alg" : "RS256",
			  "value" : "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0m59l2u9iDnMbrXHfqkO\nrn2dVQ3vfBJqcDuFUK03d+1PZGbVlNCqnkpIJ8syFppW8ljnWweP7+LiWpRoz0I7\nfYb3d8TjhV86Y997Fl4DBrxgM6KTJOuE/uxnoDhZQ14LgOU2ckXjOzOdTsnGMKQB\nLCl0vpcXBtFLMaSbpv1ozi8h7DJyVZ6EnFQZUWGdgTMhDrmqevfx95U/16c5WBDO\nkqwIn7Glry9n9Suxygbf8g5AzpWcusZgDLIIZ7JTUldBb8qU2a0Dl4mvLZOn4wPo\njfj9Cw2QICsc5+Pwf21fP+hzf+1WSRHbnYv8uanRO0gZ8ekGaghM/2H6gqJbo2nI\nJwIDAQAB\n-----END PUBLIC KEY-----",
			  "n" : "ANJufZdrvYg5zG61x36pDq59nVUN73wSanA7hVCtN3ftT2Rm1ZTQqp5KSCfLMhaaVvJY51sHj-_i4lqUaM9CO32G93fE44VfOmPfexZeAwa8YDOikyTrhP7sZ6A4WUNeC4DlNnJF4zsznU7JxjCkASwpdL6XFwbRSzGkm6b9aM4vIewyclWehJxUGVFhnYEzIQ65qnr38feVP9enOVgQzpKsCJ-xpa8vZ_UrscoG3_IOQM6VnLrGYAyyCGeyU1JXQW_KlNmtA5eJry2Tp-MD6I34_QsNkCArHOfj8H9tXz_oc3_tVkkR252L_Lmp0TtIGfHpBmoITP9h-oKiW6NpyCc"
			}]`))
			Expect(session).Should(Exit(0))
		})

		It("shows symmetric keys with unused JWK fields omitted", func() {
			symmetricKeyList := `{"keys": [{
			  "kty" : "MAC",
			  "alg" : "HS256",
			  "value" : "key",
			  "use" : "sig",
			  "kid" : "testKey"
			}]}`

			server.RouteToHandler("GET", "/token_keys",
				CombineHandlers(
					RespondWith(http.StatusOK, symmetricKeyList, contentTypeJson),
				),
			)

			session := runCommand("get-token-keys")

			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(`[{
			  "kty" : "MAC",
			  "alg" : "HS256",
			  "value" : "key",
			  "use" : "sig",
			  "kid" : "testKey"
			}]`))
			Expect(session).Should(Exit(0))
		})
	})

	Describe("Validations", func() {
		It("it requires a target to have been set", func() {
			config.WriteConfig(config.NewConfig())

			session := runCommand("get-token-keys")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cli.MISSING_TARGET))
		})
	})
})
