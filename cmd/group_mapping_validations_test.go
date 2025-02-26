package cmd_test

import (
	"code.cloudfoundry.org/uaa-cli/cmd"
	"code.cloudfoundry.org/uaa-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("validations", func() {
	Describe("without a target and context", func() {
		It("requires a target", func() {
			err := cmd.GroupMappingValidations(config.Config{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must set a target in order to use this command."))
		})

		It("requires a context", func() {
			cfg := config.NewConfigWithServerURL("http://localhost:9090")

			err := cmd.GroupMappingValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("You must have a token in your context to perform this command."))
		})
	})

	Describe("without required params", func() {
		It("requires a external_group_name", func() {
			cfg := buildConfig("http://localhost:9090")

			err := cmd.GroupMappingValidations(cfg, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
		})

		It("requires a group_name", func() {
			cfg := buildConfig("http://localhost:9090")

			err := cmd.GroupMappingValidations(cfg, []string{"external_group"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified."))
		})
	})

	Describe("with totally valid data", func() {
		It("does not complain", func() {
			cfg := buildConfig("http://localhost:9090")

			err := cmd.GroupMappingValidations(cfg, []string{"external_groupname", "groupname"})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
