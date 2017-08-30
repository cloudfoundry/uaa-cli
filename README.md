# uaa-cli

Experimental UAA cli. Not really ready for public consumption.

## Setting up Go

If you don't have a working Go setup

```
brew update
brew install go
brew install dep

echo 'export GOPATH="$HOME/go"' >> ~/.bash_profile
echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.bash_profile
```

## Trying out the latest code

```
go get code.cloudfoundry.org/uaa-cli
cd $GOPATH/src/code.cloudfoundry.org/uaa-cli
dep ensure
go build && ./uaa-cli
```

## Running the tests

```
cd $GOPATH/src/code.cloudfoundry.org/uaa-cli
ginkgo -r
```
