package utils_test

import (
	. "code.cloudfoundry.org/uaa-cli/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contains", func() {
	list := []string{"do", "re", "mi"}

	It("returns true if present", func() {
		Expect(Contains(list, "re")).To(BeTrue())
	})

	It("returns false if not present", func() {
		Expect(Contains(list, "fa")).To(BeFalse())
	})

	It("handles empty list", func() {
		Expect(Contains([]string{}, "fa")).To(BeFalse())
	})
})
