package uaa_test

import (
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("TokenKey", func() {
	var (
		server            *ghttp.Server
		client            *http.Client
		config            Config
		asymmetricKeyJson string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = &http.Client{}
		config = NewConfigWithServerURL(server.URL())

		asymmetricKeyJson = `{
		  "kty": "RSA",
		  "e": "AQAB",
		  "use": "sig",
		  "kid": "sha2-2017-01-20-key",
		  "alg": "RS256",
		  "value": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----",
		  "n": "AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"
		}`
	})

	AfterEach(func() {
		server.Close()
	})

	It("calls the /token_key endpoint", func() {
		server.RouteToHandler("GET", "/token_key", ghttp.CombineHandlers(
			ghttp.RespondWith(200, asymmetricKeyJson),
			ghttp.VerifyRequest("GET", "/token_key"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
		))

		key, _ := TokenKey(client, config)

		Expect(server.ReceivedRequests()).To(HaveLen(1))
		Expect(key.Kty).To(Equal("RSA"))
		Expect(key.E).To(Equal("AQAB"))
		Expect(key.Use).To(Equal("sig"))
		Expect(key.Kid).To(Equal("sha2-2017-01-20-key"))
		Expect(key.Alg).To(Equal("RS256"))
		Expect(key.Value).To(Equal("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyH6kYCP29faDAUPKtei3\nV/Zh8eCHyHRDHrD0iosvgHuaakK1AFHjD19ojuPiTQm8r8nEeQtHb6mDi1LvZ03e\nEWxpvWwFfFVtCyBqWr5wn6IkY+ZFXfERLn2NCn6sMVxcFV12sUtuqD+jrW8MnTG7\nhofQqxmVVKKsZiXCvUSzfiKxDgoiRuD3MJSoZ0nQTHVmYxlFHuhTEETuTqSPmOXd\n/xJBVRi5WYCjt1aKRRZEz04zVEBVhVkr2H84qcVJHcfXFu4JM6dg0nmTjgd5cZUN\ncwA1KhK2/Qru9N0xlk9FGD2cvrVCCPWFPvZ1W7U7PBWOSBBH6GergA+dk2vQr7Ho\nlQIDAQAB\n-----END PUBLIC KEY-----"))
		Expect(key.N).To(Equal("AMh-pGAj9vX2gwFDyrXot1f2YfHgh8h0Qx6w9IqLL4B7mmpCtQBR4w9faI7j4k0JvK_JxHkLR2-pg4tS72dN3hFsab1sBXxVbQsgalq-cJ-iJGPmRV3xES59jQp-rDFcXBVddrFLbqg_o61vDJ0xu4aH0KsZlVSirGYlwr1Es34isQ4KIkbg9zCUqGdJ0Ex1ZmMZRR7oUxBE7k6kj5jl3f8SQVUYuVmAo7dWikUWRM9OM1RAVYVZK9h_OKnFSR3H1xbuCTOnYNJ5k44HeXGVDXMANSoStv0K7vTdMZZPRRg9nL61Qgj1hT72dVu1OzwVjkgQR-hnq4APnZNr0K-x6JU"))
	})

	It("returns helpful error when /token_key request fails", func() {
		server.RouteToHandler("GET", "/token_key", ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/token_key"),
			ghttp.RespondWith(500, "error response"),
			ghttp.VerifyRequest("GET", "/token_key"),
		))

		_, err := TokenKey(client, config)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
	})

	It("returns helpful error when /token_key response can't be parsed", func() {
		server.RouteToHandler("GET", "/token_key", ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/token_key"),
			ghttp.RespondWith(200, "{unparsable-json-response}"),
			ghttp.VerifyRequest("GET", "/token_key"),
		))

		_, err := TokenKey(client, config)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
		Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
	})

	It("can handle symmetric keys", func() {
		symmetricKeyJson := `{
		  "kty" : "MAC",
		  "alg" : "HS256",
		  "value" : "key",
		  "use" : "sig",
		  "kid" : "testKey"
		}`

		server.RouteToHandler("GET", "/token_key", ghttp.CombineHandlers(
			ghttp.RespondWith(200, symmetricKeyJson),
			ghttp.VerifyRequest("GET", "/token_key"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
		))

		key, _ := TokenKey(client, config)

		Expect(server.ReceivedRequests()).To(HaveLen(1))
		Expect(key.Kty).To(Equal("MAC"))
		Expect(key.Alg).To(Equal("HS256"))
		Expect(key.Value).To(Equal("key"))
		Expect(key.Use).To(Equal("sig"))
		Expect(key.Kid).To(Equal("testKey"))
	})
})
