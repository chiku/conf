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

MKDIR = mkdir -p
RM = rm -rvf
GO = go
GLIDE = glide

sources := $(wildcard *.go)

.PHONY: all
all: prereq fmt vet test

.PHONY: prereq
prereqs:
	${GLIDE} install

.PHONY: fmt
fmt:
	${GO} fmt

.PHONY: vet
vet:
	${GO} vet

.PHONY: test
test: coverage/coverage.html

coverage:
	${MKDIR} coverage

.PHONY: clean
clean:
	${RM} coverage

coverage/coverage.out: $(sources) coverage
	${GO} test -coverprofile=coverage/coverage.out

coverage/coverage.html: coverage/coverage.out
	${GO} tool cover -func=coverage/coverage.out
	${GO} tool cover -html=coverage/coverage.out -o coverage/coverage.html
