# See https://tech.davis-hansson.com/p/make/
SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory
BIN := .tmp/bin
export PATH := $(abspath $(BIN)):$(PATH)
export GOBIN := $(abspath $(BIN))

GO_MOD_GOTOOLCHAIN := go1.24.0
GOLANGCI_LINT_VERSION := v1.60.3
# https://github.com/golangci/golangci-lint/issues/4837
GOLANGCI_LINT_GOTOOLCHAIN := $(GO_MOD_GOTOOLCHAIN)
# If any pins to specific dependency versions are needed, add them here
GO_GET_PKGS :=

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: all
all: ## Build, test, and lint (default)
	$(MAKE) test

.PHONY: clean
clean: ## Delete intermediate build artifacts
	@# -X only removes untracked files, -d recurses into directories, -f actually removes files/dirs
	git clean -Xdf

.PHONY: test
test: build ## Run unit tests
	go test -vet=off -race -cover ./...

.PHONY: build
build: ## Build all packages
	go build ./...

.PHONY: install
install: ## Install all binaries
	go install ./...

.PHONY: lint
lint: $(BIN)/golangci-lint ## Lint
	go vet ./...
	GOTOOLCHAIN=$(GOLANGCI_LINT_GOTOOLCHAIN) golangci-lint run --modules-download-mode=readonly --timeout=3m0s
	go run ./internal/cmd/quality-checker

.PHONY: lintfix
lintfix: $(BIN)/golangci-lint ## Automatically fix some lint errors
	GOTOOLCHAIN=$(GOLANGCI_LINT_GOTOOLCHAIN) golangci-lint run --fix --modules-download-mode=readonly --timeout=3m0s

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	go mod edit -toolchain=$(GO_MOD_GOTOOLCHAIN)
	go get -u -t ./... $(GO_GET_PKGS)
	go mod tidy -v

$(BIN)/golangci-lint: Makefile
	@mkdir -p $(@D)
	GOTOOLCHAIN=$(GOLANGCI_LINT_GOTOOLCHAIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
