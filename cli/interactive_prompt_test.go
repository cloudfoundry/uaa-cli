package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"bufio"
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interactive inputs", func() {
	var inbuf *bytes.Buffer
	var outbuf *bytes.Buffer

	BeforeEach(func() {
		inbuf = bytes.NewBuffer([]byte{})
		outbuf = bytes.NewBuffer([]byte{})
		InteractiveOutput = outbuf
		InteractiveInput = inbuf
	})

	Describe("InteractivePrompt", func() {
		It("prints the prompt for the user", func() {
			ip := InteractivePrompt{Prompt: "Username"}

			ip.Get()

			outreader := bufio.NewReader(outbuf)
			printed, err := outreader.ReadString(':')
			Expect(err).NotTo(HaveOccurred())
			Expect(printed).To(ContainSubstring("Username:"))
		})

		It("gets user input", func() {
			inbuf.WriteString("woodstock\n")

			ip := InteractivePrompt{Prompt: "Username"}
			input, err := ip.Get()

			Expect(err).NotTo(HaveOccurred())
			Expect(input).To(Equal("woodstock"))
		})
	})

	Describe("InteractiveSecret", func() {
		It("gets user input with terminal.ReadPassword", func() {
			called := false
			ReadPassword = func(fd int) ([]byte, error) {
				called = true
				return []byte("somepassword"), nil
			}

			ip := InteractiveSecret{Prompt: "Password"}
			input, err := ip.Get()

			Expect(err).NotTo(HaveOccurred())
			Expect(input).To(Equal("somepassword"))
			Expect(called).To(BeTrue())
		})
	})
})
