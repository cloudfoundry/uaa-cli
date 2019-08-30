module code.cloudfoundry.org/uaa-cli

go 1.12

require (
	github.com/cloudfoundry-community/go-uaa v0.2.5
	github.com/fatih/color v1.7.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.2 // indirect
	github.com/olekukonko/tablewriter v0.0.0-20170719101040-be5337e7b39e
	github.com/onsi/ginkgo v1.10.0
	github.com/onsi/gomega v1.6.0
	github.com/pkg/errors v0.8.1
	github.com/skratchdot/open-golang v0.0.0-20160302144031-75fb7ed4208c
	github.com/spf13/cobra v0.0.0-20170820023359-4a7b7e65864c
	github.com/spf13/pflag v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

//replace github.com/cloudfoundry-community/go-uaa => /Users/pivotal/workspace/team-forks/go-uaa
