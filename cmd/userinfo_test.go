package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
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
		  "phone_number": "",
		  "previous_logon_time": 1503123277743,
		  "name": "Charlie Brown"
		}`

		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			c.AddContext(config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"))
			config.WriteConfig(c)
		})

		It("shows the info response", func() {
			server.RouteToHandler("GET", "/userinfo", CombineHandlers(
				RespondWith(200, userinfoJson),
				VerifyRequest("GET", "/userinfo", "scheme=openid"),
				VerifyHeaderKV("Accept", "application/json"),
				VerifyHeaderKV("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"),
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
			Expect(session.Err).To(Say("An error occurred while calling " + server.URL() + "/userinfo"))
		})
	})

	Describe("Validations", func() {
		It("requires a target", func() {
			err := cmd.UserinfoValidations(config.Config{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := cmd.UserinfoValidations(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})
	})
})
