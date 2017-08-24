package uaa_test

import (
	. "github.com/jhamon/uaa-cli/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Health", func() {
	var (
		server *ghttp.Server
		context UaaContext
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		context = UaaContext{}
		context.BaseUrl = server.URL()
	})

	AfterEach(func() {
		server.Close()
	})

	It("calls the /healthz endpoint", func() {
		server.RouteToHandler("GET", "/healthz", ghttp.RespondWith(200, "ok"))
		server.AppendHandlers(ghttp.VerifyRequest("GET", "/healthz"))

		status, _ := Health(context)

		Expect(status).To(Equal(OK))
	})

	It("returns error status when non-200 response", func() {
		server.RouteToHandler("GET", "/healthz", ghttp.RespondWith(400, "ok"))
		server.AppendHandlers(ghttp.VerifyRequest("GET", "/healthz"))

		status, _ := Health(context)

		Expect(status).To(Equal(ERROR))
	})
})
