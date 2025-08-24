defalut: help

SERVICES := $(shell find ./api -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

.PHONY: help api

help:
	@echo 'usage: make [target]'
	@echo 'targets:'
	@echo 'download_lint - Downloading linter binary'
	@echo 'check_lint - Verify linter version (>= 2)'
	@echo 'verify_lint_config - Verifies linter config'
	@echo 'lint - running linter'
	@echo "api          - Сгенерировать Go-код из всех openapi.yml файлов."
	@echo "install-ogen - Установить или обновить генератор кода ogen."


api: install_ogen
	@echo "Начинаю генерацию кода для сервисов: $(SERVICES)"
	$(foreach service,$(SERVICES),$(call generate_service,$(service)))
	@echo "Генерация кода успешно завершена."

install_ogen:
	@if ! command -v ogen &> /dev/null; then \
		echo "ogen не найден. Устанавливаю..."; \
		go install github.com/ogen-go/ogen/cmd/ogen@latest; \
	fi

define generate_service
	@echo "--- Генерирую код для сервиса: $(1) ---"
	@# Определяем пути
	$(eval INPUT_FILE := ./api/$(1)/openapi.yaml)
	$(eval OUTPUT_DIR := ./internal/gen/$(1))

	@# Проверяем наличие исходного файла
	@if [ ! -f "$(INPUT_FILE)" ]; then \
		echo "Ошибка: Файл спецификации $(INPUT_FILE) не найден!"; \
		exit 1; \
	fi

	@# Запускаем ogen
	ogen --target "$(OUTPUT_DIR)" --package "$(1)" -clean "$(INPUT_FILE)"
endef

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
	@gci write . --skip-generated --skip-vendor < /dev/null

format_check:
	@gci diff . --skip-generated --skip-vendor < /dev/null

check: lint format_check

example-run:
	@go run cmd/main/main.go
example-run-local:  ## Запустить в local режиме
	@go run cmd/main/main.go --env=local --config="./configs/config-local.yaml"

example-run-prod:  ## Запустить в production режиме
	@go run cmd/main/main.go --env=prod