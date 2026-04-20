.PHONY: help go-test smoke-test lint dev-up dev-down dev-build dev-logs setup-hooks dev-migrate hooks-run install-gosec swagger

# =========================
# Project Commands (make help)
# =========================
help:
	@echo "Available targets:"
	@echo ""
	@echo "Setup:"
	@echo "  make setup-hooks"
	@echo "  make install-gosec"
	@echo ""
	@echo "Development:"
	@echo "  make dev-up"
	@echo "  make dev-down"
	@echo "  make dev-build"
	@echo "  make dev-logs"
	@echo "  make dev-migrate"
	@echo ""
	@echo "Testing:"
	@echo "  make go-test"
	@echo "  make smoke-test"
	@echo "  make full-e2e"
	@echo ""
	@echo "Quality:"
	@echo "  make lint"
	@echo "  make hooks-run"
	@echo ""
	@echo "Other:"
	@echo "  make swagger"

# =========================
# Setup (run once)
# =========================
setup-hooks:
	prek install --hook-type pre-commit
	prek install --hook-type commit-msg

install-gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# =========================
# Development Environment
# =========================
dev-up:
	docker compose -f docker-compose-dev.yml up -d
	docker compose -f docker-compose-dev.yml exec -T db sh -c "until pg_isready -U postgres; do sleep 2; done"
	$(MAKE) dev-migrate

SEED ?= dev_seed.js

dev-migrate:
	docker compose -f docker-compose-dev.yml run --rm migrations sh -c "npx knex migrate:latest && \
npx knex seed:run --specific=$(SEED)"

dev-down:
	docker compose -f docker-compose-dev.yml down

dev-build:
	docker compose -f docker-compose-dev.yml build

dev-logs:
	docker compose -f docker-compose-dev.yml logs -f

# =========================
# Testing
# =========================
go-test:
	go test ./... -count=1 -v

smoke-test: ## Run Playwright smoke tests (requires running app)
	npx playwright test --project=smoke --reporter=line

full-e2e: ## Run the full version of e2e playwright tests that runs on staging environment
	$(MAKE) dev-down
	$(MAKE) dev-up SEED=staging_seed.js
	npx playwright test --project=full --reporter=line
	$(MAKE) dev-down
	$(MAKE) dev-up

# =========================
# Code Quality
# =========================
lint:
	golangci-lint run

hooks-run:
	prek run
	setup-hooks.sh

# =========================
# Tooling / Docs
# =========================
swagger:
	rm -rf docs
	swag init -g cmd/server/main.go -o docs