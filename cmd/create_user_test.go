package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("CreateUser", func() {
	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
	})

	Describe("Validations", func() {
		It("requires a target to have been set", func() {
			config.WriteConfig(uaa.NewConfig())

			session := runCommand("create-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cmd.MISSING_TARGET))
		})

		It("requires a token in context", func() {
			config.WriteConfig(uaa.NewConfigWithServerURL(server.URL()))

			session := runCommand("create-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cmd.MISSING_CONTEXT))
		})

		It("requires a username", func() {
			session := runCommand("create-user")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("The positional argument USERNAME must be specified."))
		})

		It("requires a family name (last name)", func() {
			session := runCommand("create-user", "woodstock")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("Missing argument `familyName` must be specified."))
		})

		It("requires a given name (first name)", func() {
			session := runCommand("create-user",
				"woodstock",
				"--familyName", "Bird",
			)

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("Missing argument `givenName` must be specified."))
		})

		It("requires an email address", func() {
			session := runCommand("create-user",
				"woodstock",
				"--familyName", "Bird",
				"--givenName", "Woodstock",
			)

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("Missing argument `email` must be specified."))
		})
	})

	Describe("CreateUserCmd", func() {
		It("performs POST with user data and bearer token", func() {
			server.RouteToHandler("POST", "/Users", CombineHandlers(
				RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
				VerifyRequest("POST", "/Users"),
				VerifyHeaderKV("Authorization", "bearer access_token"),
				VerifyHeaderKV("Accept", "application/json"),
				VerifyHeaderKV("Content-Type", "application/json"),
				VerifyJSON(`
					{
						"userName": "marcus",
						"password": "secret",
						"origin": "uaa",
						"name" : { "givenName" : "Marcus", "familyName" : "Aurelius" },
						"emails": [
							{
								"value": "marcus@philosophy.com",
								"primary": true
							},
							{
								"value": "marcusA@gmail.com",
								"primary": false
							}
						],
						"phoneNumbers": [{
							"value": "555-5555"
						},
						{
							"value": "666-6666"
						}]
					}
				`),
			))

			session := runCommand("create-user", "marcus",
				"--givenName", "Marcus",
				"--familyName", "Aurelius",
				"--email", "marcus@philosophy.com",
				"--email", "marcusA@gmail.com",
				"--phone", "555-5555",
				"--phone", "666-6666",
				"--password", "secret",
				"--origin", "uaa",
			)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(0))
		})

		It("prints the created user json", func() {
			server.RouteToHandler("POST", "/Users", CombineHandlers(
				RespondWith(http.StatusOK, fixtures.MarcusUserResponse),
				VerifyRequest("POST", "/Users"),
			))

			session := runCommand("create-user", "marcus",
				"--givenName", "Marcus",
				"--familyName", "Aurelius",
				"--email", "marcus@philosophy.com",
				"--email", "marcusA@gmail.com",
			)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(fixtures.MarcusUserResponse))
		})

		It("displays an error if there is a problem during create", func() {
			server.RouteToHandler("POST", "/Users", CombineHandlers(
				RespondWith(http.StatusBadRequest, ""),
				VerifyRequest("POST", "/Users"),
			))

			session := runCommand("create-user", "marcus",
				"--givenName", "Marcus",
				"--familyName", "Aurelius",
				"--email", "marcus@philosophy.com",
				"--email", "marcusA@gmail.com",
			)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(1))
		})
	})
})
