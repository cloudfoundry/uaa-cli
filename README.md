# UAA Command Line Interface

[![Travis CI](https://travis-ci.org/cloudfoundry-incubator/uaa-cli.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/uaa-cli)

Experimental CLI for [UAA](https://github.com/cloudfoundry/uaa) written in golang. At this time it performs a limited subset of the features provided by the [uaac](https://github.com/cloudfoundry/cf-uaac) gem.

### Goals

- To provide a CLI which can be easily installed in environments without a functioning Ruby setup
- To more closely conform to the style of other widely used CLIs in the CF ecosystem, e.g. the cf CLI. Commands should be of the form VERB-NOUN, similar to `cf delete-app`.
- To provide outputs that are machine-parseable whenever possible.
- To improve the quality of help strings and error messages so that users can self-diagnose problems and unblock themselves.

### Roadmap

The immediate goal is to reach feature-parity with the uaac. Right now 
tasks are being tracked in this [trello board](https://trello.com/b/Hw4Pz0Jd/uaa-cli).

### Trying out the latest code

```
go get code.cloudfoundry.org/uaa-cli
cd $GOPATH/src/code.cloudfoundry.org/uaa-cli
make && make install
uaa -h
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
