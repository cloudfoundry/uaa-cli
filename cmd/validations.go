package cmd

import (
	"os"
	"code.cloudfoundry.org/uaa-cli/utils"
	"strings"
	"github.com/spf13/cobra"
)

func avalableFormats() []string {
	return []string { "jwt", "opaque" }
}

func availableFormatsStr() string {
	return "[" + strings.Join(avalableFormats(), ", ") + "]"
}

func validateTokenFormat(cmd *cobra.Command, tokenFormat string) {
	if !utils.Contains(avalableFormats(), tokenFormat) {
		log.Errorf(`The token format "%v" is unknown. Available formats: %v`, tokenFormat, availableFormatsStr())
		cmd.Usage()
		os.Exit(1)
	}
}
