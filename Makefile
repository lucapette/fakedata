SOURCE_FILES?=$$(go list ./pkg/...)
TEST_PATTERN?=.
TEST_OPTIONS?=

unit: ## Run tests
	@go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

integration: build-with-cover ## Run integration tests
	@go test $(TEST_OPTIONS) $$(go list ./integration/...) -timeout=30s
	@go tool covdata percent -i=.coverdata

test: unit integration

bench: ## Run benchmarks
	@go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -bench $(TEST_PATTERN) -timeout=30s

lint: ## Run linters
	@golangci-lint run

build: ## Build a dev version of fakedata
	@go build

build-debug-image:
	@GOOS=linux GOARCH=amd64 go build
	docker build -t fakedata .

build-with-cover: ## Build a cover version of fakedata
	@rm -fr .coverdata
	@mkdir .coverdata
	@go build -cover -o ./fakedata-with-cover

import: ## Import or update data from dariusk/corpora
	@go run cmd/import/main.go

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
