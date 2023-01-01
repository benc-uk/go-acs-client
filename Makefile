.EXPORT_ALL_VARIABLES:
.PHONY: help lint lint-fix test
.DEFAULT_GOAL := help

help:  ## 💬 This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## 📝 Lint & format, check to be run in CI, sets exit code on error
	golangci-lint run --modules-download-mode=mod --timeout=4m ./...

lint-fix: ## 🧙 Lint & format, fixes errors and modifies code
	golangci-lint run --modules-download-mode=mod --timeout=4m --fix ./...

test:  ## 🎯 Run integration tests
	@echo -e "WARNING: This will run integration tests\nThis will send several real emails!"
	go test -v ./...
