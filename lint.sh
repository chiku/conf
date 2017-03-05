#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

go get -u github.com/alecthomas/gometalinter
gometalinter --install
gometalinter -t --vendor ./... --deadline 1m
