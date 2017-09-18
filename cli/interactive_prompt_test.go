package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"bufio"
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InteractivePrompt", func() {
	var inbuf *bytes.Buffer
	var outbuf *bytes.Buffer

	BeforeEach(func() {
		inbuf = bytes.NewBuffer([]byte{})
		outbuf = bytes.NewBuffer([]byte{})
		InteractiveOutput = outbuf
		InteractiveInput = inbuf
	})

	It("prints the prompt for the user", func() {
		ip := InteractivePrompt{Prompt: "Username"}

		ip.Get()

		outreader := bufio.NewReader(outbuf)
		printed, err := outreader.ReadString(':')
		Expect(err).NotTo(HaveOccurred())
		Expect(printed).To(Equal("Username:"))
	})

	It("gets user input", func() {
		inbuf.WriteString("woodstock\n")

		ip := InteractivePrompt{Prompt: "Username"}
		input, err := ip.Get()

		Expect(err).NotTo(HaveOccurred())
		Expect(input).To(Equal("woodstock"))
	})
})
