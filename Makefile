defalut: help

help:
	@echo ''
	@echo 'usage: make [target]'
	@echo ''
	@echo 'targets:'
	@echo '	download_lint - Downloading linter binary'
	@echo '	check_lint - Verify linter version (>= 2)'
	@echo '	verify_lint_config - Verifies linter config'
	@echo '	lint - running linter'
	@echo '	download_gci - Downloading import formatter'
	@echo '	install - Download all dev tools (linter, formatter)'
	@echo '	format - Format go import statements'
	@echo '	format_check - Check go import statements formatting'
	@echo '	check - Run all checks (lint, format_check)'

download_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.3.1

download_gci:
	go install github.com/daixiang0/gci@v0.13.4

install: download_lint download_gci

check_lint:
	golangci-lint --version

verify_lint_config:
	golangci-lint config verify

lint:
	golangci-lint run

format:
	gci write .

format_check:
	gci diff .

check: lint format_check

example-run:
	@go run cmd/main/main.go
example-run-local:  ## Запустить в local режиме
	@go run cmd/main/main.go --env=local --config="./configs/config-local.yaml"

example-run-prod:  ## Запустить в production режиме
	@go run cmd/main/main.go --env=prod
