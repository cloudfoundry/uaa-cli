package cmd_test

import (
	. "code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/fixtures"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("UnmapGroup", func() {
	buildConfig := func(target string) config.Config {
		cfg := config.NewConfigWithServerURL(target)
		ctx := config.NewContextWithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		cfg.AddContext(ctx)

		return cfg
	}

	mockGroupLookup := func(id, groupname string) {
		server.RouteToHandler("GET", "/Groups", CombineHandlers(
			VerifyRequest("GET", "/Groups", fmt.Sprintf("filter=displayName+eq+%%22%s%%22&startIndex=1&count=100", groupname)),
			RespondWith(http.StatusOK, fixtures.PaginatedResponse(uaa.Group{ID: id, DisplayName: groupname})),
		))
	}

	mockExternalGroupUnmapping := func(externalGroupname, internalGroupId, internalGroupname string) {
		origin := "ldap"
		path := fmt.Sprintf("/Groups/External/groupId/%s/externalGroup/%s/origin/%s", internalGroupId, externalGroupname, origin)
		server.RouteToHandler("DELETE", path, CombineHandlers(
			VerifyRequest("DELETE", path),
			RespondWith(http.StatusOK, fixtures.EntityResponse(
				uaa.GroupMapping{
					GroupID:       internalGroupId,
					ExternalGroup: externalGroupname,
					DisplayName:   internalGroupname,
					Origin:        origin,
					Schemas:       []string{"urn:scim:schemas:core:1.0"},
				})),
		))
	}

	Describe("by default", func() {
		BeforeEach(func() {
			config.WriteConfig(buildConfig(server.URL()))
		})

		It("Resolves the group name and performs the unmapping", func() {
			mockGroupLookup("internal-group-id", "internal-group")
			mockExternalGroupUnmapping("external-group", "internal-group-id", "internal-group")

			session := runCommand("unmap-group", "external-group", "internal-group")
			Eventually(session).Should(Say(`Successfully unmapped internal-group from external-group for origin ldap`))
			Eventually(session).Should(Exit(0))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Describe("validations", func() {
		Describe("without a target and context", func() {
			It("requires a target", func() {
				err := UnmapGroupValidations(config.Config{}, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
			})

			It("requires a context", func() {
				cfg := config.NewConfigWithServerURL("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
			})
		})

		Describe("without required params", func() {
			It("requires a external_group_name", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
			})

			It("requires a group_name", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{"external_group"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
			})
		})

		Describe("with totally valid data", func() {
			It("does not complain", func() {
				cfg := buildConfig("http://localhost:9090")

				err := UnmapGroupValidations(cfg, []string{"external_groupname", "groupname"})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

