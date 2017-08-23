package utils_test

import (
	. "github.com/jhamon/guac/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UrlHelpers", func() {
	Describe("BuildUrl", func() {
		It("adds path to base url", func() {
			url, _ := BuildUrl("http://localhost:8080", "foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = BuildUrl("http://localhost:8080/", "foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = BuildUrl("http://localhost:8080/", "/foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))

			url, _ = BuildUrl("http://localhost:8080", "/foo")
			Expect(url.String()).To(Equal("http://localhost:8080/foo"))
		})
	})
})
