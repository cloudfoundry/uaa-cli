package uaa_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testCurl(t *testing.T, when spec.G, it spec.S) {
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

	it("gets a user from the UAA by id", func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.URL.Path).To(Equal("/Users/00000000-0000-0000-0000-000000000001"))
			Expect(req.Method).To(Equal(http.MethodGet))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userResponse))
		})

		_, resBody, err := a.Curl("/Users/00000000-0000-0000-0000-000000000001", "GET", "", []string{"Accept: application/json"})
		Expect(err).NotTo(HaveOccurred())

		var user uaa.User
		err = json.Unmarshal([]byte(resBody), &user)
		Expect(err).NotTo(HaveOccurred())

		Expect(user.ID).To(Equal("00000000-0000-0000-0000-000000000001"))
	})

	it("can POST body and multiple headers", func() {
		reqBody := map[string]interface{}{
			"externalID": "marcus-user",
			"userName":   "marcus@stoicism.com",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
			Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(req.Method).To(Equal(http.MethodPost))
			Expect(req.URL.Path).To(Equal("/Users"))
			defer req.Body.Close()
			body, _ := ioutil.ReadAll(req.Body)
			Expect(body).To(MatchJSON(reqBodyBytes))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userResponse))
		})

		Expect(err).NotTo(HaveOccurred())

		_, resBody, _ := a.Curl("/Users", "POST", string(reqBodyBytes), []string{"Content-Type: application/json", "Accept: application/json"})

		var user uaa.User
		err = json.Unmarshal([]byte(resBody), &user)
		Expect(err).NotTo(HaveOccurred())

		Expect(user.ID).To(Equal("00000000-0000-0000-0000-000000000001"))
	})
}
