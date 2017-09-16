package utils_test

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UrlHelpers", func() {
	Describe("BuildUrl", func() {
		It("adds path to base url", func() {
			url, _ := utils.BuildUrl("http://localhost:8080", "foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = utils.BuildUrl("http://localhost:8080/", "foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = utils.BuildUrl("http://localhost:8080/", "/foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = utils.BuildUrl("http://localhost:8080", "/foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))
		})
	})
})
