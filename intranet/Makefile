ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DC ?= docker compose
DEV_FILE ?= docker-compose.dev.yml
PROD_FILE ?= docker-compose.prod.yml

# DEV (локальная разработка, macOS/ARM с Postgres)
up-dev:
	$(DC) -f $(DEV_FILE) up -d --build

down-dev:
	$(DC) -f $(DEV_FILE) down

ps-dev:
	$(DC) -f $(DEV_FILE) ps

SERVICE ?= parser_app
logs-dev:
	$(DC) -f $(DEV_FILE) logs -f $(SERVICE)

restart-dev:
	$(DC) -f $(DEV_FILE) restart $(SERVICE)

# PROD (для использования на сервере)
up-prod:
	$(DC) -f $(PROD_FILE) up -d --build

down-prod:
	$(DC) -f $(PROD_FILE) down

ps-prod:
	$(DC) -f $(PROD_FILE) ps

logs-prod:
	$(DC) -f $(PROD_FILE) logs -f $(SERVICE)

restart-prod:
	$(DC) -f $(PROD_FILE) restart $(SERVICE)

CMD ?= /bin/sh
exec-dev:
	$(DC) -f $(DEV_FILE) exec -it $(SERVICE) $(CMD)

exec-prod:
	$(DC) -f $(PROD_FILE) exec -it $(SERVICE) $(CMD)

help:
	@echo 'Targets:'
	@echo '  make up-dev|down-dev    - Up/Down dev stack'
	@echo '  make logs-dev SERVICE=X - Follow logs in dev (default SERVICE=parser_app)'
	@echo '  make up-prod|down-prod  - Up/Down prod stack'
	@echo '  make logs-prod SERVICE=X- Follow logs in prod'
	@echo '  make ps-dev|ps-prod     - Show containers'
	@echo '  make restart-*(SERVICE) - Restart service'
	@echo '  make exec-*(SERVICE,CMD)- Exec into container'

.PHONY: up-dev down-dev ps-dev logs-dev restart-dev up-prod down-prod ps-prod logs-prod restart-prod exec-dev exec-prod help


# Test anti-detection in Docker
.PHONY: test-detection
test-detection:
	@echo "Building test container..."
	docker build -f intranet/Dockerfile.test -t intranet-test .
	@echo "Running anti-detection tests..."
	docker run --rm \
		-v $(PWD)/Webshare-10-proxies.txt:/app/Webshare-10-proxies.txt:ro \
		-v $(PWD)/intranet:/app/intranet:ro \
		-e PYTHONUNBUFFERED=1 \
		intranet-test

# Run test with docker-compose
.PHONY: test-detection-compose
test-detection-compose:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Clean test containers
.PHONY: test-clean
test-clean:
	docker-compose -f docker-compose.test.yml down -v
	docker rmi intranet-test 2>/dev/null || true