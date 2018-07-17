package uaa

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestContains(t *testing.T) {
	spec.Run(t, "Contains", testContains, spec.Report(report.Terminal{}))
}

func testContains(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	list := []string{"do", "re", "mi"}

	it("returns true if present", func() {
		Expect(contains(list, "re")).To(BeTrue())
	})

	it("returns false if not present", func() {
		Expect(contains(list, "fa")).To(BeFalse())
	})

	it("handles empty list", func() {
		Expect(contains([]string{}, "fa")).To(BeFalse())
	})
}
