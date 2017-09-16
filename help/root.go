package help

import "fmt"

func Root(version string) string {
	return fmt.Sprintf(`UAA Command Line Interface, version %v

Feedback:
  Email cf-identity-eng@pivotal.io with your thoughts on the experience of using this
  tool. Bugs or other issues can be filed on github.com/cloudfoundry-incubator/uaa-cli
  `, version)
}
