# UAA Command Line Interface, PR test

CLI for [UAA](https://github.com/cloudfoundry/uaa) written in golang. This is an alterntive to using uaac which is wirtten in Ruby.  At this time it performs a limited subset of the features provided by the [uaac](https://github.com/cloudfoundry/cf-uaac) gem.  The team plans to continue development on the golang CLI going forward, and once it's considered fully GA, intends to place it alongside uaac with a long-term intention of one day deprecating uaac.

### Goals

- To provide a CLI which can be easily installed in environments without a functioning Ruby setup
- To more closely conform to the style of other widely used CLIs in the CF ecosystem, e.g. the cf CLI. Commands should be of the form VERB-NOUN, similar to `cf delete-app`.
- To provide outputs that are machine-parseable whenever possible.
- To improve the quality of help strings and error messages so that users can self-diagnose problems and unblock themselves.
- To provide only the essential, highly used and/or required command options.

### Trying out the latest code

```
go get code.cloudfoundry.org/uaa-cli
cd $GOPATH/src/code.cloudfoundry.org/uaa-cli
make && make install
uaa -h
```
Or, install it using brew.  It's been made available as part of the [cloudfoundry tap](https://github.com/cloudfoundry/homebrew-tap)

```
brew install cloudfoundry/tap/uaa-cli
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
