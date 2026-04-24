package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Targets", func() {
	Describe("when no targets have been set", func() {
		BeforeEach(func() {
			c := config.NewConfig()
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and indicates no targets", func() {
			session := runCommand("targets")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say("No targets set."))
		})
	})

	Describe("when one target has been set", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0, lists the target URL, and marks it as active", func() {
			session := runCommand("targets")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(`\*`))
			Expect(session.Out).To(Say(server.URL()))
		})

		It("exits 0 with --verbose flag", func() {
			session := runCommand("targets", "--verbose")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(server.URL()))
		})
	})

	Describe("when multiple targets have been set", func() {
		var secondURL string

		BeforeEach(func() {
			secondURL = "http://second-uaa.example.com"

			c := config.NewConfigWithServerURL(server.URL())
			t2 := config.NewTarget()
			t2.BaseUrl = secondURL
			c.AddTarget(t2)
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and lists all targets with active marker on the second", func() {
			session := runCommand("targets")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(server.URL()))
			Expect(session.Out).To(Say(secondURL))
		})

		It("marks only the active target with *", func() {
			session := runCommand("targets")

			Eventually(session).Should(Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("* "))
			Expect(output).To(ContainSubstring(secondURL))
		})
	})
})
