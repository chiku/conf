# Makefile
#
# Author::    Chirantan Mitra
# Copyright:: Copyright (c) 2015-2016. All rights reserved
# License::   MIT

MAKEFLAGS += --warn-undefined-variables
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: all prereq fmt vet test clean

sources := $(wildcard *.go)

all: prereq fmt vet test

prereqs:
	glide install

fmt:
	go fmt

vet:
	go vet

test: out/coverage.html

out:
	mkdir -p out

clean:
	rm -rfv out

out/coverage.out: $(sources) out
	go test -coverprofile=out/coverage.out

out/coverage.html: $(sources) out/coverage.out
	go tool cover -func=out/coverage.out
	go tool cover -html=out/coverage.out -o out/coverage.html
