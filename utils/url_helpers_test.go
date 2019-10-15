package utils_test

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UrlHelpers", func() {
	Describe("BuildUrl", func() {
		It("adds path to base url", func() {
			url, _ := utils.BuildUrl("http://localhost:9090", "foo")
			Expect(url.String()).To(Equal("http://localhost:9090/foo"))

			url, _ = utils.BuildUrl("http://localhost:9090/", "foo")
			Expect(url.String()).To(Equal("http://localhost:9090/foo"))

			url, _ = utils.BuildUrl("http://localhost:9090/", "/foo")
			Expect(url.String()).To(Equal("http://localhost:9090/foo"))

			url, _ = utils.BuildUrl("http://localhost:9090", "/foo")
			Expect(url.String()).To(Equal("http://localhost:9090/foo"))
		})

		It("preserves the base path", func() {
			url, _ := utils.BuildUrl("http://localhost:9090/uaa", "foo")
			Expect(url.String()).To(Equal("http://localhost:9090/uaa/foo"))
		})
	})
})
