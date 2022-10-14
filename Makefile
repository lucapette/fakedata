SOURCE_FILES?=$$(go list ./...)
TEST_PATTERN?=.
TEST_OPTIONS?=

test: ## Run all the tests
	@go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

lint: ## Run all the linters
	@golangci-lint run

build: ## Build a dev version of fakedata
	@go build

import: ## Import or update data from dariusk/corpora
	@go run cmd/import/main.go

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
