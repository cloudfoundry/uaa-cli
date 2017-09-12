package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("GetUser", func() {
	var stdOut io.Writer
	var userManager uaa.TestUserCrud
	var printer cli.TestPrinter

	BeforeEach(func() {
		stdOut = GinkgoWriter
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

})
