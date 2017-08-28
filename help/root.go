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
  Email jhamon@pivotal.io with your thoughts on the experience of using this
  tool. Bugs or other issues can be filed on github.com/jhamon/uaa-cli

  `
}
