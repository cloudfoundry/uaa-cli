package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Info", func() {
	Describe("and a target was previously set", func() {
		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			config.WriteConfig(c)
		})

		ItSupportsTheVerboseFlagWhenGet("info", "/info", InfoResponseJson)

		It("shows the info response", func() {
			server.RouteToHandler("GET", "/info",
				RespondWith(http.StatusOK, InfoResponseJson),
			)

			session := runCommand("info")

			Eventually(session).Should(Exit(0))
			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(InfoResponseJson))
		})

		It("handles request errors", func() {
			server.RouteToHandler("GET", "/info",
				RespondWith(http.StatusBadRequest, ""),
			)

			session := runCommand("info")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("An unknown error occurred while calling " + server.URL() + "/info"))
		})
	})

	Describe("when no target was previously set", func() {
		BeforeEach(func() {
			c := uaa.Config{}
			config.WriteConfig(c)
		})

		It("tells the user to set a target", func() {
			session := runCommand("info")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("You must set a target in order to use this command."))
		})
	})
})
