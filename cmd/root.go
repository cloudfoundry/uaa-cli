package cmd

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/version"
)

var cfgFile uaa.Config
var log cli.Logger

// Token flags
var (
	password    string
	username    string
	tokenFormat string
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
	accessTokenValidity  int64
	refreshTokenValidity int64
	displayName          string
	scope                string
	redirectUri          string
	clone                string
	zoneSubdomain        string
	port                 int
)

// Users
var (
	origin string
)

// SCIM filters
var (
	filter     string
	sortBy     string
	sortOrder  string
	attributes string
	count      int
	startIndex int
)

// Command categories
const (
	INTRO_CATEGORY       = "Getting Started"
	TOKEN_CATEGORY       = "Getting Tokens"
	CLIENT_CRUD_CATEGORY = "Managing Clients"
	USER_CRUD_CATEGORY   = "Managing Users"
	MISC_CATEGORY        = "Miscellaneous"
)

var RootCmd = cobra.Command{
	Use:   "uaa",
	Short: "A cli for interacting with UAAs",
	Long:  help.Root(version.VersionString()),
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
	RootCmd.Annotations = make(map[string]string)
	RootCmd.Annotations[INTRO_CATEGORY] = "true"
	RootCmd.Annotations[TOKEN_CATEGORY] = "true"
	RootCmd.Annotations[CLIENT_CRUD_CATEGORY] = "true"
	RootCmd.Annotations[USER_CRUD_CATEGORY] = "true"
	RootCmd.Annotations[MISC_CATEGORY] = "true"
	RootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{ $cmds := .Commands }}{{ range $category, $value := .Annotations }}

{{ $category }}:  {{range $cmd := $cmds}}{{if (or (and $cmd.IsAvailableCommand (eq (index $cmd.Annotations $category) "true")) (and (eq $cmd.Name "help") (eq $category "Getting Started")))}}
  {{rpad $cmd.Name $cmd.NamePadding }}  {{$cmd.Short}}{{end}}{{end}}{{ end }}{{ end }}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	log = cli.NewLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func initConfig() {
	// Startup tasks
}

func GetLogger() *cli.Logger {
	return &log
}

func GetSavedConfig() uaa.Config {
	cfgFile = config.ReadConfig()
	cfgFile.Trace = trace
	cfgFile.ZoneSubdomain = zoneSubdomain
	return cfgFile
}
