package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const MISSING_TARGET = "You must set a target in order to use this command."
const MISSING_CONTEXT = "You must have a token in your context to perform this command."

func MissingArgumentError(argName string) error {
	return MissingArgumentWithExplanationError(argName, "")
}

func MissingArgumentWithExplanationError(argName string, explanation string) error {
	return errors.New(fmt.Sprintf("Missing argument `%v` must be specified. %v", argName, explanation))
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

func NotifyValidationErrors(err error, cmd *cobra.Command, log cli.Logger) {
	if err != nil {
		log.Error(err.Error())
		cmd.Usage()
		os.Exit(1)
	}
}

func NotifyErrorsWithRetry(err error, cfg uaa.Config, log cli.Logger) {
	if err != nil {
		log.Error(err.Error())
		TraceRetryMsg(GetSavedConfig())
		os.Exit(1)
	}
}

func TraceRetryMsg(c uaa.Config) {
	if !c.Trace {
		log.Info("Retry with --trace for more information.")
	}
}
