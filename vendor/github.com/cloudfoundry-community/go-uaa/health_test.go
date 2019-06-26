package uaa_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testIsHealthy(t *testing.T, when spec.G, it spec.S) {
	var (
		s       *httptest.Server
		handler http.Handler
		called  int
		a       *uaa.API
	)

	it.Before(func() {
		RegisterTestingT(t)
		called = 0
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			called = called + 1
			Expect(handler).NotTo(BeNil())
			handler.ServeHTTP(w, req)
		}))
		c := &http.Client{Transport: http.DefaultTransport}
		u, _ := url.Parse(s.URL)
		a = &uaa.API{
			TargetURL:             u,
			AuthenticatedClient:   c,
			UnauthenticatedClient: c,
		}
	})

	it.After(func() {
		if s != nil {
			s.Close()
		}
	})

	it("is healthy when a 200 response is received", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.URL.Path).To(Equal("/healthz"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		status, err := a.IsHealthy()
		Expect(status).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	it("is unhealthy when a non-200 response is received", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.URL.Path).To(Equal("/healthz"))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ok"))
		})
		status, err := a.IsHealthy()
		Expect(status).To(BeFalse())
		Expect(err).NotTo(HaveOccurred())
	})
}
