#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

mkdir -p out/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./out -d v2.11.4
./out/golangci-lint run ./... --timeout 1m
