package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
	"github.com/jhamon/uaa-cli/config"
	"github.com/jhamon/uaa-cli/uaa"
)

const InfoResponseJson string = `{
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
				c := config.Config{}
				c.Context = uaa.UaaContext{}
				c.Context.BaseUrl = server.URL()
				config.WriteConfig(c)
			})

			It("shows the currently set target and UAA version", func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusOK, InfoResponseJson),
				)

				session := runCommand("target")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Target: " + server.URL()))
				Eventually(session.Out).Should(Say("Status: OK"))
				Eventually(session.Out).Should(Say("UAA Version: 4.5.0"))
			})

			It("shows <unknown version> when UAA can't be reached", func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusBadRequest, ""),
				)

				session := runCommand("target")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Target: " + server.URL()))
				Eventually(session.Out).Should(Say("Status: ERROR"))
				Eventually(session.Out).Should(Say("UAA Version: unknown"))
			})
		})

		Describe("and a target was never set", func() {
			It("displays empty target", func() {
				session := runCommand("target")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Target:"))
				Eventually(session.Out).Should(Say("Status:"))
				Eventually(session.Out).Should(Say("UAA Version:"))
			})
		})
	})

	Describe("when a new url is provided", func() {
		Describe("when the url is a valid UAA", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusOK, InfoResponseJson),
				)

				c := config.Config{}
				c.Context = uaa.UaaContext{}
				config.WriteConfig(c)
			})

			It("updates the saved context", func() {
				Expect(config.ReadConfig().Context.BaseUrl).To(Equal(""))

				runCommand("target", server.URL())

				Expect(config.ReadConfig().Context.BaseUrl).NotTo(Equal(""))
				Expect(config.ReadConfig().Context.BaseUrl).To(Equal(server.URL()))
			})

			It("displays a success message", func() {
				session := runCommand("target", server.URL())

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Target set to " + server.URL()))
			})
		})

		Describe("when the UAA cannot be reached", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/info",
					RespondWith(http.StatusNotFound, ""),
				)

				c := config.Config{}
				c.Context = uaa.UaaContext{}
				c.Context.BaseUrl = "http://someuaa.com"
				config.WriteConfig(c)
			})

			It("does not update the saved context", func() {
				runCommand("target", server.URL())

				Expect(config.ReadConfig().Context.BaseUrl).To(Equal("http://someuaa.com"))
			})

			It("displays an error message", func() {
				session := runCommand("target", server.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Out).Should(Say("The target " + server.URL() + " is not responding and could not be set."))
			})
		})
	})
})
