defalut: help

SERVICES := $(shell find ./api -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

help:
	@echo ''
	@echo 'usage: make [target]'
	@echo ''
	@echo 'targets:'
	@echo '	download-lint - Downloading linter binary'
	@echo '	check-lint - Verify linter version (>= 2)'
	@echo '	verify-lint-config - Verifies linter config'
	@echo '	lint - running linter'
	@echo '	download-gci - Downloading import formatter'
	@echo '	install - Download all dev tools (linter, formatter)'
	@echo '	format - Format go import statements'
	@echo '	format-check - Check go import statements formatting'
	@echo '	check - Run all checks (lint, format-check)'
	@echo "api          - Сгенерировать Go-код из всех openapi.yml файлов."
	@echo "install-ogen - Установить или обновить генератор кода ogen."


api: install-ogen
	@echo "Начинаю генерацию кода для сервисов: $(SERVICES)"
	$(foreach service,$(SERVICES),$(call generate-service,$(service)))
	@echo "Генерация кода успешно завершена."

install-ogen:
	@if ! command -v ogen &> /dev/null; then \
		echo "ogen не найден. Устанавливаю..."; \
		go install github.com/ogen-go/ogen/cmd/ogen@latest; \
	fi

define generate-service
	@echo "--- Генерирую код для сервиса: $(1) ---"
	@# Определяем пути
	$(eval INPUT_FILE := ./docs/$(1)/openapi.yaml)
	$(eval OUTPUT_DIR := ./internal/gen/$(1))

	@# Проверяем наличие исходного файла
	@if [ ! -f "$(INPUT_FILE)" ]; then \
		echo "Ошибка: Файл спецификации $(INPUT_FILE) не найден!"; \
		exit 1; \
	fi

	@# Запускаем ogen
	ogen --target "$(OUTPUT_DIR)" --package "$(1)" -clean "$(INPUT_FILE)"
endef


download-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.3.1

download-gci:
	go install github.com/daixiang0/gci@v0.13.4

install: download-lint download-gci

check-lint:
	golangci-lint --version

verify-lint-config:
	golangci-lint config verify

lint:
	golangci-lint run

format:
	@gci write . --skip-generated --skip-vendor < /dev/null

format-check:
	@gci diff . --skip-generated --skip-vendor < /dev/null

check: lint format-check

example-run:
	@go run cmd/main/main.go
example-run-local:  ## Запустить в local режиме
	@go run cmd/main/main.go --env=local --config="./configs/local-config.yaml"

example-run-prod:  ## Запустить в production режиме
	@go run cmd/main/main.go --env=prod


.PHONY: download-lint download-gci lint format format-check check help api
