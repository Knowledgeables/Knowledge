.PHONY: help test lint dev-up dev-down dev-build dev-logs setup-hooks

help:
	@echo Available targets:
	@echo   make test
	@echo   make lint
	@echo   make dev-up
	@echo   make dev-down
	@echo   make dev-build
	@echo   make dev-logs
	@echo   make setup-hooks

test:
	go test ./...

lint:
	golangci-lint run

dev-up:
	docker compose -f docker-compose-dev.yml up -d

dev-down:
	docker compose -f docker-compose-dev.yml down

dev-build:
	docker compose -f docker-compose-dev.yml build

dev-logs:
	docker compose -f docker-compose-dev.yml logs -f

setup-hooks:
	setup-hooks.sh

swagger:
	rm -rf docs
	swag init -g cmd/server/main.go -o docs