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
	"code.cloudfoundry.org/uaa-cli/cmd"
)

var _ = Describe("Userinfo", func() {
	Describe("and a target was previously set", func() {
		userinfoJson := `{
		  "user_id": "d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87",
		  "sub": "d6ef6c2e-02f6-477a-a7c6-18e27f9a6e87",
		  "user_name": "charlieb",
		  "given_name": "Charlie",
		  "family_name": "Brown",
		  "email": "charlieb@peanuts.com",
		  "phone_number": null,
		  "previous_logon_time": 1503123277743,
		  "name": "Charlie Brown"
		}`

		BeforeEach(func() {
			c := uaa.NewConfigWithServerURL(server.URL())
			c.AddContext(uaa.NewContextWithToken("access_token"))
			config.WriteConfig(c)
		})

		It("shows the info response", func() {
			server.RouteToHandler("GET", "/userinfo", CombineHandlers(
				RespondWith(200, userinfoJson),
				VerifyRequest("GET", "/userinfo", "scheme=openid"),
				VerifyHeaderKV("Accept", "application/json"),
				VerifyHeaderKV("Authorization", "bearer access_token"),
			))

			session := runCommand("userinfo")

			Eventually(session).Should(Exit(0))
			outputBytes := session.Out.Contents()
			Expect(outputBytes).To(MatchJSON(userinfoJson))
		})

		It("handles request errors", func() {
			server.RouteToHandler("GET", "/userinfo",
				RespondWith(http.StatusBadRequest, ""),
			)

			session := runCommand("userinfo")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("An unknown error occurred while calling " + server.URL() + "/userinfo"))
		})
	})

	Describe("Validations", func() {
		It("requires a target", func() {
			err := cmd.UserinfoValidations(uaa.Config{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := cmd.UserinfoValidations(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})
	})
})
