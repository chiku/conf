#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

mkdir -p out/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./out v1.40.1
./out/golangci-lint run ./... --deadline 1m
