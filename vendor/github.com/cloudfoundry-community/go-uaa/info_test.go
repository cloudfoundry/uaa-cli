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

const InfoResponseJSON string = `{
	"app": {
	"version": "4.5.0"
	},
	"links": {
	"uaa": "https://uaa.run.pivotal.io",
	"passwd": "https://account.run.pivotal.io/forgot-password",
	"login": "https://login.run.pivotal.io",
	"register": "https://account.run.pivotal.io/sign-up"
	},
	"zone_name": "uaa",
	"entityID": "login.run.pivotal.io",
	"commit_id": "df80f63",
	"idpDefinitions": {
	 "SAML" : "http://localhost:8080/uaa/saml/discovery?returnIDParam=idp&entityID=cloudfoundry-saml-login&idp=SAML&isPassive=true"
	},
	"prompts": {
	"username": [
		"text",
		"Email"
	],
	"password": [
		"password",
		"Password"
	]
	},
	"timestamp": "2017-07-21T22:45:01+0000"
}`

func testInfo(t *testing.T, when spec.G, it spec.S) {
	var (
		s            *httptest.Server
		a            *uaa.API
		h            http.Handler
		handlerCalls int
	)

	it.Before(func() {
		RegisterTestingT(t)
		handlerCalls = 0
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handlerCalls = handlerCalls + 1
			h.ServeHTTP(w, req)
		}))
		url, _ := url.Parse(s.URL)
		client := &http.Client{}
		a = &uaa.API{
			UnauthenticatedClient: client,
			AuthenticatedClient:   client,
			TargetURL:             url,
		}
	})

	it.After(func() {
		if s != nil {
			s.Close()
		}
	})

	when("the info endpoint responds with valid info", func() {
		it.Before(func() {
			h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.Path).To(Equal("/info"))
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				w.Write([]byte(InfoResponseJSON))
			})
		})

		it("calls the /info endpoint", func() {
			infoResponse, _ := a.GetInfo()
			Expect(handlerCalls).To(Equal(1))
			Expect(infoResponse.App.Version).To(Equal("4.5.0"))
			Expect(infoResponse.Links.ForgotPassword).To(Equal("https://account.run.pivotal.io/forgot-password"))
			Expect(infoResponse.Links.Uaa).To(Equal("https://uaa.run.pivotal.io"))
			Expect(infoResponse.Links.Registration).To(Equal("https://account.run.pivotal.io/sign-up"))
			Expect(infoResponse.Links.Login).To(Equal("https://login.run.pivotal.io"))
			Expect(infoResponse.ZoneName).To(Equal("uaa"))
			Expect(infoResponse.EntityID).To(Equal("login.run.pivotal.io"))
			Expect(infoResponse.CommitID).To(Equal("df80f63"))
			Expect(infoResponse.IdpDefinitions["SAML"]).To(Equal("http://localhost:8080/uaa/saml/discovery?returnIDParam=idp&entityID=cloudfoundry-saml-login&idp=SAML&isPassive=true"))
			Expect(infoResponse.Prompts["username"]).To(Equal([]string{"text", "Email"}))
			Expect(infoResponse.Prompts["password"]).To(Equal([]string{"password", "Password"}))
			Expect(infoResponse.Timestamp).To(Equal("2017-07-21T22:45:01+0000"))
		})
	})

	when("the info endpoint responds with an error", func() {
		it.Before(func() {
			h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.Path).To(Equal("/info"))
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				w.WriteHeader(http.StatusInternalServerError)
			})
		})

		it("returns a helpful error", func() {
			_, err := a.GetInfo()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
		})
	})

	when("the info endpoint responds with unparsable JSON", func() {
		it.Before(func() {
			h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.Path).To(Equal("/info"))
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				w.Write([]byte("{unparsable-json-response}"))
			})
		})

		it("returns a helpful error", func() {
			_, err := a.GetInfo()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
			Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
		})
	})
}
