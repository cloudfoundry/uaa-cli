package help

func Root() string {
	return `      __  _____   ___
     / / / / _ | / _ |  Universal Authentication and Authorization
    / /_/ / __ |/ __ |            Command-Line Interface
    \____/_/ |_/_/ |_|

This command-line interface has been designed to help app developers and
platform operators troubleshoot and administer their deployed UAAs.

This tool may be used to:

  * Discover UAA status, token signing keys, and other metadata
  * Create, test, and validate OAuth2 client configurations
  * Manage UAA resources such as users, groups, and memberships

Feedback:
  Email cf-identity-eng@pivotal.io with your thoughts on the experience of using this
  tool. Bugs or other issues can be filed on github.com/cloudfoundry-incubator/uaa-cli

Usage:
  uaa [command]

Getting Started:
  help                         Help about any command
  target                       Set the url of the UAA you'd like to target
  context                      See information about the currently active CLI context
  info                         See version and global configurations for the targeted UAA

Get Tokens:
  get-client-credentials-token Obtain a token using the client_credentials grant type
  get-password-token           Obtain a token using the password grant type

Client Management:
  create-client                Create an OAuth client registration in the UAA
  get-client                   View client registration
  update-client                Update an OAuth client registration in the UAA
  delete-client                Delete a client registration

User Management:
  me                           See claims about the authenticated user

Flags:
  -h, --help    help for uaa
  -t, --trace   See additional info on HTTP requests

Use "uaa [command] --help" for more information about a command.

  `
}
