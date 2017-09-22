package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"

	"code.cloudfoundry.org/uaa-cli/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var (
	commandPath string
	homeDir     string
	server      *Server
)

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "uaa-test")
	Expect(err).NotTo(HaveOccurred())
	server = NewServer()

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", homeDir)
	} else {
		os.Setenv("HOME", homeDir)
	}
})

var _ = AfterEach(func() {
	os.RemoveAll(homeDir)
	config.RemoveConfig()
	server.Close()
})

var _ = SynchronizedBeforeSuite(func() []byte {
	executable_path, err := Build("code.cloudfoundry.org/uaa-cli", "-ldflags", "-X code.cloudfoundry.org/uaa-cli/version.Version=test-version")
	Expect(err).NotTo(HaveOccurred())
	return []byte(executable_path)
}, func(data []byte) {
	commandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func runCommandWithEnv(env []string, args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	existing := os.Environ()
	for _, env_var := range env {
		existing = append(existing, env_var)
	}
	cmd.Env = existing
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func runCommandWithStdin(stdin io.Reader, args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	cmd.Stdin = stdin
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func ItBehavesLikeHelp(command string, alias string, validate func(*Session)) {
	It("displays help", func() {
		session := runCommand(command, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})

	It("displays help using the alias", func() {
		session := runCommand(alias, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})
}

func ItSupportsTheVerboseFlagWhenGet(command string, endpoint string, responseJson string) {
	It("shows extra output about the request on success", func() {
		server.RouteToHandler("GET", endpoint,
			RespondWith(http.StatusOK, responseJson),
		)

		session := runCommand(command, "--verbose")

		Eventually(session).Should(Exit(0))
		Expect(session.Out).To(Say("GET " + server.URL() + endpoint))
		Expect(session.Out).To(Say("Accept: application/json"))
		Expect(session.Out).To(Say("200 OK"))
	})

	It("shows extra output about the request on error", func() {
		server.RouteToHandler("GET", endpoint,
			RespondWith(http.StatusBadRequest, "garbage response"),
		)

		session := runCommand(command, "--verbose")

		Eventually(session).Should(Exit(1))
		Expect(session.Out).To(Say("GET " + server.URL() + endpoint))
		Expect(session.Out).To(Say("Accept: application/json"))
		Expect(session.Out).To(Say("400 Bad Request"))
		Expect(session.Out).To(Say("garbage response"))
	})
}
