package help

import "fmt"

func Root(version string) string {
	return fmt.Sprintf(`UAA Command Line Interface, version %v

Feedback:
  Bugs or other issues can be filed on github.com/cloudfoundry/uaa-cli
  `, version)
}
