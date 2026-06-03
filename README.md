# UAA Command Line Interface

CLI for [UAA](https://github.com/cloudfoundry/uaa) written in golang. This is an alterntive to using uaac which is wirtten in Ruby.  At this time it performs a limited subset of the features provided by the [uaac](https://github.com/cloudfoundry/cf-uaac) gem.  The team plans to continue development on the golang CLI going forward, and once it's considered fully GA, intends to place it alongside uaac with a long-term intention of one day deprecating uaac.

### Command Reference

See the [Command Reference](docs/commands.md) for the full list of commands and their options.

> **Migrating from uaac?** See the [Migrating from uaac](docs/migrating-from-uaac.md) guide for a side-by-side command reference.

### Goals

- To provide a CLI which can be easily installed in environments without a functioning Ruby setup
- To more closely conform to the style of other widely used CLIs in the CF ecosystem, e.g. the cf CLI. Commands should be of the form VERB-NOUN, similar to `cf delete-app`.
- To provide outputs that are machine-parseable whenever possible.
- To improve the quality of help strings and error messages so that users can self-diagnose problems and unblock themselves.
- To provide only the essential, highly used and/or required command options.

## Installation

### Install with Go

```
go install code.cloudfoundry.org/uaa-cli@latest
uaa -h
```

### Install with Homebrew

Install using the Homebrew cask from the [cloudfoundry tap](https://github.com/cloudfoundry/homebrew-tap):

```bash
# Install the cask
brew install --cask cloudfoundry/tap/uaa-cli
```

If upgrading from the old formula-based version, uninstall it first
```bash
brew uninstall --force uaa-cli
```

### Building from source

```bash
git clone https://github.com/cloudfoundry/uaa-cli.git
cd uaa-cli
make && make install
uaa -h
```

### Troubleshooting Installation

If you encounter trust warnings with Homebrew, you can either:

1. Trust the specific tap (recommended):
   ```bash
   brew trust cloudfoundry/tap
   ```

2. Or temporarily disable trust requirements:
   ```bash
   export HOMEBREW_NO_REQUIRE_TAP_TRUST=1
   brew install --cask cloudfoundry/tap/uaa-cli
   ```

## Development notes

### Setting up Go

If you don't have a working Go setup

```
brew update
brew install go

echo 'export GOPATH="$HOME/go"' >> ~/.bash_profile
echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.bash_profile
```

### Running the tests

```
cd $GOPATH/src/code.cloudfoundry.org/uaa-cli
ginkgo -r -randomizeAllSpecs -randomizeSuites
```
