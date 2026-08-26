//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("clients", Ordered, func() {
	clientID := "uaa-cli-integration-test-client"

	AfterAll(func() {
		// Other Describes assume an admin token is active; assert this one.
		runOK("get-client-credentials-token", "admin", "-s", "adminsecret")
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
		// client_federated_jwt_trust has no allowpublic key at all, so unlike
		// list-clients above this doesn't hit the go-uaa bug -- it's here to
		// confirm get-client for other legacy clients still works today.
		assertValidJSON(runOK("get-client", "client_federated_jwt_trust"))
	})

	It("creates a client", func() {
		runOK("create-client", clientID,
			"-s", "test-secret",
			"--authorized_grant_types", "client_credentials,refresh_token",
			"--scope", "uaa.none",
			// clients.secret is required for a client to change its own secret
			// via change-client-secret below.
			"--authorities", "uaa.none,clients.secret",
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
