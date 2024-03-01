# Overview

# Steps for building locally

```shell
go build
```

## Run unit tests

To run all unit tests from command line

```shell
go test ./...
```
# Install Go

## Download

Get latest go from:
https://go.dev/dl/

Installation instructions:
https://go.dev/doc/install

## Install on Linux
(based on https://go.dev/doc/install)

If previous version of go exists, remove it using `rm -rf /usr/local/go`.

```shell
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
```

Add `/usr/local/go/bin` to the PATH environment variable in your `$HOME/.bashrc` file.

```shell
export PATH=$PATH:/usr/local/go/bin
```

Source the new setting:

```shell
source $HOME/.bashrc
```

Verify

```shell
go version
```
