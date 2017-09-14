package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetUser", func() {
	var userManager uaa.TestUserCrud
	var printer cli.TestPrinter

	BeforeEach(func() {
		printer = cli.NewTestPrinter()
		userManager = uaa.NewTestUserCrud()
	})

	It("uses the UserManager to get a given userId", func() {
		GetUserCmd("jen", userManager, printer)

		Expect(userManager.CallData["GetId"]).To(Equal("jen"))
	})

	It("prints the user", func() {
		GetUserCmd("jen", userManager, printer)

		Expect(printer.CallData["Print"]).To(Equal(uaa.ScimUser{Id: "jen"}))
	})

	Describe("validations", func() {
		It("requires a target", func() {
			err := GetUserValidations(uaa.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})

		It("requires a user id", func() {
			cfg := uaa.NewConfigWithServerURL("http://localhost:8080")
			ctx := uaa.NewContextWithToken("access_token")
			cfg.AddContext(ctx)

			err := GetUserValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional argument USER_ID must be specified."))

			err = GetUserValidations(cfg, []string{"userid"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
