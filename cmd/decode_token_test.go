package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

// testRS256JWT is a real RS256-signed JWT with payload:
//
//	{"sub":"abc123","iss":"https://uaa.example.com/oauth/token",
//	 "iat":1505079823,"exp":1507671823,"jti":"test-jti"}
//
// Signed with testRS256PublicKey's corresponding private key.
const testRS256JWT = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9" +
	".eyJzdWIiOiJhYmMxMjMiLCJpc3MiOiJodHRwczovL3VhYS5leGFtcGxlLmNvbS9vYXV0aC90b2tlbiIsImlhdCI6MTUwNTA3OTgyMywiZXhwIjoxNTA3NjcxODIzLCJqdGkiOiJ0ZXN0LWp0aSJ9" +
	".hE_Re8Exrxe86wNfzYeDCXohy-d-QaqUCS3WoZ5wtUs7GSRbXeZubwT62MxO2FXl2iVx3NYphsaJ7P7IABvMP6UsPEYZ3oOdgGyXOBIki6GqnHiMueE5hwNKyETsmovRJmVo7PHPceYKM1J2RugQan1np8ELMMLGgWJQAjOD4TqOUPQC6CYfIZcaVATClV_lXFYun7_6hox6N7QooEB27-YquZYK88gSspvRC9m19VnuK4UuGqa0VpPHBfde_k1aPl-cabj-kvcywYXnwY4pZYwJas3Prvm8cxEQK49V7Bemf89qwNHhn-0VEXtPcwLn3Xcc7DWM0fFyN1qx2MepPg"

const testRS256PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzj9b79BofXFY7pcLrDjv
njB/vHHbLwc33sAa0eiPu3n6tgi7z7k/VpXt7s5/CiGhcJt1oHb7b8zpS0Vg96Is
p4KbPIwA+GTAYqEPDLLIzxPxB13lVUXMUmhqc2oY/ryRpYs87WywT4maaaabwV1H
2v2LWUXnNO7zeFz7amXAwLjpombEEGZbA6AGt21PWnLfhRxyb/+Vk0h1enF+SSH8
B8nG6D55RaNCgQ/GduTll9cQu2GHV2npmyEagr80dGNS0fvprWf9zZjzfsXZUH+Z
5PabXmqiYnHPAiRkPgv1BkLeYt72t4IykRB7Cws6h7ltP8zY89x2YgxY+OtSBX8q
AQIDAQAB
-----END PUBLIC KEY-----`

const testRS256WrongPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2a2rwplBQLF29amygykE
MmYz0+Kcj3bKBp29P2rFj7bQMBqBaMhBNjRnSFBqHHDEIxJBBGJFCFBJCBBBBBB
BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBw
IDAQAB
-----END PUBLIC KEY-----`

var _ = Describe("DecodeToken", func() {
	Describe("when no token is in context and no arg is given", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 1 with a helpful error", func() {
			session := runCommand("decode-token")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("no token"))
		})
	})

	Describe("when a token is in the active context", func() {
		BeforeEach(func() {
			c := config.NewConfigWithServerURL(server.URL())
			ctx := config.NewContextWithToken(testRS256JWT)
			c.AddContext(ctx)
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and prints decoded claims as JSON", func() {
			session := runCommand("decode-token")

			Eventually(session).Should(Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring(`"sub"`))
			Expect(output).To(ContainSubstring(`"abc123"`))
			Expect(output).To(ContainSubstring(`"iss"`))
			Expect(output).To(ContainSubstring(`uaa.example.com`))
		})

		It("exits 0 with --verbose flag", func() {
			session := runCommand("decode-token", "--verbose")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(`"sub"`))
		})
	})

	Describe("when a token is passed as a positional argument", func() {
		BeforeEach(func() {
			c := config.NewConfig()
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and prints decoded claims", func() {
			session := runCommand("decode-token", testRS256JWT)

			Eventually(session).Should(Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring(`"sub"`))
			Expect(output).To(ContainSubstring(`"abc123"`))
			Expect(output).To(ContainSubstring(`"jti"`))
			Expect(output).To(ContainSubstring(`"test-jti"`))
		})

		It("exits 0 when a token type is also passed as a second argument", func() {
			session := runCommand("decode-token", testRS256JWT, "bearer")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(`"sub"`))
		})
	})

	Describe("--key flag", func() {
		BeforeEach(func() {
			c := config.NewConfig()
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and prints 'Valid token signature' when key matches", func() {
			session := runCommand("decode-token", testRS256JWT, "--key", testRS256PublicKey)

			Eventually(session).Should(Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("Valid token signature"))
			Expect(output).To(ContainSubstring(`"sub"`))
		})

		It("exits 1 with a signature error when key does not match", func() {
			session := runCommand("decode-token", testRS256JWT, "--key", testRS256WrongPublicKey)

			Eventually(session).Should(Exit(1))
		})
	})

	Describe("--decode-times flag", func() {
		BeforeEach(func() {
			c := config.NewConfig()
			Expect(config.WriteConfig(c)).Error().ShouldNot(HaveOccurred())
		})

		It("exits 0 and prints a human-readable timestamp section", func() {
			session := runCommand("decode-token", testRS256JWT, "--decode-times")

			Eventually(session).Should(Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("Decoded timestamps"))
			Expect(output).To(ContainSubstring("iat"))
			Expect(output).To(ContainSubstring("Issued At"))
			Expect(output).To(ContainSubstring("exp"))
			Expect(output).To(ContainSubstring("Expires At"))
		})

		It("does not print timestamp section without the flag", func() {
			session := runCommand("decode-token", testRS256JWT)

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("Decoded timestamps"))
		})
	})
})
