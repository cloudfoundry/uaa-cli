package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

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
		log.Error("You must set a target in order to use this command.")
		os.Exit(1)
	}
}

func EnsureContext() {
	EnsureTarget()
	c := GetSavedConfig()

	if c.GetActiveTarget().ActiveContextName == "" {
		log.Error("You must have a token in your context to perform this command.")
		os.Exit(1)
	}
}
