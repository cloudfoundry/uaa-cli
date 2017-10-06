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

var _ = Describe("CreateGroup", func() {
	BeforeEach(func() {
		cfg := uaa.NewConfigWithServerURL(server.URL())
		cfg.AddContext(uaa.NewContextWithToken("access_token"))
		config.WriteConfig(cfg)
	})

	Describe("Validations", func() {
		It("requires a target to have been set", func() {
			config.WriteConfig(uaa.NewConfig())

			session := runCommand("create-group")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cmd.MISSING_TARGET))
		})

		It("requires a token in context", func() {
			config.WriteConfig(uaa.NewConfigWithServerURL(server.URL()))

			session := runCommand("create-group")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say(cmd.MISSING_CONTEXT))
		})

		It("requires a groupname", func() {
			session := runCommand("create-group")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("The positional argument GROUPNAME must be specified."))
		})
	})

	Describe("CreateGroupCmd", func() {
		It("performs POST with user data and bearer token", func() {
			reqBody := map[string]interface{}{
				"displayName": "admin",
				"members":     nil,
			}
			server.RouteToHandler("POST", "/Groups", CombineHandlers(
				RespondWith(http.StatusOK, fixtures.AdminGroupResponse),
				VerifyRequest("POST", "/Groups"),
				VerifyHeaderKV("Authorization", "bearer access_token"),
				VerifyHeaderKV("Accept", "application/json"),
				VerifyHeaderKV("Content-Type", "application/json"),
				VerifyJSONRepresenting(reqBody),
			))

			session := runCommand("create-group", "admin")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(0))
		})

		It("prints the created user json", func() {
			server.RouteToHandler("POST", "/Groups", CombineHandlers(
				RespondWith(http.StatusOK, fixtures.AdminGroupResponse),
				VerifyRequest("POST", "/Groups"),
			))

			session := runCommand("create-group", "admin")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(fixtures.AdminGroupResponse))
		})

		It("displays an error if there is a problem during create", func() {
			server.RouteToHandler("POST", "/Groups", CombineHandlers(
				RespondWith(http.StatusBadRequest, ""),
				VerifyRequest("POST", "/Groups"),
			))

			session := runCommand("create-group", "admin")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(session).To(Exit(1))
		})
	})
})
