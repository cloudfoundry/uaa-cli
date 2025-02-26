package cli_test

import (
	"bytes"
	. "code.cloudfoundry.org/uaa-cli/cli"
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("JsonPrinter", func() {
	var infoLogBuf, robotLogBuf, warnLogBuf, errorLogBuf *Buffer
	var printer JsonPrinter

	BeforeEach(func() {
		infoLogBuf, robotLogBuf, warnLogBuf, errorLogBuf = NewBuffer(), NewBuffer(), NewBuffer(), NewBuffer()
		printer = NewJsonPrinter(NewLogger(infoLogBuf, robotLogBuf, warnLogBuf, errorLogBuf))
	})
	Describe("Print", func() {
		It("prints things to the Robots log", func() {
			printer.Print(struct {
				Foo string
				Bar string
			}{"foo", "bar"})

			Expect(robotLogBuf.Contents()).To(MatchJSON(`{"Foo":"foo","Bar":"bar"}`))
		})

		It("returns error when cannot marhsal into json", func() {
			unJsonifiableObj := make(chan bool)
			err := printer.Print(unJsonifiableObj)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PrintError", func() {
		It("prints a json buffer to Error log", func() {
			jsonData := struct {
				Foo string
				Bar string
			}{"foo", "bar"}
			jsonRaw, _ := json.Marshal(jsonData)
			var out bytes.Buffer
			_ = json.Indent(&out, jsonRaw, "", "  ")

			printer.PrintError(jsonRaw)

			Expect(string(errorLogBuf.Contents())).To(ContainSubstring(out.String()))
		})
	})
})
