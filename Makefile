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
DIRS=$(shell go list -f {{.Dir}} ./...)

.PHONY: goagen

all: depend lint cyclo goagen test

docs:
	@go get -v github.com/spf13/hugo
	@git clone https://github.com/shogo82148/goa-v1.design
	@rm -rf goa.design/content/reference goa.design/public
	@mdc --exclude goa.design github.com/shogo82148/goa-v1 goa.design/content/reference
	@cd goa.design && hugo
	@rm -rf public
	@mv goa.design/public public
	@rm -rf goa.design

.PHONY: depend
depend:
	go mod download
	GO111MODULE=off go get github.com/onsi/ginkgo/ginkgo

lint:
	@for d in $(DIRS) ; do \
		if [ "`goimports -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
		fi \
	done
	@if [ "`golint ./... | grep -vf .golint_exclude | tee /dev/stderr`" ]; then \
		echo "^ - Lint errors!" && echo && exit 1; \
	fi

cyclo:
	@if [ "`gocyclo -over 20 . | grep -v _integration_tests | grep -v _test.go | tee /dev/stderr`" ]; then \
		echo "^ - Cyclomatic complexity exceeds 20, refactor the code!" && echo && exit 1; \
	fi

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites -race
	go test ./_integration_tests

goagen:
	@cd goagen && \
	go install
