package cmd

import (
	"errors"
	"fmt"
	"os"
)

func MissingArgument(argName string) error {
	return errors.New(fmt.Sprintf("Missing argument `%v` must be specified.\n", argName))
}

func MissingArgumentForGrantType(argName, grantType string) error {
	return errors.New(fmt.Sprintf("Missing argument `%v` must be specified for %v grant type.\n", argName, grantType))
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
