package cmd

import (
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"os"
)

func printTarget(target uaa.Target, status string, version string) {
	log.Info("Target: " + target.BaseUrl)
	log.Info("Status: " + status)
	log.Info("UAA Version: " + version)
	log.Infof("SkipSSLValidation: %v", target.SkipSSLValidation)
}

func TraceRetryMsg(c uaa.Config) {
	if !c.Trace {
		log.Info("Retry with --trace for more information.")
	}
}

func showTarget() {
	c := GetSavedConfig()
	target := c.GetActiveTarget()

	if target.BaseUrl == "" {
		printTarget(target, "", "")
		return
	}

	info, err := uaa.Info(GetHttpClient(), c)
	if err != nil {
		printTarget(target, "ERROR", "unknown")
		os.Exit(1)
	}

	printTarget(target, "OK", info.App.Version)
}

func updateTarget(newTarget string) {
	c := GetSavedConfig()

	target := uaa.Target{
		SkipSSLValidation: skipSSLValidation,
		BaseUrl:           newTarget,
	}

	c.AddTarget(target)
	_, err := uaa.Info(GetHttpClientWithConfig(c), c)
	if err != nil {
		log.Errorf("The target %s could not be set.", newTarget)
		TraceRetryMsg(c)
		os.Exit(1)
	}

	config.WriteConfig(c)
	log.Info("Target set to " + newTarget)
}

var targetCmd = &cobra.Command{
	Use:   "target UAA_URL",
	Short: "Set the url of the UAA you'd like to target",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showTarget()
		} else {
			updateTarget(args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVarP(&skipSSLValidation, "skip-ssl-validation", "k", false, "Disable security validation on requests to this target")
}
