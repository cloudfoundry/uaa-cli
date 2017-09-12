package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"code.cloudfoundry.org/uaa-cli/utils"
)

const MISSING_TARGET = "You must set a target in order to use this command."
const MISSING_CONTEXT = "You must have a token in your context to perform this command."

func MissingArgumentWithExplanation(argName string, cmd *cobra.Command, explanation string) {
	log.Errorf("Missing argument `%v` must be specified. %v", argName, explanation)
	cmd.Usage()
	os.Exit(1)
}

func MissingArgument(argName string, cmd *cobra.Command) {
	MissingArgumentWithExplanation(argName, cmd, "")
}

func MissingArgumentForGrantType(argName, grantType string, cmd *cobra.Command) {
	log.Errorf("Missing argument `%v` must be specified for %v grant type.", argName, grantType)
	cmd.Usage()
	os.Exit(1)
}

func EnsureTarget() {
	c := GetSavedConfig()

	if c.ActiveTargetName == "" {
		log.Error(MISSING_TARGET)
		os.Exit(1)
	}
}

func EnsureContext() {
	EnsureTarget()
	c := GetSavedConfig()

	if c.GetActiveTarget().ActiveContextName == "" {
		log.Error(MISSING_CONTEXT)
		os.Exit(1)
	}
}

func EnsureTargetInConfig(cfg uaa.Config) error {
	if cfg.ActiveTargetName == "" {
		return errors.New(MISSING_TARGET)
	}
	return nil
}

func EnsureContextInConfig(cfg uaa.Config) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if cfg.GetActiveTarget().ActiveContextName == "" {
		return errors.New(MISSING_CONTEXT)
	}
	return nil
}

func NotifyValidationErrors(err error, cmd *cobra.Command, log utils.Logger) {
	if err != nil {
		log.Error(err.Error())
		cmd.Usage()
		os.Exit(1)
	}
}

func NotifyErrorsWithRetry(err error, cfg uaa.Config, log utils.Logger) {
	if err != nil {
		log.Error(err.Error())
		TraceRetryMsg(GetSavedConfig())
		os.Exit(1)
	}
}
