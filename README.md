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
| Infrastructure (IaC) | Terraform (Azure) |
| Configuration management | Ansible |
| Metrics | Prometheus + Node Exporter |
| Logs | Loki + Promtail |
| Dashboards | Grafana |

## CI/CD Pipeline

Every change goes through the following pipeline before reaching production:

```
Pull Request → main
    ├── Go CI          (lint → build → DB setup → migrations → integration tests)
    └── Playwright     (smoke tests)

Push to main
    ├── Auto Release   (conventional commit → semver tag)
    ├── Delivery       (build & push Docker images to GHCR)
    │       └── Deploy Staging  (SSH → staging)
    │               └── Playwright Staging  (full E2E suite)
    │                       └── Deployment  (SSH → production, rollback on failure)
    └── Configure      (Ansible → install/update Docker, fail2ban, node_exporter, promtail on VMs)

Push to main (ansible/ changes)
    └── Configure      (Ansible configuration management run)

Monitoring server
    └── Monitoring Deployment  (SSH → deploy Prometheus + Loki + Grafana + Nginx stack)

Daily (03:00 UTC)
    ├── Security Scan  (Gosec → SARIF → GitHub Security tab)
    ├── Drift Check    (Ansible check mode → open GitHub issue on drift)
    └── Crawler        (Node.js crawler ingests data to Knowledge backend)

Scheduled (03:17 & 13:17 UTC)
    └── Crawler        (recurring content ingestion)
```

| Workflow file | Trigger | Purpose |
|---|---|---|
| `ci.yml` | PR → main | golangci-lint, build, DB setup, migrations, integration tests |
| `playwright_ci.yml` | PR → main | Playwright smoke tests |
| `security_scan.yml` | Daily + manual | Gosec, uploads SARIF to GitHub Security |
| `release.yml` | Push → main | Bumps semver tag from Conventional Commits |
| `delivery.yml` | Push → main | Builds multi-arch Docker images, pushes to GHCR |
| `deploy_staging.yml` | After Delivery | SSH deploy to staging server |
| `playwright_staging.yml` | After staging deploy | Full Playwright suite against staging environment |
| `deployment.yml` | After Playwright staging | SSH deploy to production; rollback on failure |
| `configure.yml` | Push → main (ansible/) + manual | Ansible: provisions Docker, fail2ban, node_exporter, promtail on all VMs |
| `monitoring-deployment.yml` | Manual | SSH deploy of monitoring stack (Prometheus, Loki, Grafana, Nginx) |
| `drift-check.yml` | Daily 03:00 UTC | Ansible check mode; opens GitHub issue if configuration drift detected |
| `crawler.yml` | Scheduled (03:17 & 13:17 UTC) + manual | Node.js web crawler ingests content to Knowledge backend |

## Infrastructure

Cloud infrastructure is defined as code in the [`Terraform/`](Terraform/) directory and targets Azure (Norway East).

| Resource | Details |
|---|---|
| Provider | Azure (`azurerm` v4.27.0) |
| VMs | 2 × `Standard_B2ats_v2` Ubuntu 22.04 LTS — `knowledge` (app) and `kmonitor` (monitoring) |
| Network | VNet `10.0.0.0/16`, subnet `10.0.2.0/24` |
| Security | NSG rules: SSH, HTTP/HTTPS, Node Exporter (subnet-only, port 9100), staging port 8081 |
| Outputs | Public IPs, SSH commands, ready-to-paste Ansible inventory and SSH config |

```bash
cd Terraform/
terraform init
terraform apply          # provisions VMs, network, security groups
terraform output         # prints IPs and SSH config
```

## Configuration Management

Server state is managed with Ansible from the [`ansible/`](ansible/) directory. The `configure.yml` GitHub Actions workflow runs the playbook automatically on pushes that touch `ansible/`, and daily via `drift-check.yml` (check mode only).

**What the playbook installs on every VM:**

- Docker + Docker Compose v2 plugin
- `fail2ban` (SSH brute-force protection)
- Node Exporter v1.8.2 (metrics, port 9100)

**App server additionally gets:**

- Promtail (ships container + system logs to Loki)
- PostgreSQL backup cron (hourly, uploads to Azure Blob via `azcopy`)

**Monitoring server additionally gets:**

- Loki backup cron (daily 03:00 UTC, uploads to Azure Blob)
- Blackbox Exporter config (deployed via Ansible, used for HTTP uptime probing)

```bash
cd ansible/
ansible-playbook playbook.yml              # full run
ansible-playbook playbook.yml --check      # drift detection (no changes)
```

## Monitoring

The monitoring stack lives in [`monitoring/`](monitoring/) and is deployed to the dedicated `kmonitor` VM via the `monitoring-deployment.yml` workflow.

| Service | Role |
|---|---|
| Prometheus | Scrapes metrics from both VMs via Node Exporter (30 s interval, 15-day retention) |
| Loki | Log aggregation backend (720 h retention, boltdb-shipper) |
| Promtail | Runs on the app VM; ships Docker container logs and system logs to Loki |
| Grafana | Dashboards for application logs, search performance/SLO, user activity, and infrastructure |
| Nginx | Reverse proxy on port 80 — routes `/loki/api/*` to Loki and `/` to Grafana |
| Blackbox Exporter | HTTP uptime probing — Prometheus scrapes it to track availability of the production `/health` endpoint |

Pre-provisioned Grafana dashboards:

- **Infrastructure** — VM CPU, memory, disk, network via Node Exporter
- **Search Performance & Quality** — latency and result quality metrics
- **Search SLO (Health)** — error-budget burn tracking
- **Container Activiy**  - shows active containers and log volume
- **User Activity (Realtime)** — live registration and login events
- **Application Logs** — structured log explorer via Loki
- **Forgot Password / Reset** — password-reset funnel

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

### Script Placement

Use this convention when adding scripts:

- `scripts/` (repo root): repository-level automation used by hooks and CI helpers
- `knowledgeable/scripts/`: application runtime scripts used by the Knowledgeable service itself

Examples:

- Commit hooks and quality checks stay in `scripts/`
- The crawler stays in `knowledgeable/scripts/`

## Project Structure

```
Knowledge/
├── .github/
│   ├── workflows/          # 13 CI/CD pipelines (app, infra, monitoring, security)
│   ├── scripts/            # deploy-staging.sh, teardown-staging.sh
│   └── ISSUE_TEMPLATE/     # Bug, feature, chore, docs templates
├── Terraform/              # Azure IaC — VMs, network, NSG, outputs
│   ├── main.tf             # Provider + resource group
│   ├── vms.tf              # 2 Linux VMs (knowledge, kmonitor)
│   ├── network.tf          # VNet, subnet, NSG rules
│   ├── variables.tf        # Input variables (VM size, location, etc.)
│   ├── outputs.tf          # IPs, SSH commands, Ansible inventory
│   └── terraform.tfvars    # Azure subscription + region config
├── ansible/                # Configuration management
│   ├── playbook.yml        # Main playbook (bootstraps all VMs)
│   ├── ansible.cfg
│   ├── inventory.ini       # Static inventory (knowledge + kmonitor)
│   ├── group_vars/all.yml  # Shared variables (node_exporter version, Loki URL)
│   ├── tasks/              # packages, docker, node_exporter, promtail, fail2ban, backups
│   ├── templates/          # db_backup.sh.j2, loki_backup.sh.j2, promtail-config.yml.j2
│   └── files/              # promtail-docker-compose.yml
├── monitoring/             # Monitoring stack (deployed to kmonitor VM)
│   ├── docker-compose.yml  # Loki, Prometheus, Grafana, Nginx, Blackbox Exporter
│   ├── prometheus/         # prometheus.yml (scrapes both VMs + blackbox HTTP probes)
│   ├── loki/               # config.yaml
│   ├── blackbox/           # blackbox_config.yml (HTTP uptime probe modules)
│   ├── grafana/
│   │   └── provisioning/
│   │       ├── datasources/    # Prometheus + Loki datasource configs
│   │       └── dashboards/     # 6 pre-provisioned dashboard JSON files
│   ├── nginx/              # nginx.conf (reverse proxy)
│   └── .env.example
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
| Infrastructure | Manually provisioned | Terraform (Azure IaC) |
| Server config | Ad-hoc / undocumented | Ansible (idempotent, drift-checked daily) |
| Observability | None | Prometheus + Loki + Grafana + 6 dashboards |
| Log shipping | None | Promtail → Loki |
| Backups | None | Hourly DB + daily Loki backups to Azure Blob |
