package uaa

import (
	"net"
	"net/http"
	"testing"
	"time"

	"net/http/httptest"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"golang.org/x/oauth2"
)

func TestEnsureTransport(t *testing.T) {
	spec.Run(t, "EnsureTransport", testEnsureTransport, spec.Report(report.Terminal{}))
}

func testEnsureTransport(t *testing.T, when spec.G, it spec.S) {
	var a *API
	it.Before(func() {
		RegisterTestingT(t)
		a = &API{}
	})

	when("the client is nil", func() {
		it("is a no-op", func() {
			a.ensureTransport(a.UnauthenticatedClient)
			Expect(a.UnauthenticatedClient).To(BeNil())
		})
	})

	when("the authenticated client is not set but the unauthenticated client is set", func() {
		var s *httptest.Server

		it.Before(func() {
			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))
			a.UnauthenticatedClient = &http.Client{}
		})

		it.After(func() {
			if s != nil {
				s.Close()
			}
		})

		it("will make a http call with the unauthenticated client", func() {
			req, err := http.NewRequest("GET", s.URL, nil)
			Expect(err).NotTo(HaveOccurred())
			_, err = a.doAndRead(req, false)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	when("the client transport is not set", func() {
		it.Before(func() {
			a.UnauthenticatedClient = &http.Client{}
		})

		it("is a no-op", func() {
			a.ensureTransport(a.UnauthenticatedClient)
			Expect(a.UnauthenticatedClient).NotTo(BeNil())
			Expect(a.UnauthenticatedClient.Transport).To(BeNil())
		})
	})

	when("the client transport is an http.Transport", func() {
		it.Before(func() {
			a.UnauthenticatedClient = &http.Client{Transport: &http.Transport{}}
		})

		when("SkipSSLValidation is false", func() {
			it.Before(func() {
				a.SkipSSLValidation = false
			})

			it("will not initialize the TLS client config", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*http.Transport)
				Expect(t.TLSClientConfig).To(BeNil())
			})
		})

		when("SkipSSLValidation is true", func() {
			it.Before(func() {
				a.SkipSSLValidation = true
			})

			it("will initialize the TLS client config and set InsecureSkipVerify", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*http.Transport)
				Expect(t.TLSClientConfig).NotTo(BeNil())
				Expect(t.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
			})
		})
	})

	when("the client transport is a tokenTransport", func() {
		it.Before(func() {
			a.UnauthenticatedClient = &http.Client{Transport: &tokenTransport{
				underlyingTransport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
						DualStack: true,
					}).DialContext,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				},
			}}
		})

		when("SkipSSLValidation is false", func() {
			it.Before(func() {
				a.SkipSSLValidation = false
			})

			it("will not initialize the TLS client config", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*tokenTransport)
				Expect(t.underlyingTransport.TLSClientConfig).To(BeNil())
			})
		})

		when("SkipSSLValidation is true", func() {
			it.Before(func() {
				a.SkipSSLValidation = true
			})

			it("will initialize the TLS client config and set InsecureSkipVerify", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*tokenTransport)
				Expect(t.underlyingTransport.TLSClientConfig).NotTo(BeNil())
				Expect(t.underlyingTransport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
			})
		})
	})

	when("the client transport is an oauth2.Transport but the Base transport is nil", func() {
		it.Before(func() {
			a.UnauthenticatedClient = &http.Client{Transport: &oauth2.Transport{}}
		})

		it("is a no-op", func() {
			a.ensureTransport(a.UnauthenticatedClient)
			Expect(a.UnauthenticatedClient).NotTo(BeNil())
			Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
			t := a.UnauthenticatedClient.Transport.(*oauth2.Transport)
			Expect(t.Base).To(BeNil())
		})
	})

	when("the client transport is an oauth2.Transport with a Base transport", func() {
		it.Before(func() {
			a.UnauthenticatedClient = &http.Client{Transport: &oauth2.Transport{
				Base: &http.Transport{},
			}}
		})

		when("SkipSSLValidation is false", func() {
			it.Before(func() {
				a.SkipSSLValidation = false
			})

			it("will not initialize the TLS client config if SkipSSLValidation is false", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*oauth2.Transport)
				Expect(t.Base).NotTo(BeNil())
				b := t.Base.(*http.Transport)
				Expect(b.TLSClientConfig).To(BeNil())
			})
		})

		when("SkipSSLValidation is true", func() {
			it.Before(func() {
				a.SkipSSLValidation = true
			})

			it("will initialize the TLS client config and set InsecureSkipVerify", func() {
				a.ensureTransport(a.UnauthenticatedClient)
				Expect(a.UnauthenticatedClient).NotTo(BeNil())
				Expect(a.UnauthenticatedClient.Transport).NotTo(BeNil())
				t := a.UnauthenticatedClient.Transport.(*oauth2.Transport)
				Expect(t.Base).NotTo(BeNil())
				b := t.Base.(*http.Transport)
				Expect(b.TLSClientConfig).NotTo(BeNil())
				Expect(b.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
			})
		})
	})
}
