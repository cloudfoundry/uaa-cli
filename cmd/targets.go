package cmd

import (
	"fmt"
	"sort"

	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/spf13/cobra"
)

func ListTargetsCmd(cfg config.Config) {
	if len(cfg.Targets) == 0 {
		log.Info("No targets set.")
		return
	}

	keys := make([]string, 0, len(cfg.Targets))
	for k := range cfg.Targets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		target := cfg.Targets[k]
		marker := "  "
		if k == cfg.ActiveTargetName {
			marker = "* "
		}
		log.Info(fmt.Sprintf("%s%d: %s", marker, i+1, target.BaseUrl))
	}
}

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List all registered targets",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		ListTargetsCmd(cfg)
	},
}

func init() {
	RootCmd.AddCommand(targetsCmd)
	targetsCmd.Annotations = make(map[string]string)
	targetsCmd.Annotations[INTRO_CATEGORY] = "true"
}
