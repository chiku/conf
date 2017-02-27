#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

go get -v github.com/mattn/goveralls
goveralls -package github.com/chiku/conf -repotoken $COVERALLS_TOKEN -service=circle-ci
