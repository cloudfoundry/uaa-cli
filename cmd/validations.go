package cmd

import (
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"fmt"
	"strings"
)

func avalableFormats() []string {
	return []string{"jwt", "opaque"}
}

func availableFormatsStr() string {
	return "[" + strings.Join(avalableFormats(), ", ") + "]"
}

func validateTokenFormatError(tokenFormat string) error {
	if !utils.Contains(avalableFormats(), tokenFormat) {
		return errors.New(fmt.Sprintf(`The token format "%v" is unknown. Available formats: %v`, tokenFormat, availableFormatsStr()))
	}
	return nil
}
