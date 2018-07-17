package uaa_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestTokenKeys(t *testing.T) {
	spec.Run(t, "TokenKey", testTokenKeys, spec.Report(report.Terminal{}))
}

func testTokenKeys(t *testing.T, when spec.G, it spec.S) {
	var (
		s             *httptest.Server
		handler       http.Handler
		called        int
		a             *uaa.API
		tokenKeysJSON string
	)

	it.Before(func() {
		RegisterTestingT(t)
		tokenKeysJSON = `{
			"keys": [
			{
				"kty": "RSA",
				"e": "AQAB",
				"use": "sig",
				"kid": "sha2-2017-01-20-key",
				"alg": "RS256",
				"value": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----",
				"n": "AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"
			},
			{
				"kty": "RSA",
				"e": "AQAB",
				"use": "sig",
				"kid": "legacy-token-key",
				"alg": "RS256",
				"value": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8/aXmEImpdwWHJlYc4G8\n3BgZVmyhCdy7SCL0kM7wV5xCvRKK0k4nKjH0QW2E+0GIKzIj4JQhYU+MeZHrArfC\nrfthIXcio/Ll6NvoTPY77XA7U6vBGCiLdGYSGrN8y064cF2uM8d3AEgTT0RzLK3E\n688Ltq38GxnoXOUuLZmXS2HeHNd2bW/k6Eyd9Z3ymmdpVZXMyLwepNxU38WQS2bJ\nPYXYvRkzoZ1ru/deExwbecI18NCeO/GKp3f8bwKuC2j3LKFJIAwW3zFoDrcAxpC/\nJDG2RSTj//CRvhtd7JkeQLVKGyIHNtACaPT3tFT6scvVXHGPB5fRTLB8Lr+mK4RI\nBwIDAQAB\n-----END PUBLIC KEY-----",
				"n": "APP2l5hCJqXcFhyZWHOBvNwYGVZsoQncu0gi9JDO8FecQr0SitJOJyox9EFthPtBiCsyI-CUIWFPjHmR6wK3wq37YSF3IqPy5ejb6Ez2O-1wO1OrwRgoi3RmEhqzfMtOuHBdrjPHdwBIE09EcyytxOvPC7at_BsZ6FzlLi2Zl0th3hzXdm1v5OhMnfWd8ppnaVWVzMi8HqTcVN_FkEtmyT2F2L0ZM6Gda7v3XhMcG3nCNfDQnjvxiqd3_G8Crgto9yyhSSAMFt8xaA63AMaQvyQxtkUk4__wkb4bXeyZHkC1ShsiBzbQAmj097RU-rHL1VxxjweX0UywfC6_piuESAc"
			}
			]
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

	it("calls the /token_keys endpoint", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/token_keys"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenKeysJSON))
		})
		keys, _ := a.TokenKeys()
		Expect(called).To(Equal(1))
		Expect(keys[0].Kty).To(Equal("RSA"))
		Expect(keys[0].E).To(Equal("AQAB"))
		Expect(keys[0].Use).To(Equal("sig"))
		Expect(keys[0].Kid).To(Equal("sha2-2017-01-20-key"))
		Expect(keys[0].Alg).To(Equal("RS256"))
		Expect(keys[0].Value).To(Equal("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----"))
		Expect(keys[0].N).To(Equal("AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"))
		Expect(keys[1].Kid).To(Equal("legacy-token-key"))
	})

	it("returns a helpful error when response cannot be parsed", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{unparsable}"))
		})
		_, err := a.TokenKeys()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
	})

	when("the server is an older UAA that is missing the /token_keys endpoint", func() {
		var tokenKeyJSON string = `{
		  "kty": "RSA",
		  "e": "AQAB",
		  "use": "sig",
		  "kid": "sha2-2017-01-20-key",
		  "alg": "RS256",
		  "value": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----",
		  "n": "AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"
		}`

		it("falls back to /token_key endpoint", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				if req.URL.Path == "/token_keys" {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tokenKeyJSON))
				}
			})
			keys, _ := a.TokenKeys()
			Expect(called).To(Equal(2))
			Expect(keys).To(HaveLen(1))
			Expect(keys[0].Kty).To(Equal("RSA"))
			Expect(keys[0].E).To(Equal("AQAB"))
			Expect(keys[0].Use).To(Equal("sig"))
			Expect(keys[0].Kid).To(Equal("sha2-2017-01-20-key"))
			Expect(keys[0].Alg).To(Equal("RS256"))
			Expect(keys[0].Value).To(Equal("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----"))
			Expect(keys[0].N).To(Equal("AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"))
		})
	})
}
