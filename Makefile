

run:
	@go run cmd/main/main.go
run-local:  ## Запустить в local режиме
	@go run cmd/main/main.go --env=local --config="./configs/config-local.yaml"

run-prod:  ## Запустить в production режиме
	@go run cmd/main/main.go --env=prod
