package uaa_test

import (
	"encoding/json"
	"net/http"

	. "code.cloudfoundry.org/uaa-cli/fixtures"
	. "code.cloudfoundry.org/uaa-cli/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Curl", func() {
	var (
		cm        CurlManager
		uaaServer *ghttp.Server
	)

	BeforeEach(func() {
		uaaServer = ghttp.NewServer()
		config := NewConfigWithServerURL(uaaServer.URL())
		config.AddContext(NewContextWithToken("access_token"))
		cm = CurlManager{&http.Client{}, config}
	})

	Describe("CurlManager#Curl", func() {
		It("gets a user from the UAA by id", func() {
			uaaServer.RouteToHandler("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, MarcusUserResponse),
			))

			_, resBody, err := cm.Curl("/Users/fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70", "GET", "", []string{"Accept: application/json"})
			Expect(err).NotTo(HaveOccurred())

			var user ScimUser
			err = json.Unmarshal([]byte(resBody), &user)
			Expect(err).NotTo(HaveOccurred())

			Expect(user.Id).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
		})

		It("can POST body and multiple headers", func() {
			reqBody := map[string]interface{}{
				"externalId": "marcus-user",
				"userName":   "marcus@stoicism.com",
			}
			uaaServer.RouteToHandler("POST", "/Users", ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/Users"),
				ghttp.VerifyHeaderKV("Authorization", "bearer access_token"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSONRepresenting(reqBody),
				ghttp.RespondWith(http.StatusCreated, MarcusUserResponse),
			))

			reqBodyBytes, err := json.Marshal(reqBody)
			Expect(err).NotTo(HaveOccurred())

			_, resBody, err := cm.Curl("/Users", "POST", string(reqBodyBytes), []string{"Content-Type: application/json", "Accept: application/json"})

			var user ScimUser
			err = json.Unmarshal([]byte(resBody), &user)
			Expect(err).NotTo(HaveOccurred())

			Expect(user.Id).To(Equal("fb5f32e1-5cb3-49e6-93df-6df9c8c8bd70"))
		})
	})

})
