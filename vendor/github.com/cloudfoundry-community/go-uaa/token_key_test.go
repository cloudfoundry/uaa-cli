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

func testTokenKey(t *testing.T, when spec.G, it spec.S) {
	var (
		s                 *httptest.Server
		handler           http.Handler
		called            int
		a                 *uaa.API
		asymmetricKeyJSON string
	)

	it.Before(func() {
		RegisterTestingT(t)
		asymmetricKeyJSON = `{
		  "kty": "RSA",
		  "e": "AQAB",
		  "use": "sig",
		  "kid": "sha2-2017-01-20-key",
		  "alg": "RS256",
		  "value": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----",
		  "n": "AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"
		}`
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

	it("calls the /token_key endpoint", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/token_key"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(asymmetricKeyJSON))
		})

		key, _ := a.TokenKey()

		Expect(called).To(Equal(1))
		Expect(key.Kty).To(Equal("RSA"))
		Expect(key.E).To(Equal("AQAB"))
		Expect(key.Use).To(Equal("sig"))
		Expect(key.Kid).To(Equal("sha2-2017-01-20-key"))
		Expect(key.Alg).To(Equal("RS256"))
		Expect(key.Value).To(Equal("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----"))
		Expect(key.N).To(Equal("AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"))
	})

	it("returns helpful error when /token_key request fails", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/token_key"))
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := a.TokenKey()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
	})

	it("returns helpful error when /token_key response can't be parsed", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/token_key"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{unparsable-json-response}"))
		})
		_, err := a.TokenKey()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
		Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
	})

	it("can handle symmetric keys", func() {
		symmetricKeyJSON := `{
		  "kty" : "MAC",
		  "alg" : "HS256",
		  "value" : "key",
		  "use" : "sig",
		  "kid" : "testKey"
		}`
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/token_key"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(symmetricKeyJSON))
		})
		key, _ := a.TokenKey()
		Expect(called).To(Equal(1))
		Expect(key.Kty).To(Equal("MAC"))
		Expect(key.Alg).To(Equal("HS256"))
		Expect(key.Value).To(Equal("key"))
		Expect(key.Use).To(Equal("sig"))
		Expect(key.Kid).To(Equal("testKey"))
	})
}
