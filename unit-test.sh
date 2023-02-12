#!/usr/bin/env bash

for dir in service/*/; do
  pushd "$dir" || exit
    go test -gcflags=-l -v -cover ./...
  popd || exit
done
