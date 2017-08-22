package info_test

import (
	"github.com/jhamon/uaa/info"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Health", func() {
	var (
		server *ghttp.Server
		client info.UaaClient
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = info.UaaClient{server.URL()}
	})

	AfterEach(func() {
		server.Close()
	})

	It("calls the /healthz endpoint", func() {
		server.RouteToHandler("GET", "/healthz", ghttp.RespondWith(200, "ok"))
		server.AppendHandlers(ghttp.VerifyRequest("GET", "/healthz"))

		status := info.Health(client)

		Expect(status).To(Equal(info.OK))
	})

	It("returns error status when non-200 response", func() {
		server.RouteToHandler("GET", "/healthz", ghttp.RespondWith(400, "ok"))
		server.AppendHandlers(ghttp.VerifyRequest("GET", "/healthz"))

		status := info.Health(client)

		Expect(status).To(Equal(info.ERROR))
	})
})
