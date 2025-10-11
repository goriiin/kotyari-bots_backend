DOCKER_NETWORK := public-gateway-network

defalut: help

SERVICES := $(shell find ./docs -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: help up down reboot test

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

profiles-reboot: profiles-down profiles-up


# --- Вспомогательные и внутренние команды ---

.PHONY: setup-network teardown-network copy-env


setup-network:
	@docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || \
		(echo "Создаю общую Docker-сеть: $(DOCKER_NETWORK)..." && docker network create $(DOCKER_NETWORK))

# Удаляет общую сеть
teardown-network:
	@docker network rm $(DOCKER_NETWORK) >/dev/null 2>&1 || true

copy-env:
	@if [ ! -f .env ]; then \
		echo "Создаю .env файл из .env.example..."; \
		cp .env.example .env; \
	fi

# --- Кодогенерация и статический анализ ---

PROTO_DIR := ./api/protos
GEN_DIR := gen
PROTOC := protoc
ENTITIES := $(shell find $(PROTO_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

proto-build: $(ENTITIES)

$(ENTITIES):
	@echo "Генерация кода для сущности $@..."
	@mkdir -p $(PROTO_DIR)/$@/$(GEN_DIR)
	@$(PROTOC) \
		--proto_path=$(PROTO_DIR)/$@/proto \
		--go_out=$(PROTO_DIR)/$@/$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR)/$@/$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/$@/proto/*.proto
	@echo "Генерация для $@ завершена."

api: install-ogen
	@echo "Начинаю генерацию кода для сервисов: $(SERVICES)"
	$(foreach service,$(SERVICES),$(call generate-service,$(service)))
	@echo "Генерация кода успешно завершена."

install-ogen:
	go install github.com/ogen-go/ogen/cmd/ogen@latest

define generate-service
	@echo "--- Генерирую код для сервиса: $(1) ---"
	$(eval INPUT_FILE := ./docs/$(1)/openapi.yaml)
	$(eval OUTPUT_DIR := ./internal/gen/$(1))
	@if [ ! -f "$(INPUT_FILE)" ]; then \
		echo "Ошибка: Файл спецификации $(INPUT_FILE) не найден!"; \
		exit 1; \
	fi
	ogen --target "$(OUTPUT_DIR)" --package "$(1)" -clean "$(INPUT_FILE)"
endef

download-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.3.1

download-gci:
	go install github.com/daixiang0/gci@v0.13.4

install: download-lint download-gci

lint:
	golangci-lint run

format:
	@gci write . --skip-generated --skip-vendor < /dev/null

format-check:
	@gci diff . --skip-generated --skip-vendor < /dev/null

check: lint format-check

# параллельно
up: copy-env setup-network
	@echo "Starting services in parallel..."
	@$(MAKE) bots-up & \
	$(MAKE) profiles-up & \
	wait
	@echo "All services are up and running."

bots-up: setup-network
	@echo "Starting bots service and dependencies..."
	@docker compose -f docker-compose.bots.yml up -d --build

bots-down:
	@echo "Stopping bots service and dependencies..."
	@docker compose -f docker-compose.bots.yml down

bots-reboot:
	@echo "Rebooting bots service and dependencies..."
	$(MAKE) bots-down
	$(MAKE) bots-up

profiles-up: setup-network
	@echo "Starting profiles service and dependencies..."
	docker-compose -f docker-compose.profiles.yml up -d --build

profiles-down:
	@echo "Stopping profiles service and dependencies..."
	@docker-compose -f docker-compose.profiles.yml down

profiles-reboot:
	@echo "Rebooting profiles service and dependencies..."
	$(MAKE) profiles-down
	$(MAKE) profiles-up

example-run:
	@go run cmd/example/main.go
example-run-local:  ## Запустить в local режиме
	@go run cmd/example/main.go --env=local --config="./configs/local-config.yaml"

example-run-prod:
	@go run cmd/example/main.go --env=prod

install-migrate:
	@if ! command -v migrate &> /dev/null; then \
		echo "migrate CLI not found. Installing..."; \
		go install -tags 'pgx5' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi

.PHONY: download-lint download-gci lint format format-check check help api

INTRANET_DIR := ./intranet

intranet-up:
	$(MAKE) -C $(INTRANET_DIR) up

intranet-down:
	$(MAKE) -C $(INTRANET_DIR) down

intranet-parser-logs:
	$(MAKE) -C $(INTRANET_DIR) parser-logs