package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/utils"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

type TargetStatus struct {
	Target            string
	Status            string
	UaaVersion        string
	SkipSSLValidation bool
}

func printTarget(log cli.Logger, target config.Target, status string, version string) error {
	return cli.NewJsonPrinter(log).Print(TargetStatus{target.BaseUrl, status, version, target.SkipSSLValidation})
}

func ShowTargetCmd(api *uaa.API, cfg config.Config, log cli.Logger) error {
	target := cfg.GetActiveTarget()

	if target.BaseUrl == "" {
		return printTarget(log, target, "", "")
	}

	info, err := api.GetInfo()
	if err != nil {
		_ = printTarget(log, target, "ERROR", "unknown")
		return errors.New("There was an error while fetching status info about the current target.")
	}

	return printTarget(log, target, "OK", info.App.Version)
}

func UpdateTargetCmd(cfg config.Config, newTarget string, log cli.Logger) error {
	target := config.Target{
		SkipSSLValidation: skipSSLValidation,
		BaseUrl:           newTarget,
	}
	cfg.AddTarget(target)

	api := GetUnauthenticatedAPIFromConfig(cfg)
	_, err := api.GetInfo()

	if err != nil {
		return errors.New(fmt.Sprintf("The target %s could not be set: %v", newTarget, err.Error()))
	}

	config.WriteConfig(cfg)
	log.Info("Target set to " + utils.Emphasize(newTarget))
	return nil
}

var targetCmd = &cobra.Command{
	Use:     "target UAA_URL",
	Aliases: []string{"api"},
	Short:   "Set the url of the UAA you'd like to target",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		if len(args) == 0 {
			NotifyErrorsWithRetry(ShowTargetCmd(GetUnauthenticatedAPI(), cfg, log), log)
		} else {
			NotifyErrorsWithRetry(UpdateTargetCmd(cfg, args[0], log), log)
		}
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVarP(&skipSSLValidation, "skip-ssl-validation", "k", false, "Disable security validation on requests to this target")
	targetCmd.Annotations = make(map[string]string)
	targetCmd.Annotations[INTRO_CATEGORY] = "true"
}
