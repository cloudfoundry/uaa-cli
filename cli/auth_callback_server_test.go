package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/utils"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func serverUrl(port int) string {
	return fmt.Sprintf("http://localhost:%v/", port)
}

var _ = Describe("AuthCallbackServer", func() {
	var (
		httpClient *http.Client
		acs        AuthCallbackServer
		done       chan string
		randPort   int
		logger     utils.Logger
	)

	BeforeSuite(func() {
		rand.Seed(time.Now().Unix())
	})

	BeforeEach(func() {
		randPort = rand.Intn(50000-8000) + 8000

		httpClient = &http.Client{}
		logger := utils.NewLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
		acs = NewAuthCallbackServer(
			"<h1>Hello There</h1>",
			"<style> h1 { background: #F00 } </style>",
			"<script>console.log('Objective judgement, now, at this very moment.')</script>",
			logger,
			randPort)

		done = make(chan string)
	})

	AfterEach(func() {
		close(done)
	})

	Describe("Start", func() {
		It("serves static content on the configured port", func() {
			acs.Hangup = func(donedone chan string, values url.Values) {
				donedone <- "Yep, exit immediately after first call"
			}

			go acs.Start(done)

			resp, err := httpClient.Get(serverUrl(randPort))
			if err != nil {
				Fail(err.Error())
			}

			<-done
			parsedBody, _ := ioutil.ReadAll(resp.Body)
			Expect(string(parsedBody)).To(ContainSubstring("Hello There"))
			Expect(string(parsedBody)).To(ContainSubstring("background: #F00"))
			Expect(string(parsedBody)).To(ContainSubstring("Objective judgement"))
		})
	})

	It("uses the Hangup func to decide when to close the server", func() {
		acs.Hangup = func(donedone chan string, values url.Values) {
			if values.Get("code") != "" {
				donedone <- values.Get("code")
			}
		}

		go acs.Start(done)

		_, err := httpClient.Get(serverUrl(randPort))
		Expect(err).NotTo(HaveOccurred())
		_, err = httpClient.Get(serverUrl(randPort) + "?foo=not_the_code")
		Expect(err).NotTo(HaveOccurred())
		_, err = httpClient.Get(serverUrl(randPort) + "?code=secret_code")
		Expect(err).NotTo(HaveOccurred())

		// Server should close after first request with "code" param
		// Sleep so we don't call while the server is still closing
		time.Sleep(20 * time.Millisecond)
		_, err = httpClient.Get(serverUrl(randPort))
		Expect(err).To(HaveOccurred())

		Expect(<-done).To(Equal("secret_code"))
	})
})
