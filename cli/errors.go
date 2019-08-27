package cli

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
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

func EnsureTargetInConfig(cfg config.Config) error {
	if cfg.ActiveTargetName == "" {
		return errors.New(MISSING_TARGET)
	}
	return nil
}

func EnsureContextInConfig(cfg config.Config) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if cfg.GetActiveTarget().ActiveContextName == "" {
		return errors.New(MISSING_CONTEXT)
	}
	return nil
}

func NotifyValidationErrors(err error, cmd *cobra.Command, log Logger) {
	if err != nil {
		log.Error(err.Error())
		cmd.Usage()
		os.Exit(1)
	}
}

func NotifyErrorsWithRetry(err error, log Logger, c config.Config) {
	if err != nil {
		switch t := err.(type) {
		case uaa.RequestError:
			log.Error(err.Error())
			NewJsonPrinter(log).PrintError(t.ErrorResponse)
		default:
			log.Error(err.Error())
		}
		verboseRetryMsg(log, c)
		os.Exit(1)
	}
}

func verboseRetryMsg(log Logger, c config.Config) {
	if !c.Verbose {
		log.Info("Retry with --verbose for more information.")
	}
}
