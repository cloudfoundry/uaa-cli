package utils_test

import (
	. "code.cloudfoundry.org/uaa-cli/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stringifier", func() {
	It("presents a string slice as a string", func() {
		Expect(StringSliceStringifier([]string{"foo", "bar", "baz"})).To(Equal("[foo, bar, baz]"))
		Expect(StringSliceStringifier([]string{"foo"})).To(Equal("[foo]"))
		Expect(StringSliceStringifier([]string{})).To(Equal("[]"))
		Expect(StringSliceStringifier([]string{" "})).To(Equal("[ ]"))
	})
})
