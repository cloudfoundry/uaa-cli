package uaa

import (
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testURLWithPath(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
		log.SetFlags(log.Lshortfile)
	})

	it("returns a URL which retains the path", func() {
		url, err := url.Parse("http://example.com/uaa")
		Expect(url).NotTo(BeNil())
		Expect(err).To(BeNil())

		withPath := urlWithPath(*url, "path")
		Expect(withPath.String()).To(Equal("http://example.com/uaa/path"))
	})
}
