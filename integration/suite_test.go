//go:build integration

// Package integration exercises the compiled uaa-cli binary against a real,
// already-running UAA server (see scripts/start-uaa.sh) instead of a mocked
// HTTP server. It exists to catch bugs that only manifest against real
// server responses -- e.g. https://github.com/cloudfoundry-community/go-uaa
// unmarshalling a field UAA legitimately serializes as more than one JSON
// type -- which a mocked-response unit test can't surface because the mock
// only ever returns what the test author expects.
//
// Each resource area (clients, users, groups, tokens, misc) is its own
// top-level Describe in its own file, rather than one giant Ordered block.
// Ginkgo's Ordered containers skip every remaining spec in the container
// once one fails, so a single failure anywhere would otherwise hide the
// pass/fail status of everything unrelated to it. Splitting means a failure
// in "clients" still lets "users"/"groups"/"tokens"/"misc" report real
// results in the same run.
package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Live UAA Integration Suite")
}

var (
	commandPath string
	homeDir     string
	uaaTarget   string
)

var _ = BeforeSuite(func() {
	var err error
	commandPath, err = Build("code.cloudfoundry.org/uaa-cli")
	Expect(err).NotTo(HaveOccurred())

	homeDir, err = os.MkdirTemp("", "uaa-cli-integration")
	Expect(err).NotTo(HaveOccurred())
	Expect(os.Setenv("HOME", homeDir)).To(Succeed())

	uaaTarget = os.Getenv("UAA_TARGET")
	if uaaTarget == "" {
		uaaTarget = "http://localhost:8080/uaa"
	}

	waitForUAA(uaaTarget)

	// Every resource-area Describe below is independent and may run in any
	// order (Ginkgo randomizes top-level container order by default), so
	// targeting and authenticating as admin happens once here rather than as
	// a spec any one area could own.
	runOK("target", uaaTarget, "-k")
	runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
})

var _ = AfterSuite(func() {
	os.RemoveAll(homeDir)
	CleanupBuildArtifacts()
})

func waitForUAA(target string) {
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(target + "/login")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
	Fail(fmt.Sprintf("UAA at %s never became ready", target))
}

// run executes the compiled uaa binary and waits for it to exit, without
// asserting on the outcome -- for cleanup steps where a failure shouldn't
// fail the test itself.
func run(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 30).Should(Exit())
	return session
}

// runOK executes the compiled uaa binary and asserts it succeeded and didn't
// hit the class of unmarshal/response-parsing failure this suite guards
// against.
func runOK(args ...string) *Session {
	session := run(args...)
	stderr := string(session.Err.Contents())
	Expect(stderr).NotTo(ContainSubstring("cannot unmarshal"), "command %v produced an unmarshal error", args)
	Expect(stderr).NotTo(ContainSubstring("unknown error"), "command %v produced an unexpected error", args)
	Expect(session.ExitCode()).To(Equal(0), "command %v failed with stderr: %s", args, stderr)
	return session
}

func assertValidJSON(session *Session) {
	var v interface{}
	Expect(json.Unmarshal(session.Out.Contents(), &v)).To(Succeed(), "expected valid JSON, got: %s", string(session.Out.Contents()))
}

// groupID looks up a group's SCIM id by name, returning "" if it can't be
// found -- used for best-effort cleanup, not assertions.
func groupID(name string) string {
	session := run("get-group", name)
	if session.ExitCode() != 0 {
		return ""
	}
	var group struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(session.Out.Contents(), &group); err != nil {
		return ""
	}
	return group.ID
}
