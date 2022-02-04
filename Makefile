# List of effective go files
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" | egrep -v "^\./\.go" | grep -v _test.go)

# List of packages except testsutils
PACKAGES ?= $(shell go list ./... | egrep -v "testutils" )

COVERAGE_BUILD_FOLDER = build/coverage
COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/cov.out
COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/index.html

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

test:
	@mkdir -p build/coverage
	@go test -covermode=count -coverprofile $(COVERAGE_OUT) $(PACKAGES)

test-and-generate-coverage-html: test
	@go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	
