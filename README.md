# Knowledgeable

A knowledge management web application — rewritten from a legacy Flash/Python 2 monolith to a modern Go service as part of a school project exploring DevOps principles end-to-end.

[![Go CI](https://github.com/Knowledgeables/Knowledge/actions/workflows/ci.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/ci.yml)
[![Security Scan](https://github.com/Knowledgeables/Knowledge/actions/workflows/security_scan.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/security_scan.yml)
[![E2E Tests](https://github.com/Knowledgeables/Knowledge/actions/workflows/playwright_ci.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/playwright_ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/Knowledgeables/Knowledge)](https://github.com/Knowledgeables/Knowledge/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](https://go.dev/)

---

## Overview

Knowledgeable is a searchable knowledge-base platform supporting multiple languages (English & Danish). Users can register, log in, and search pages. The rewrite addresses every major flaw found in the legacy system — SQL injection, MD5 password hashing, hardcoded secrets, and zero test coverage — while introducing a full DevOps lifecycle: automated CI, container builds, security scanning, semantic versioning, and zero-downtime production deploys.

## Features

- User registration and authentication with `bcrypt` password hashing
- Full-text page search with multi-language support (EN / DA)
- Swagger/OpenAPI documentation (`/swagger/`)
- Health check endpoint
- Prometheus metrics endpoint (`/metrics`)
- Docker-first development and production environments
- Automated semantic versioning driven by Conventional Commits

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Database | PostgreSQL 15 |
| Frontend | HTML templates + Tailwind CSS 3.4.19 |
| E2E Testing | Playwright (TypeScript) |
| Containerisation | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Reverse Proxy | Nginx |
| Security scanning | Gosec + GitHub SARIF |
| API docs | Swagger |

## CI/CD Pipeline

Every change goes through the following pipeline before reaching production:

```
Pull Request → main
    ├── Go CI          (lint → build → DB setup → migrations → go test)
    └── Playwright     (smoke tests)

Push to main
    ├── Auto Release   (conventional commit → semver tag)
    ├── Delivery       (build & push Docker images to GHCR)
    │       └── Deployment      (SSH → pull image → redeploy production)
    └── Build Migration Image   (build & push Dockerfile.migrations to GHCR)

Staging
    ├── Deploy Staging  (SSH deploy to staging server)
    └── Playwright      (full Playwright suite on staging)

Daily (00:00 UTC)
    └── Security Scan  (Gosec → SARIF → GitHub Security tab)
```

| Workflow file | Trigger | Purpose |
|---|---|---|
| `ci.yml` | PR → main | golangci-lint, build, DB setup, migrations, go test |
| `playwright_ci.yml` | PR → main | Playwright smoke tests |
| `security_scan.yml` | Daily + manual | Gosec, uploads SARIF to GitHub Security |
| `release.yml` | Push → main | Bumps semver tag from Conventional Commits |
| `delivery.yml` | Push → main | Builds multi-stage Docker images, pushes to GHCR |
| `deployment.yml` | After Delivery | SSH deploy to production server |
| `deploy_staging.yml` | Manual / staging trigger | SSH deploy to staging server |
| `playwright_staging.yml` | After staging deploy | Playwright tests against staging environment |
| `build-and-push-migration.yml` | Push → main | Builds and pushes migration Docker image to GHCR |

## Getting Started

> **Before making your first commit, read [CONTRIBUTING.md](CONTRIBUTING.md).** It covers pre-commit hooks, commit conventions, branch naming, and the full dev/prod workflow.

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- `make`
- Go 1.25
- A bash-compatible terminal

### Running locally

```bash
# Clone the repo
git clone https://github.com/Knowledgeables/Knowledge.git
cd Knowledge/knowledgeable

# Start the development environment (Docker Compose + live reload)
make dev-up

# App is available at http://localhost:8080
# Database is seeded automatically in dev mode

# Stop
make dev-down
```

### Install pre-commit hooks

```bash
# Install prek (macOS / Linux)
brew install prek
# or
pip install prek

# Wire up hooks
make setup-hooks
```

The hooks run automatically before every commit and check:

1. Commit size — warns if the diff exceeds 400 LOC
2. `golangci-lint` — static analysis
3. `gosec` — security scanner
4. Commit message format — enforces Conventional Commits
5. `go test ./...` — full test suite

### Other useful commands

```bash
make lint        # Run golangci-lint manually
make go-test     # Run all Go tests
make smoke-test  # Run Playwright smoke tests (requires running app)
make full-e2e    # Run full Playwright integration tests (staging)
make swagger     # Regenerate Swagger docs
```

## Project Structure

```
Knowledge/
├── .github/
│   ├── workflows/          # 8+ CI/CD pipelines
│   └── ISSUE_TEMPLATE/     # Bug, feature, chore, docs templates
├── docs/
│   ├── legacy-analysis.md     # Issues found in the old system
│   └── legacy-architecture.md # Old monolith architecture notes
├── knowledgeable/
│   ├── cmd/server/main.go  # Application entry point
│   ├── internal/
│   │   ├── auth/           # Login / session handling
│   │   ├── db/             # PostgreSQL init & connection
│   │   ├── pages/          # Page search — handler/service/repo
│   │   ├── users/          # User management — handler/service/repo
│   │   └── web/            # Router and middleware
│   ├── db-migrations/      # Knex.js migrations & seeds
│   ├── e2e/
│   │   ├── smoke/          # Playwright smoke tests
│   │   └── full/           # Playwright full integration tests
│   ├── frontend/           # Tailwind source
│   ├── templates/          # HTML templates
│   ├── static/             # Compiled CSS + JS
│   ├── monitoring/         # Prometheus config
│   ├── nginx/              # Nginx reverse proxy config
│   ├── Dockerfile.dev
│   ├── Dockerfile.prod
│   ├── Dockerfile.migrations
│   ├── docker-compose-dev.yml
│   ├── docker-compose-prod.yml
│   ├── docker-compose-staging.yml
│   ├── docker-compose-build.yml
│   ├── knowledge.sql       # Legacy schema (replaced by Knex.js migrations)
│   ├── Makefile
│   └── playwright.config.js
├── CONTRIBUTING.md
└── README.md
```

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) **before** opening a PR. Key points:

- Every change must be linked to an issue
- Branch naming: `type/scope-keywords` (e.g. `feat/user-registration`)
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) — this drives automatic versioning
- No direct pushes to `main`
- Keep PRs small and focused

## Legacy vs. Rewrite

| | Legacy (Flash / Python 2) | Knowledgeable (Go 1.25) |
|---|---|---|
| Passwords | MD5 (no salt) | bcrypt |
| SQL | Raw string interpolation | Parameterised queries |
| Secrets | Hardcoded in source | Environment variables |
| Tests | None | Unit + E2E (Playwright) |
| Deployments | Manual | Automated via GitHub Actions |
| Security scanning | None | Gosec (daily, SARIF) |
