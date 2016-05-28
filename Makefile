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

ifndef GOPATH
$(error GOPATH not set)
endif

MKDIR = mkdir -p
RM = rm -rvf
GO = go
GLIDE = $(GOPATH)/bin/glide

sources := $(wildcard *.go)
coverage = coverage
coverage_out = $(coverage)/coverage.out
coverage_html = $(coverage)/coverage.html

.PHONY: all
all: fmt vet test

.PHONY: fmt
fmt:
	${GO} fmt

.PHONY: vet
vet:
	${GO} vet

.PHONY: test
test: $(coverage_html)

.PHONY: fuzz
fuzz:
	${GO} test -v ./fuzz

$(coverage):
	${MKDIR} $(coverage)

.PHONY: clean
clean:
	${RM} $(coverage)

$(coverage_out): $(sources) $(coverage)
	${GO} test -coverprofile=$(coverage_out)

$(coverage_html): $(coverage_out)
	${GO} tool cover -func=$(coverage_out)
	${GO} tool cover -html=$(coverage_out) -o $(coverage_html)
