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
