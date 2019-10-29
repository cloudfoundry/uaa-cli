package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"errors"
)

func GroupMappingValidations(cfg config.Config, args []string) error {
	if err := cli.EnsureContextInConfig(cfg); err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("The positional arguments EXTERNAL_GROUPNAME and GROUPNAME must be specified.")
	}

	return nil
}
