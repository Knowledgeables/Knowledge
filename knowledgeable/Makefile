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
	@echo   make swagger

test:
	go test ./... -count=1 -v

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

## This target sets up the pre-commit hooks defined in .pre-commit-config.yaml
## Prek install by default installs pre-commit hooks, but we also want to install commit-msg hooks for validating commit messages
setup-hooks:
	prek install --hook-type pre-commit
	prek install --hook-type commit-msg

## This target is for running the pre-commit hooks manually, e.g. in CI pipelines
hooks-run:
	prek run
	setup-hooks.sh

swagger:
	rm -rf docs
	swag init -g cmd/server/main.go -o docs
