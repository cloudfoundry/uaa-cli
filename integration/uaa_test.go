//go:build integration

// Package integration exercises the compiled uaa-cli binary against a real,
// already-running UAA server (see scripts/start-uaa.sh) instead of a mocked
// HTTP server. It exists to catch bugs that only manifest against real
// server responses -- e.g. https://github.com/cloudfoundry-community/go-uaa
// unmarshalling a field UAA legitimately serializes as more than one JSON
// type -- which a mocked-response unit test can't surface because the mock
// only ever returns what the test author expects.
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
})

var _ = AfterSuite(func() {
	os.RemoveAll(homeDir)
	CleanupBuildArtifacts()
})

func waitForUAA(target string) {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(target + "/login")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
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

var _ = Describe("uaa-cli against a live UAA server", Ordered, func() {
	It("targets the server", func() {
		runOK("target", uaaTarget, "-k")
	})

	It("gets an admin client-credentials token", func() {
		runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
	})

	It("shows server info", func() {
		assertValidJSON(runOK("info"))
	})

	It("shows the current context", func() {
		runOK("context")
	})

	It("lists saved contexts", func() {
		runOK("contexts")
	})

	It("prints the version", func() {
		runOK("version")
	})

	Describe("clients", func() {
		clientID := "uaa-cli-integration-test-client"

		AfterAll(func() {
			run("get-client-credentials-token", "admin", "-s", "adminsecret")
			run("delete-client", clientID)
		})

		It("lists clients, including the legacy string-typed additionalInformation clients seeded by scripts/boot/uaa.yml", func() {
			// login, client_federated_jwt_trust, client_with_allowpublic_and_jwks_uri_trust,
			// and oauth_showcase_saml2_bearer are all seeded with allowpublic/autoapprove as
			// YAML strings, which UAA stores verbatim and serializes back as JSON strings
			// rather than booleans -- the exact shape that broke `uaa list-clients` originally.
			session := runOK("list-clients")
			assertValidJSON(session)
			Expect(string(session.Out.Contents())).To(ContainSubstring("client_federated_jwt_trust"))
		})

		It("gets one of the legacy string-typed clients directly", func() {
			assertValidJSON(runOK("get-client", "client_federated_jwt_trust"))
		})

		It("creates a client", func() {
			runOK("create-client", clientID,
				"-s", "test-secret",
				"--authorized_grant_types", "client_credentials,refresh_token",
				"--scope", "uaa.none",
				"--authorities", "uaa.none",
			)
		})

		It("gets the new client", func() {
			assertValidJSON(runOK("get-client", clientID))
		})

		It("updates the client", func() {
			runOK("update-client", clientID, "--scope", "uaa.none,openid")
		})

		It("sets the client secret", func() {
			runOK("set-client-secret", clientID, "-s", "second-test-secret")
		})

		It("gets a token as the client and changes its own secret", func() {
			runOK("get-client-credentials-token", clientID, "-s", "second-test-secret")
			runOK("change-client-secret", "--old_secret", "second-test-secret", "--secret", "third-test-secret")
		})
	})

	Describe("users", func() {
		username := "uaa-cli-integration-test-user"

		BeforeAll(func() {
			runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
		})

		AfterAll(func() {
			run("delete-user", username)
		})

		It("creates a user", func() {
			runOK("create-user", username,
				"--givenName", "Integration",
				"--familyName", "Test",
				"--email", username+"@example.com",
				"-p", "S0meSecur3Pass!",
			)
		})

		It("gets the user", func() {
			assertValidJSON(runOK("get-user", username))
		})

		It("lists users with a filter", func() {
			assertValidJSON(runOK("list-users", "--filter", fmt.Sprintf(`userName eq %q`, username)))
		})

		It("updates the user", func() {
			runOK("update-user", username, "--givenName", "Updated")
		})

		It("deactivates and reactivates the user", func() {
			runOK("deactivate-user", username)
			runOK("activate-user", username)
		})

		It("unlocks the user", func() {
			runOK("unlock-user", username)
		})
	})

	Describe("groups", func() {
		groupName := "uaa.cli.integration.test.group"
		memberUsername := "uaa-cli-integration-test-group-member"

		BeforeAll(func() {
			runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
			runOK("create-user", memberUsername,
				"--givenName", "Group", "--familyName", "Member",
				"--email", memberUsername+"@example.com", "-p", "S0meSecur3Pass!",
			)
		})

		AfterAll(func() {
			run("delete-user", memberUsername)
		})

		It("creates a group", func() {
			runOK("create-group", groupName, "-d", "uaa-cli integration test group")
		})

		It("gets the group", func() {
			assertValidJSON(runOK("get-group", groupName))
		})

		It("lists groups", func() {
			assertValidJSON(runOK("list-groups"))
		})

		It("adds and removes a member", func() {
			runOK("add-member", groupName, memberUsername)
			runOK("remove-member", groupName, memberUsername)
		})

		It("maps and unmaps an external group", func() {
			runOK("map-group", "cn=integration-test,ou=groups,dc=example,dc=com", groupName)
			assertValidJSON(runOK("list-group-mappings"))
			runOK("unmap-group", "cn=integration-test,ou=groups,dc=example,dc=com", groupName)
		})
	})

	Describe("tokens", func() {
		BeforeEach(func() {
			runOK("get-password-token", "cf", "-s", "", "-u", "marissa", "-p", "koala")
		})

		AfterEach(func() {
			run("get-client-credentials-token", "admin", "-s", "adminsecret")
		})

		It("decodes the current access token", func() {
			token := string(runOK("context", "--access_token").Out.Contents())
			runOK("decode-token", token)
		})

		It("gets the token signing key(s)", func() {
			runOK("get-token-key")
			runOK("get-token-keys")
		})

		It("refreshes the token", func() {
			runOK("refresh-token", "-s", "")
		})
	})

	Describe("misc", func() {
		It("curls an arbitrary endpoint", func() {
			assertValidJSON(runOK("curl", "/info"))
		})

		It("gets userinfo for a user-context token", func() {
			runOK("get-password-token", "cf", "-s", "", "-u", "marissa", "-p", "koala")
			assertValidJSON(runOK("userinfo"))
			runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
		})
	})
})
