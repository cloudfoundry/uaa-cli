package utils_test

import (
	. "code.cloudfoundry.org/uaa-cli/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Arrayify", func() {
	It("Can handle whitespace", func() {
		Expect(Arrayify("")).To(Equal([]string{}))
		Expect(Arrayify(" ")).To(Equal([]string{}))
		Expect(Arrayify("  ")).To(Equal([]string{}))
	})

	It("Can handle a single item", func() {
		Expect(Arrayify("foo")).To(Equal([]string{"foo"}))
		Expect(Arrayify("foo ")).To(Equal([]string{"foo"}))
	})

	It("Can split comma-separated string into slice of strings", func() {
		Expect(Arrayify("foo,bar,baz")).To(Equal([]string{"foo", "bar", "baz"}))
	})

	It("Can split space-separated string into slice of strings", func() {
		Expect(Arrayify("foo bar baz")).To(Equal([]string{"foo", "bar", "baz"}))
	})

	It("Can split comma-space-separated string into slice of strings", func() {
		Expect(Arrayify("foo, bar, baz")).To(Equal([]string{"foo", "bar", "baz"}))
	})

	It("Can trim extra whitespace", func() {
		Expect(Arrayify(" foo bar baz ")).To(Equal([]string{"foo", "bar", "baz"}))
		Expect(Arrayify("foo,bar,baz ")).To(Equal([]string{"foo", "bar", "baz"}))
	})

	It("Can trim extra spaces between entries", func() {
		Expect(Arrayify("foo bar    baz")).To(Equal([]string{"foo", "bar", "baz"}))
		Expect(Arrayify("foo,   bar, baz")).To(Equal([]string{"foo", "bar", "baz"}))
	})

	It("Can remove empty entries", func() {
		Expect(Arrayify("foo,,bar,baz")).To(Equal([]string{"foo", "bar", "baz"}))
	})
})
