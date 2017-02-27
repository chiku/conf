#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

go get -v github.com/google/gofuzz
go test -v . -timeout 1h -tags fuzz
