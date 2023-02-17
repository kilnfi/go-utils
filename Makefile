# List of effective go files
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" | egrep -v "^\./\.go" | grep -v _test.go)

# List of packages except testsutils
PACKAGES ?= $(shell go list ./... | grep -v "mock" )

# Build folder
BUILD_FOLDER = build

# Test coverage variables
COVERAGE_BUILD_FOLDER = $(BUILD_FOLDER)/coverage

UNIT_COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/ut_cov.out
UNIT_COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/ut_index.html

INTEGRATION_COVERAGE_OUT = $(COVERAGE_BUILD_FOLDER)/it_cov.out
INTEGRATION_COVERAGE_HTML =$(COVERAGE_BUILD_FOLDER)/it_index.html

# Test lint variables
GOLANGCI_VERSION = v1.48.0

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

build/coverage:
	@mkdir -p build/coverage

unit-test: build/coverage
	@go test -covermode=count -coverprofile $(UNIT_COVERAGE_OUT) -v $(PACKAGES)

unit-test-cov: unit-test
	@go tool cover -html=$(UNIT_COVERAGE_OUT) -o $(UNIT_COVERAGE_HTML)

fix-lint: ## Run linter to fix issues
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run --fix

# @misspell -error $(GOFILES)
test-lint: ## Check linting
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run -v

integration-test: build/coverage
	@go test -covermode=count -coverprofile $(INTEGRATION_COVERAGE_OUT) -v --tags integration ${PACKAGES}

integration-test-cov: integration-test
	@go tool cover -html=$(INTEGRATION_COVERAGE_OUT) -o $(INTEGRATION_COVERAGE_HTML)
