# List of effective go files
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" | egrep -v "^\./\.go" | grep -v _test.go)

# List of packages except testsutils
PACKAGES ?= $(shell go list ./... | egrep -v "testutils" )

# Test coverage variables
COVERAGE_BUILD_FOLDER = build/coverage
COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/cov.out
COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/index.html

# Test lint variables
GOLANGCI_VERSION = v1.44.0

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


fix-lint: ## Run linter to fix issues
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run --fix

# @misspell -error $(GOFILES)
test-lint: ## Check linting
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run -v

