# guac

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
go get github.com/jhamon/guac
cd $GOPATH/src/github.com/jhamon/guac
deps ensure
go build && ./guac
```

## Running the tests

```
cd $GOPATH/src/github.com/jhamon/guac
ginkgo -r
```