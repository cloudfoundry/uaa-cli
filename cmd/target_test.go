package cmd_test

import (
	"net/http"

	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

const InfoResponseJson = `{
  "app": {
	"version": "4.5.0"
  },
  "links": {
	"uaa": "https://uaa.run.pivotal.io",
	"passwd": "https://account.run.pivotal.io/forgot-password",
	"login": "https://login.run.pivotal.io",
	"register": "https://account.run.pivotal.io/sign-up"
  },
  "zone_name": "uaa",
  "entityID": "login.run.pivotal.io",
  "commit_id": "df80f63",
  "idpDefinitions": {},
  "prompts": {
	"username": [
	  "text",
	  "Email"
	],
	"password": [
	  "password",
	  "Password"
	]
  },
  "timestamp": "2017-07-21T22:45:01+0000"
}`

var _ = Describe("Target", func() {
	Describe("when no new url is provided", func() {
		Describe("and a target was previously set", func() {
			BeforeEach(func() {
				c := config.NewConfigWithServerURL(server.URL())
				config.WriteConfig(c)
			})

			It("shows the currently set target and UAA version", func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusOK, InfoResponseJson, contentTypeJson),
				)

				session := runCommand("target")

				Eventually(session).Should(Exit(0))
				expectedJson := `{ "Target": "` + server.URL() + `", "Status": "OK", "UaaVersion": "4.5.0", "SkipSSLValidation": false }`
				Eventually(session.Out.Contents()).Should(MatchJSON(expectedJson))
			})

			It("shows <unknown version> when UAA can't be reached", func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusBadRequest, ""),
				)

				session := runCommand("target")

				Eventually(session.Out).Should(Say(`"Status": "ERROR"`))
				Eventually(session.Out).Should(Say(`"UaaVersion": "unknown"`))
				Eventually(session).Should(Exit(1))
			})
		})

		Describe("and a target was never set", func() {
			It("displays empty target", func() {
				session := runCommand("target")

				Eventually(session).Should(Exit(0))
				expectedJson := `{ "Target": "", "Status": "", "UaaVersion": "", "SkipSSLValidation": false }`
				Eventually(session.Out.Contents()).Should(MatchJSON(expectedJson))
			})
		})
	})

	Describe("when a new url is provided", func() {
		Describe("when the url is a valid UAA", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusOK, InfoResponseJson, contentTypeJson),
				)

				c := config.NewConfig()
				config.WriteConfig(c)
			})

			It("updates the saved context", func() {
				Expect(config.ReadConfig().GetActiveTarget().BaseUrl).To(Equal(""))

				runCommand("target", server.URL())

				Expect(config.ReadConfig().GetActiveTarget().BaseUrl).NotTo(Equal(""))
				Expect(config.ReadConfig().GetActiveTarget().BaseUrl).To(Equal(server.URL()))
			})

			It("displays a success message", func() {
				session := runCommand("target", server.URL())

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Target set to " + server.URL()))
			})

			It("respects the --skip-ssl-validation flag", func() {
				runCommand("target", server.URL())
				Expect(config.ReadConfig().GetActiveTarget().SkipSSLValidation).To(BeFalse())

				runCommand("target", server.URL(), "--skip-ssl-validation")
				Expect(config.ReadConfig().GetActiveTarget().SkipSSLValidation).To(BeTrue())
			})
		})

		Describe("when the UAA cannot be reached", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusNotFound, ""),
				)

				c := config.NewConfigWithServerURL("http://someuaa.com")
				config.WriteConfig(c)
			})

			It("does not update the saved context", func() {
				runCommand("target", server.URL())

				Expect(config.ReadConfig().GetActiveTarget().BaseUrl).To(Equal("http://someuaa.com"))
			})

			It("displays an error message", func() {
				session := runCommand("target", server.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The target " + server.URL() + " could not be set."))
			})
		})
	})
})
