package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"errors"
	"fmt"
	"net/http"
	"code.cloudfoundry.org/uaa-cli/utils"
	"code.cloudfoundry.org/uaa-cli/cli"
)

type TargetStatus struct {
	Target string
	Status string
	UaaVersion string
	SkipSSLValidation bool
}

func printTarget(log utils.Logger, target uaa.Target, status string, version string) error {
	return cli.NewJsonPrinter(log).Print(TargetStatus{target.BaseUrl, status, version, target.SkipSSLValidation})
}

func ShowTargetCmd(cfg uaa.Config, httpClient *http.Client, log utils.Logger) error {
	target := cfg.GetActiveTarget()

	if target.BaseUrl == "" {
		return printTarget(log, target, "", "")
	}

	info, err := uaa.Info(httpClient, cfg)
	if err != nil {
		_ = printTarget(log, target, "ERROR", "unknown")
		return errors.New("There was an error while fetching status info about the current target.")
	}

	return printTarget(log, target, "OK", info.App.Version)
}

func UpdateTargetCmd(cfg uaa.Config, newTarget string, log utils.Logger) error {
	target := uaa.Target{
		SkipSSLValidation: skipSSLValidation,
		BaseUrl:           newTarget,
	}

	cfg.AddTarget(target)
	_, err := uaa.Info(GetHttpClientWithConfig(cfg), cfg)
	if err != nil {
		return errors.New(fmt.Sprintf("The target %s could not be set.", newTarget))
	}

	config.WriteConfig(cfg)
	log.Info("Target set to " + newTarget)
	return nil
}

var targetCmd = &cobra.Command{
	Use:   "target UAA_URL",
	Short: "Set the url of the UAA you'd like to target",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		if len(args) == 0 {
			NotifyErrorsWithRetry(ShowTargetCmd(cfg, GetHttpClient(), log), cfg, log)
		} else {
			NotifyErrorsWithRetry(UpdateTargetCmd(cfg, args[0], log), cfg, log)
		}
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVarP(&skipSSLValidation, "skip-ssl-validation", "k", false, "Disable security validation on requests to this target")
	targetCmd.Annotations = make(map[string]string)
	targetCmd.Annotations[INTRO_CATEGORY] = "true"
}
