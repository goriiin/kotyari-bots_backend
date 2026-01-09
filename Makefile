DOCKER_NETWORK := public-gateway-network
DOMAIN := writehub.space
EMAIL := admin@writehub.space
NGINX_COMPOSE := docker-compose.nginx.yml
FRONTEND_DIR := ../kotyari-bots_frontend

# Автоматический поиск сервисов для Ogen
SERVICES := $(shell find ./docs -mindepth 2 -maxdepth 3 -type f -name 'openapi.yaml' -print \
                              | sed -e 's|^./docs/||' -e 's|/openapi.yaml$$||' | sort -u)

export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: help up down bots-up profiles-up posts-up gateway-up ssl-install frontend-build

default: help

help:
	@echo ''
	@echo 'usage: make [target]'
	@echo ''
	@echo 'MAIN TARGETS:'
	@echo '  up              - Поднять весь бэкенд (БД, Kafka, Go-сервисы) без Nginx'
	@echo '  down            - Остановить весь бэкенд'
	@echo '  gateway-up      - Поднять Nginx (Gateway) + Certbot'
	@echo '  ssl-install     - Полная настройка HTTPS (с генерацией сертификатов)'
	@echo '  frontend-build  - Собрать статику Nuxt и исправить права доступа'
	@echo ''
	@echo 'DEV TOOLS:'
	@echo '  lint            - Запустить линтер'
	@echo '  format          - Отформатировать импорты'
	@echo '  api             - Сгенерировать Go-код (Ogen) из OpenAPI'
	@echo '  proto-build     - Сгенерировать gRPC код из .proto'

# --- СЕТЬ И ОКРУЖЕНИЕ ---

setup-network:
	@docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || \
		(echo "Создаю общую Docker-сеть: $(DOCKER_NETWORK)..." && docker network create $(DOCKER_NETWORK))

copy-env:
	@if [ ! -f .env ]; then \
		echo "Создаю .env файл из .env.example..."; \
		cp .env.example .env; \
	fi

# --- БЭКЕНД (Docker Compose) ---

# Параллельный запуск основных сервисов
up: copy-env setup-network
	@echo "Starting backend services..."
	@$(MAKE) bots-up & \
	 $(MAKE) profiles-up & \
	 $(MAKE) posts-up & \
	wait
	@echo "Backend services are up."

down:
	@echo "Stopping backend services..."
	@$(MAKE) bots-down & \
	 $(MAKE) profiles-down & \
	 $(MAKE) posts-down & \
	wait
	@echo "Backend services stopped."

bots-up: setup-network
	docker compose -f docker-compose.bots.yml up -d --build

bots-down:
	docker compose -f docker-compose.bots.yml down

profiles-up: setup-network
	docker compose -f docker-compose.profiles.yml up -d --build

profiles-down:
	docker compose -f docker-compose.profiles.yml down

posts-up: setup-network
	docker compose -f docker-compose.posts.yml up -d --build

posts-down:
	docker compose -f docker-compose.posts.yml down

# --- FRONTEND ---

frontend-build:
	@echo "Building Frontend Static Site..."
	cd $(FRONTEND_DIR) && npm run generate
	@echo "Fixing permissions for Nginx..."
	chmod -R 755 $(FRONTEND_DIR)/.output/public
	@echo "Frontend built successfully."

# --- GATEWAY & SSL (NGINX) ---

gateway-up: setup-network
	docker compose -f $(NGINX_COMPOSE) up -d

gateway-down:
	docker compose -f $(NGINX_COMPOSE) down

gateway-restart:
	docker compose -f $(NGINX_COMPOSE) restart gateway

gateway-logs:
	docker compose -f $(NGINX_COMPOSE) logs -f

ssl-install:
	@if [ ! -f nginx.conf.http ] || [ ! -f nginx.conf.https ]; then \
		echo "Ошибка: Файлы nginx.conf.http и nginx.conf.https должны существовать."; \
		exit 1; \
	fi
	@echo ">>> [1/4] Применяем HTTP конфигурацию (для валидации)..."
	cp nginx.conf.http nginx.conf
	$(MAKE) gateway-up
	@echo ">>> Ожидание запуска Nginx..."
	@sleep 5
	@echo ">>> [2/4] Генерация сертификатов через Let's Encrypt..."
	docker compose -f $(NGINX_COMPOSE) run --rm --entrypoint certbot certbot certonly --webroot --webroot-path /var/www/certbot \
		-d $(DOMAIN) -d www.$(DOMAIN) \
		--email $(EMAIL) \
		--agree-tos --no-eff-email --force-renewal
	@echo ">>> [3/4] Применяем HTTPS конфигурацию (боевую)..."
	cp nginx.conf.https nginx.conf
	@echo ">>> [4/4] Перезагрузка Nginx..."
	docker compose -f $(NGINX_COMPOSE) exec gateway nginx -s reload
	@echo ">>> Готово. HTTPS настроен."

cert-renew:
	docker compose -f $(NGINX_COMPOSE) run --rm --entrypoint certbot certbot renew
	docker compose -f $(NGINX_COMPOSE) exec gateway nginx -s reload

# --- CODE GEN & LINTING ---

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

install-ogen:
	go install github.com/ogen-go/ogen/cmd/ogen@v1.16.0

api: install-ogen
	@echo "Начинаю генерацию кода для сервисов: $(SERVICES)"
	$(foreach service,$(SERVICES),$(call generate-service,$(service)))
	@echo "Генерация кода успешно завершена."

define generate-service
	@echo "--- Генерирую код для сервиса: $(1) ---"
	$(eval INPUT_FILE := ./docs/$(1)/openapi.yaml)
	$(eval OUTPUT_DIR := ./internal/gen/$(1))
	$(eval PKG := $(notdir $(1)))
	$(eval OGEN_CFG  := ./docs/ogen-config.yaml)
	@mkdir -p "$(OUTPUT_DIR)"
	ogen --config "$(OGEN_CFG)" --target "$(OUTPUT_DIR)" --package "$(PKG)" -clean "$(INPUT_FILE)"
endef

download-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0

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

# --- INTRANET (Parsers) ---

INTRANET_DIR := ./intranet

intranet-up-dev:
	$(MAKE) -C $(INTRANET_DIR) up-dev

intranet-down-dev:
	$(MAKE) -C $(INTRANET_DIR) down-dev

intranet-up-prod:
	$(MAKE) -C $(INTRANET_DIR) up-prod

intranet-down-prod:
	$(MAKE) -C $(INTRANET_DIR) down-prod