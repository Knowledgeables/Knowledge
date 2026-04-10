# Knowledgeable

A knowledge management web application — rewritten from a legacy Flash/Python 2 monolith to a modern Go service as part of a school project exploring DevOps principles end-to-end.

[![Go CI](https://github.com/Knowledgeables/Knowledge/actions/workflows/ci.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/ci.yml)
[![Security Scan](https://github.com/Knowledgeables/Knowledge/actions/workflows/security_scan.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/security_scan.yml)
[![E2E Tests](https://github.com/Knowledgeables/Knowledge/actions/workflows/playwright.yml/badge.svg)](https://github.com/Knowledgeables/Knowledge/actions/workflows/playwright.yml)
[![Latest Release](https://img.shields.io/github/v/release/Knowledgeables/Knowledge)](https://github.com/Knowledgeables/Knowledge/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](https://go.dev/)

---

## Overview

Knowledgeable is a searchable knowledge-base platform supporting multiple languages (English & Danish). Users can register, log in, and search pages. The rewrite addresses every major flaw found in the legacy system — SQL injection, MD5 password hashing, hardcoded secrets, and zero test coverage — while introducing a full DevOps lifecycle: automated CI, container builds, security scanning, semantic versioning, and zero-downtime production deploys.

## Features

- User registration and authentication with `bcrypt` password hashing
- Full-text page search with multi-language support (EN / DA)
- Swagger/OpenAPI documentation (`/swagger/`)
- Docker-first development and production environments
- Automated semantic versioning driven by Conventional Commits

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Database | SQLite (`modernc.org/sqlite`) |
| Frontend | HTML templates + Tailwind CSS |
| E2E Testing | Playwright |
| Containerisation | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Security scanning | Gosec + GitHub SARIF |
| API docs | Swagger |

## CI/CD Pipeline

Every change goes through the following pipeline before reaching production:

```
Pull Request
    └── Go CI (lint → build → test)
            │
Push to main
    ├── Auto Release  (conventional commit → semver tag)
    ├── E2E Tests     (Playwright smoke tests)
    └── Delivery      (build & push Docker image to GHCR)
            └── Deployment (SSH → pull image → redeploy)

Daily (00:00 UTC)
    └── Security Scan (Gosec → SARIF → GitHub Security tab)
```

| Workflow | Trigger | Purpose |
|---|---|---|
| `Go CI` | PR → main | golangci-lint, build, go test |
| `Security scan` | Daily + manual | Gosec, uploads SARIF to GitHub Security |
| `E2E Tests` | Push / PR → main | Playwright smoke tests |
| `Auto Release` | Push → main | Bumps semver tag from commit type |
| `Delivery` | Push → main | Builds multi-stage Docker image, pushes to GHCR |
| `Deployment` | After Delivery | SSH deploy to production server |

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
make test        # Run all Go tests
make swagger     # Regenerate Swagger docs
```

## Project Structure

```
Knowledge/
├── .github/
│   ├── workflows/          # CI/CD pipelines
│   └── ISSUE_TEMPLATE/     # Bug, feature, chore, docs templates
├── docs/
│   ├── legacy-analysis.md     # Issues found in the old system
│   └── legacy-architecture.md # Old monolith architecture notes
├── knowledgeable/
│   ├── cmd/server/main.go  # Application entry point
│   ├── internal/
│   │   ├── auth/           # Login / session handling
│   │   ├── config/         # App configuration
│   │   ├── db/             # SQLite init & connection
│   │   ├── pages/          # Page search — handler/service/repo
│   │   ├── users/          # User management — handler/service/repo
│   │   └── web/            # Router and middleware
│   ├── frontend/           # Tailwind source
│   ├── templates/          # HTML templates
│   ├── tests/              # Integration & Playwright tests
│   ├── Dockerfile.dev
│   ├── Dockerfile.prod
│   ├── docker-compose-dev.yml
│   ├── docker-compose-prod.yml
│   ├── knowledge.sql       # Schema (auto-applied on startup)
│   └── Makefile
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
