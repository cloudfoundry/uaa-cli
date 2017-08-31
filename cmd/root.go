package cmd

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/uaa"
)

var cfgFile uaa.Config

// User flags
var (
	password string
	username string
)

// Global flags
var (
	skipSSLValidation bool
	trace             bool
)

// Client flags
var (
	clientSecret         string
	authorizedGrantTypes string
	authorities          string
	accessTokenValidity  int32
	refreshTokenValidity int32
	displayName          string
	autoapprove          string
	scope                string
	redirectUri          string
	clone                string
	zoneSubdomain        string
)

var RootCmd = cobra.Command{
	Use:   "uaa",
	Short: "A cli for interacting with UAAs",
	Long:  help.Root(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", false, "See additional info on HTTP requests")
}

func initConfig() {
	// Startup tasks
}

var configLoaded bool
func GetSavedConfig() uaa.Config {
	if configLoaded {
		return cfgFile
	}

	cfgFile = config.ReadConfig()
	cfgFile.Trace = trace
	cfgFile.ZoneSubdomain = zoneSubdomain
	configLoaded = true
	return cfgFile
}
