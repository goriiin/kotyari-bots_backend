ifneq (,$(wildcard ./.env))
    include .env
    export
endif

up:
	docker-compose -f docker-compose.yml up -d --build

down:
	docker-compose -f docker-compose.yml down

parser-logs:
	docker-compose -f docker-compose.yml logs -f parser_app

help:
	@echo 'Intranet Service Local Commands (run from intranet/ directory):'
	@echo '  make up           - Запустить сервисы'
	@echo '  make down         - Остановить сервисы'
	@echo '  make parser-logs         - Показать логи'
	@echo ''

.PHONY: up down logs cli redis-cli publish help