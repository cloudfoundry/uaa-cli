package uaa_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

const clientResponse string = `{
	"scope" : [ "clients.read", "clients.write" ],
	"client_id" : "00000000-0000-0000-0000-000000000001",
	"resource_ids" : [ "none" ],
	"authorized_grant_types" : [ "client_credentials" ],
	"redirect_uri" : [ "http://ant.path.wildcard/**/passback/*", "http://test1.com" ],
	"autoapprove" : [ "true" ],
	"authorities" : [ "clients.read", "clients.write" ],
	"token_salt" : "1SztLL",
	"allowedproviders" : [ "uaa", "ldap", "my-saml-provider" ],
	"name" : "My Client Name",
	"lastModified" : 1502816030525,
	"required_user_groups" : [ ]
}`

const clientListResponse = `{
	"resources" : [ {
	"client_id" : "00000000-0000-0000-0000-000000000001"
	},
	{
	"client_id" : "00000000-0000-0000-0000-000000000002"
	}],
	"startIndex" : 1,
	"itemsPerPage" : 2,
	"totalResults" : 6,
	"schemas" : [ "http://cloudfoundry.org/schema/scim/oauth-clients-1.0" ]
}`

var testClientValue uaa.Client = uaa.Client{
	ClientID:     "00000000-0000-0000-0000-000000000001",
	ClientSecret: "new_secret",
}

const testClientJSON string = `{"client_id": "00000000-0000-0000-0000-000000000001", "client_secret": "new_secret"}`

func TestClientExtra(t *testing.T) {
	spec.Run(t, "ClientExtra", testClientExtra, spec.Report(report.Terminal{}))
}

func testClientExtra(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("Client.Validate()", func() {
		it("rejects empty grant types", func() {
			client := uaa.Client{}
			err := client.Validate()
			Expect(err.Error()).To(Equal(`grant type must be one of [authorization_code, implicit, password, client_credentials]`))
		})

		when("when authorization_code", func() {
			it("requires client_id", func() {
				client := uaa.Client{
					AuthorizedGrantTypes: []string{"authorization_code"},
					RedirectURI:          []string{"http://localhost:8080"},
					ClientSecret:         "secret",
				}

				err := client.Validate()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("client_id must be specified in the client definition"))
			})

			it("requires redirect_uri", func() {
				client := uaa.Client{
					ClientID:             "myclient",
					AuthorizedGrantTypes: []string{"authorization_code"},
					ClientSecret:         "secret",
				}
				err := client.Validate()
				Expect(err.Error()).To(Equal("redirect_uri must be specified for authorization_code grant type"))
			})

			it("requires client_secret", func() {
				client := uaa.Client{
					ClientID:             "myclient",
					AuthorizedGrantTypes: []string{"authorization_code"},
					RedirectURI:          []string{"http://localhost:8080"},
				}
				err := client.Validate()
				Expect(err.Error()).To(Equal("client_secret must be specified for authorization_code grant type"))
			})
		})

		when("when implicit", func() {
			it("requires client_id", func() {
				client := uaa.Client{
					AuthorizedGrantTypes: []string{"implicit"},
					RedirectURI:          []string{"http://localhost:8080"},
				}
				err := client.Validate()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("client_id must be specified in the client definition"))
			})

			it("requires redirect_uri", func() {
				client := uaa.Client{
					ClientID:             "myclient",
					AuthorizedGrantTypes: []string{"implicit"},
				}
				err := client.Validate()
				Expect(err.Error()).To(Equal("redirect_uri must be specified for implicit grant type"))
			})

			it("does not require client_secret", func() {
				client := uaa.Client{
					ClientID:             "someclient",
					AuthorizedGrantTypes: []string{"implicit"},
					RedirectURI:          []string{"http://localhost:8080"},
				}
				err := client.Validate()
				Expect(err).To(BeNil())
			})
		})

		when("when client_credentials", func() {
			it("requires client_id", func() {
				client := uaa.Client{
					AuthorizedGrantTypes: []string{"client_credentials"},
					ClientSecret:         "secret",
				}
				err := client.Validate()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("client_id must be specified in the client definition"))
			})

			it("requires client_secret", func() {
				client := uaa.Client{
					ClientID:             "myclient",
					AuthorizedGrantTypes: []string{"client_credentials"},
				}
				err := client.Validate()
				Expect(err.Error()).To(Equal("client_secret must be specified for client_credentials grant type"))
			})
		})

		when("when password", func() {
			it("requires client_id", func() {
				client := uaa.Client{
					AuthorizedGrantTypes: []string{"password"},
					RedirectURI:          []string{"http://localhost:8080"},
					ClientSecret:         "secret",
				}
				err := client.Validate()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("client_id must be specified in the client definition"))
			})

			it("requires client_secret", func() {
				client := uaa.Client{
					ClientID:             "myclient",
					AuthorizedGrantTypes: []string{"password"},
				}
				err := client.Validate()
				Expect(err.Error()).To(Equal("client_secret must be specified for password grant type"))
			})
		})
	})

	when("ChangeClientSecret()", func() {
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

		it("calls the /oauth/clients/<clientid>/secret endpoint", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Method).To(Equal(http.MethodPut))
				Expect(req.URL.Path).To(Equal(uaa.ClientsEndpoint + "/00000000-0000-0000-0000-000000000001/secret"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{"clientId": "00000000-0000-0000-0000-000000000001", "secret": "new_secret"}`))
				w.WriteHeader(http.StatusOK)
			})
			a.ChangeClientSecret("00000000-0000-0000-0000-000000000001", "new_secret")
			Expect(called).To(Equal(1))
		})

		it("does not panic when error happens during network call", func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Accept")).To(Equal("application/json"))
				Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
				Expect(req.Method).To(Equal(http.MethodPut))
				Expect(req.URL.Path).To(Equal(uaa.ClientsEndpoint + "/00000000-0000-0000-0000-000000000001/secret"))
				defer req.Body.Close()
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(MatchJSON(`{"clientId": "00000000-0000-0000-0000-000000000001", "secret": "new_secret"}`))
				w.WriteHeader(http.StatusUnauthorized)
			})
			err := a.ChangeClientSecret("00000000-0000-0000-0000-000000000001", "new_secret")
			Expect(called).To(Equal(1))
			Expect(err).NotTo(BeNil())
		})
	})
}
