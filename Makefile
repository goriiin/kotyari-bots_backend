defalut: help

help:
	@echo 'usage: make [target]'
	@echo 'targets:'
	@echo 'download_lint - Downloading linter binary'
	@echo 'check_lint - Verify linter version (>= 2)'
	@echo 'verify_lint_config - Verifies linter config'
	@echo 'lint - running linter'


download_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.3.1

check_lint:
	golangci-lint --version

verify_lint_config:
	golangci-lint config verify

lint:
	golangci-lint run
