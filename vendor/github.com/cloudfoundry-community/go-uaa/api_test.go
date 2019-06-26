package uaa_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"io/ioutil"

	uaa "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"golang.org/x/oauth2"
)

func testNew(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("TokenFormat.String()", func() {
		it("prints the string representation appropriately", func() {
			var t uaa.TokenFormat
			Expect(t.String()).To(Equal("opaque"))
			t = 3
			Expect(t.String()).To(Equal(""))
			Expect(uaa.JSONWebToken.String()).To(Equal("jwt"))
			Expect(uaa.OpaqueToken.String()).To(Equal("opaque"))
		})
	})

	when("New()", func() {
		it("returns an API even if the target is an invalid URL", func() {
			api := uaa.New("(*#&^@%$&%)", "")
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL).To(BeNil())
		})

		it("sets the TargerURL and zone", func() {
			api := uaa.New("https://example.net", "zone-1")
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal("https://example.net"))
			Expect(api.ZoneID).To(Equal("zone-1"))
		})

		it("Token() fails because when there is no mechanism to get a token", func() {
			api := uaa.New("https://example.net", "zone-1")
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(t).To(BeNil())
		})
	})

	when("NewWithToken()", func() {
		it("fails if the target url is invalid", func() {
			api, err := uaa.NewWithToken("(*#&^@%$&%)", "", oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: "test-token"})
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("fails if the token is invalid", func() {
			api, err := uaa.NewWithToken("https://example.net", "", oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: ""})
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
			api, err = uaa.NewWithToken("https://example.net", "", oauth2.Token{Expiry: time.Now().Add(-20 * time.Second), AccessToken: "test-token"})
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("returns an API with a TargetURL", func() {
			api, err := uaa.NewWithToken("https://example.net", "", oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: "test-token"})
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal("https://example.net"))
		})

		it("returns an API with an HTTPClient", func() {
			api, err := uaa.NewWithToken("https://example.net", "", oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: "test-token"})
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.UnauthenticatedClient).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
			Expect(reflect.TypeOf(api.AuthenticatedClient.Transport).String()).To(Equal("*uaa.tokenTransport"))
		})

		it("sets the authorization header correctly when round tripping", func() {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Header.Get("Authorization")).To(Equal("Bearer test-token"))
				w.WriteHeader(http.StatusOK)
			}))
			api, err := uaa.NewWithToken("https://example.net", "", oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: "test-token"})
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.UnauthenticatedClient).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
			r, err := api.AuthenticatedClient.Get(s.URL)
			Expect(err).NotTo(HaveOccurred())
			Expect(r.StatusCode).To(Equal(http.StatusOK))
		})

		it("Token() fails when the mode is token and the token is invalid", func() {
			api := uaa.New("https://example.net", "").WithToken(oauth2.Token{Expiry: time.Now().Add(-20 * time.Second), AccessToken: "test-token"})
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(t).To(BeNil())
		})

		it("Token() succeeds when the mode is token and the token is valid", func() {
			api := uaa.New("https://example.net", "").WithToken(oauth2.Token{Expiry: time.Now().Add(20 * time.Second), AccessToken: "test-token"})
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(t).NotTo(BeNil())
			Expect(t.Valid()).To(BeTrue())
		})
	})

	when("NewWithClientCredentials()", func() {
		it("fails if the target url is invalid", func() {
			api, err := uaa.NewWithClientCredentials("(*#&^@%$&%)", "", "", "", uaa.OpaqueToken, true)
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("returns an API with a TargetURL", func() {
			api, err := uaa.NewWithClientCredentials("https://example.net", "", "", "", uaa.OpaqueToken, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal("https://example.net"))
		})

		it("returns an API with an HTTPClient", func() {
			api, err := uaa.NewWithClientCredentials("https://example.net", "", "", "", uaa.OpaqueToken, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
		})

		it("Token() fails when the mode is client credentials and the client credentials are invalid", func() {
			api := uaa.New("(*#&^@%$&%)", "")
			Expect(api).NotTo(BeNil())
			api.TargetURL = nil
			api = api.WithClientCredentials("client-id", "client-secret", uaa.OpaqueToken)
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(t).To(BeNil())
		})

		when("the server returns tokens", func() {
			var (
				s           *httptest.Server
				returnToken bool
				callCount   int
			)

			it.Before(func() {
				returnToken = true
				s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					callCount = callCount + 1

					w.Header().Set("Content-Type", "application/json")

					t := &oauth2.Token{
						AccessToken:  "test-access-token",
						RefreshToken: "test-refresh-token",
						TokenType:    "bearer",
						Expiry:       time.Now().Add(60 * time.Second),
					}
					if !returnToken {
						t = nil
					}
					w.WriteHeader(http.StatusOK)
					err := json.NewEncoder(w).Encode(t)
					Expect(err).NotTo(HaveOccurred())
				}))
			})

			it.After(func() {
				if s != nil {
					s.Close()
				}
			})

			it("Token() succeeds when the mode is client credentials and the client credentials are valid", func() {
				api := uaa.New(s.URL, "")
				Expect(api).NotTo(BeNil())
				api.TargetURL = nil
				api = api.WithClientCredentials("client-id", "client-secret", uaa.OpaqueToken)
				Expect(api).NotTo(BeNil())
				t, err := api.Token(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(t).NotTo(BeNil())
				Expect(t.Valid()).To(BeTrue())
			})
		})
	})

	when("NewWithPasswordCredentials()", func() {
		it("fails if the target url is invalid", func() {
			api, err := uaa.NewWithPasswordCredentials("(*#&^@%$&%)", "", "", "", "", "", uaa.OpaqueToken, true)
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("returns an API with a TargetURL", func() {
			api, err := uaa.NewWithPasswordCredentials("https://example.net", "", "", "", "", "", uaa.OpaqueToken, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal("https://example.net"))
		})

		it("returns an API with an HTTPClient", func() {
			api, err := uaa.NewWithPasswordCredentials("https://example.net", "", "", "", "", "", uaa.OpaqueToken, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
		})
	})

	when("NewWithAuthorizationCode", func() {
		var (
			s           *httptest.Server
			returnToken bool
			reqBody     []byte
			callCount   int
		)

		it.Before(func() {
			returnToken = true
			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				callCount = callCount + 1
				var err error
				reqBody, err = ioutil.ReadAll(req.Body)
				Expect(err).NotTo(HaveOccurred())

				w.Header().Set("Content-Type", "application/json")

				t := &oauth2.Token{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(60 * time.Second),
				}
				if !returnToken {
					t = nil
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(t)
				Expect(err).NotTo(HaveOccurred())
			}))
		})

		it.After(func() {
			if s != nil {
				s.Close()
			}
		})

		it("fails if the target url is invalid", func() {
			api, err := uaa.NewWithAuthorizationCode("(*#&^@%$&%)", "", "", "", "", uaa.OpaqueToken, false)
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("returns an API with a TargetURL", func() {
			api, err := uaa.NewWithAuthorizationCode(s.URL, "", "", "", "", uaa.OpaqueToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal(s.URL))
			Expect(callCount).To(Equal(1))
		})

		it("returns an API with an HTTPClient", func() {
			api, err := uaa.NewWithAuthorizationCode(s.URL, "", "", "", "", uaa.OpaqueToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
		})

		it("returns an error if the token cannot be retrieved", func() {
			returnToken = false
			api, err := uaa.NewWithAuthorizationCode(s.URL, "", "", "", "", uaa.OpaqueToken, false)
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("ensure that auth code grant type params are set correctly", func() {
			api, err := uaa.NewWithAuthorizationCode(s.URL, "", "", "", "", uaa.OpaqueToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())

			Expect(string(reqBody)).To(ContainSubstring("token_format=opaque"))
			Expect(string(reqBody)).To(ContainSubstring("response_type=token"))
			Expect(string(reqBody)).To(ContainSubstring("grant_type=authorization_code"))
		})

		it("Token() fails when the mode is authorizationcode and the authorization code is invalid", func() {
			api := uaa.New("(*#&^@%$&%)", "")
			Expect(api).NotTo(BeNil())
			api.TargetURL = nil
			api = api.WithAuthorizationCode("client-id", "client-secret", "", uaa.OpaqueToken)
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(t).To(BeNil())
		})

		it("Token() will set the UnauthenticatedClient to the default if necessary", func() {
			api := uaa.New(s.URL, "")
			Expect(api).NotTo(BeNil())
			api.TargetURL = nil
			api = api.WithAuthorizationCode("client-id", "client-secret", "valid", uaa.OpaqueToken)
			Expect(api).NotTo(BeNil())
			api.UnauthenticatedClient = nil
			t, err := api.Token(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Valid()).To(BeTrue())
		})
	})

	when("NewWithRefreshToken", func() {
		var (
			s           *httptest.Server
			returnToken bool
			rawQuery    string
			reqBody     []byte
		)

		it.Before(func() {
			returnToken = true

			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				var err error
				rawQuery = req.URL.RawQuery
				reqBody, err = ioutil.ReadAll(req.Body)
				Expect(err).NotTo(HaveOccurred())

				w.Header().Set("Content-Type", "application/json")
				t := &oauth2.Token{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(60 * time.Second),
				}
				if !returnToken {
					t = nil
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(t)
				Expect(err).NotTo(HaveOccurred())
			}))
		})

		it.After(func() {
			if s != nil {
				s.Close()
			}
		})

		it("fails if the refresh token is invalid", func() {
			invalidRefreshToken := ""
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", invalidRefreshToken, uaa.JSONWebToken, false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("oauth2: token expired and refresh token is not set"))
			Expect(api).To(BeNil())
		})

		it("fails if the target url is invalid", func() {
			api, err := uaa.NewWithRefreshToken("(*#&^@%$&%)", "", "", "", "refresh-token", uaa.JSONWebToken, false)
			Expect(err).To(HaveOccurred())
			Expect(api).To(BeNil())
		})

		it("returns an API with a TargetURL", func() {
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", "refresh-token", uaa.JSONWebToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.TargetURL.String()).To(Equal(s.URL))
		})

		it("returns an API with an HTTPClient", func() {
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", "refresh-token", uaa.JSONWebToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())
			Expect(api.AuthenticatedClient).NotTo(BeNil())
		})

		it("returns an error if the token cannot be retrieved", func() {
			returnToken = false
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", "refresh-token", uaa.JSONWebToken, false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("oauth2: server response missing access_token"))
			Expect(api).To(BeNil())
		})

		it("ensure that refresh grant type params are set correctly", func() {
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", "refresh-token", uaa.JSONWebToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())

			Expect(rawQuery).To(Equal("token_format=jwt"))
			Expect(string(reqBody)).To(ContainSubstring("grant_type=refresh_token"))
			Expect(string(reqBody)).To(ContainSubstring("refresh_token=refresh-token"))
		})

		it("ensure that refresh grant type params are set correctly for opaque tokens", func() {
			api, err := uaa.NewWithRefreshToken(s.URL, "", "", "", "refresh-token", uaa.OpaqueToken, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(api).NotTo(BeNil())

			Expect(rawQuery).To(Equal("token_format=opaque"))
			Expect(string(reqBody)).To(ContainSubstring("grant_type=refresh_token"))
			Expect(string(reqBody)).To(ContainSubstring("refresh_token=refresh-token"))
		})

		it("Token() fails when the mode is refreshtoken and the refresh token is invalid", func() {
			api := uaa.New("(*#&^@%$&%)", "")
			Expect(api).NotTo(BeNil())
			api.TargetURL = nil
			api = api.WithRefreshToken("client-id", "client-secret", "", uaa.OpaqueToken)
			Expect(api).NotTo(BeNil())
			t, err := api.Token(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(t).To(BeNil())
		})

		it("Token() will set the UnauthenticatedClient to the default if necessary", func() {
			api := uaa.New(s.URL, "")
			Expect(api).NotTo(BeNil())
			api.TargetURL = nil
			api = api.WithRefreshToken("client-id", "client-secret", "valid", uaa.OpaqueToken)
			Expect(api).NotTo(BeNil())
			api.UnauthenticatedClient = nil
			t, err := api.Token(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Valid()).To(BeTrue())
		})
	})
}
