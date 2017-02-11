# Makefile
#
# Author::    Chirantan Mitra
# Copyright:: Copyright (c) 2016-2017. All rights reserved
# License::   MIT

MAKEFLAGS += --warn-undefined-variables
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

ifndef GOPATH
$(error GOPATH not set)
endif

MKDIR = mkdir -p
RM = rm -rvf
GO = go

sources := $(wildcard *.go)
gofuzz = github.com/google/gofuzz
coverage = coverage
coverage_out = $(coverage)/coverage.out
coverage_html = $(coverage)/coverage.html

all: fmt vet test out/example
.PHONY: all

fmt:
	${GO} fmt
.PHONY: fmt

vet:
	${GO} vet
.PHONY: vet

test: $(coverage_html)
.PHONY: test

$(GOPATH)/src/$(gofuzz):
	go get $(gofuzz)

fuzz: $(GOPATH)/src/$(gofuzz)
	${GO} test -v . -timeout 1h -tags fuzz
.PHONY: fuzz

out/example: $(sources) examples/example.go
	${GO} build -o out/example ./examples

$(coverage):
	${MKDIR} $(coverage)

clean:
	${RM} $(coverage)
.PHONY: clean

$(coverage_out): $(sources) $(coverage)
	${GO} test -coverprofile=$(coverage_out)

$(coverage_html): $(coverage_out)
	${GO} tool cover -func=$(coverage_out)
	${GO} tool cover -html=$(coverage_out) -o $(coverage_html)
