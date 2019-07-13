#!/usr/bin/env bash

pushd $GOPATH/src/pharmer.dev/pre-k/hack/gendocs
go run main.go
popd
