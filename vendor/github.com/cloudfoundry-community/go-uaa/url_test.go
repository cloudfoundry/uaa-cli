package uaa_test

import (
	"log"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildSubdomainURL(t *testing.T) {
	spec.Run(t, "BuildSubdomainURL", testBuildSubdomainURL, spec.Report(report.Terminal{}))
}

func testBuildSubdomainURL(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
		log.SetFlags(log.Lshortfile)
	})

	it("returns a URL", func() {
		url, err := uaa.BuildSubdomainURL("http://test.example.com", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(url).NotTo(BeNil())
		Expect(url.String()).To(Equal("http://test.example.com"))
	})

	it("returns an error when the url is invalid", func() {
		url, err := uaa.BuildSubdomainURL("(*#&^@%$&%)", "")
		Expect(err).To(HaveOccurred())
		Expect(url).To(BeNil())
	})

	when("the zone ID is set", func() {
		it("adds the zone ID as a prefix to the target", func() {
			testCases := []struct {
				target   string
				zoneID   string
				expected string
			}{
				{"http://test.example.com", "zone1", "http://zone1.test.example.com"},
				{"https://test.example.com", "zone1", "https://zone1.test.example.com"},
				{"test.example.com", "zone1", "https://zone1.test.example.com"},
			}
			for i := range testCases {
				url, err := uaa.BuildSubdomainURL(testCases[i].target, testCases[i].zoneID)
				Expect(err).NotTo(HaveOccurred())
				Expect(url).NotTo(BeNil())
				Expect(url.String()).To(Equal(testCases[i].expected))
			}
		})

		it("returns an error when the url is invalid", func() {
			url, err := uaa.BuildSubdomainURL("(*#&^@%$&%)", "zone1")
			Expect(err).To(HaveOccurred())
			Expect(url).To(BeNil())
		})
	})
}
