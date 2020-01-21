#! /usr/bin/make
#
# Makefile for goa
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
#
# Meta targets:
# - "all" is the default target, it runs all the targets in the order above.
#

all: depend lint goagen

.PHONY: depend
depend:
	go mod download
	go get github.com/onsi/ginkgo/ginkgo

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites -race
	go test -v github.com/shogo82148/goa-v1/_integration_tests

.PHONY: goagen
goagen:
	@cd goagen && \
	go install
