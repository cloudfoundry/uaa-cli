package cmd

import (
	"errors"
	"fmt"
	"os"
)

func MissingArgument(argName string) error {
	return errors.New(fmt.Sprintf("Missing argument `%v` must be specified.\n", argName))
}

func EnsureTarget() {
	c := GetSavedConfig()

	if c.ActiveTargetName == "" {
		fmt.Println("You must set a target in order to use this command.")
		os.Exit(1)
	}
}

func EnsureContext() {
	EnsureTarget()
	c := GetSavedConfig()

	if c.GetActiveTarget().ActiveContextName == "" {
		fmt.Println("You must set have a token in your context to perform this command.")
		os.Exit(1)
	}
}