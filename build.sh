#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

rm -rfv ./out
mkdir -pv ./out/coverage

go fmt . ./examples
go vet . ./examples

go test -coverprofile=./out/coverage/coverage.out
go tool cover -func=./out/coverage/coverage.out
go tool cover -html=./out/coverage/coverage.out -o ./out/coverage/coverage.html

go build -o out/example examples/example.go
