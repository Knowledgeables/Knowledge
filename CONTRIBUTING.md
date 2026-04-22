# Contributing

## Set up local work environment with Docker compose

After cloning the repository, navigate to the project root "knowledgeable".

To start the development environment, run:

```
 make dev-up
```
This command starts the development environment using Docker Compose.

To stop the environment again:

```
 make dev-down
```

The Makefile wraps the Docker Compose commands used by the project, so developers do not need to remember the full Docker commands.

Requirements:

- Docker Desktop installed
- make installed
- Go version 1.25 installed
- Make commands must run in a bash-compatible terminal

Now you have access to the development environment used by all the developers on the team.

## Before You Start Committing

Before making your first commit, install the framework prek, to run automated checks before commits.

```
brew install prek
pip install prek
```
See the direct documentation from the developers in case of confusion
https://github.com/j178/prek?tab=readme-ov-file#quick-start

This downloads the prek framework and allows you to run the following command

```
make setup-hooks
```

This command installs 5 hooks inside `/.git/hooks` that check:
- Commit size — warns if the diff exceeds 400 lines of code
- `golangci-lint` — static analysis to catch code issues before committing
- `gosec` — security scanner to catch vulnerabilities in new code
- Commit message format — enforces Conventional Commits
- `go test ./...` — full test suite must pass

## Environment Variables

The project uses a `.env` file inside `knowledgeable/` for database credentials.
This file is gitignored — you need to create it yourself before running anything.

Create `knowledgeable/.env` with the following content:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
DATABASE_URL=postgres://postgres:postgres@db:5432/knowledge?sslmode=disable
```

The Docker Compose files fall back to these same defaults if no `.env` is present,
so the values above will work out of the box for local development. Change them
if you need a different setup.

## Dev Flow

Start the development environment from the `knowledgeable/` directory:

    make dev-up

This builds from `Dockerfile.dev`, mounts the entire source tree into the
container, and starts the app on http://localhost:8080.

The PostgreSQL database is started automatically as a container. Knex.js
migrations run on startup via the migrations container, and dev seed data
from `db-migrations/seeds/dev_seed.js` is applied automatically.

To stop:

    make dev-down



## Prod Flow

Production runs from pre-built images published to the GitHub Container Registry.
No build step happens on the server.

On the server, from the `knowledgeable/` directory:

    docker compose -f docker-compose-prod.yml pull
    docker compose -f docker-compose-prod.yml up -d

This starts four containers:
- **postgres** — PostgreSQL 15 database, data persisted in a named volume
- **migrations** — runs Knex.js migrations then exits
- **app** — the Go server on internal port 8080
- **nginx** — reverse proxy, exposed on port 80

The database persists across deploys in a named Docker volume.

To deploy a new version, run the same two commands again. The old containers
are replaced; the data volume is untouched.


## Docker Files

| File                          | Purpose                                         |
|-------------------------------|-------------------------------------------------|
| `Dockerfile.dev`              | Dev image — mounts source, supports live reload |
| `Dockerfile.prod`             | Multi-stage prod build — outputs minimal image  |
| `Dockerfile.migrations`       | Runs Knex.js migrations against PostgreSQL      |
| `docker-compose-dev.yml`      | Local development                               |
| `docker-compose-prod.yml`     | Production — uses published images              |
| `docker-compose-staging.yml`  | Staging environment                             |
| `docker-compose-build.yml`    | Builds and tags images for GHCR                 |

To build and push new images (CI does this automatically on merge to main):

    docker compose -f docker-compose-build.yml build
    docker compose -f docker-compose-build.yml push


## Database Migrations

Schema changes are managed with **Knex.js** migrations in `db-migrations/`.
The `Dockerfile.migrations` container runs all pending migrations automatically
on every deploy — no manual step is required.

Migration files live in `db-migrations/migrations/` and seed files in
`db-migrations/seeds/`. In dev, `dev_seed.js` is applied automatically.
In staging/prod, no seed runs unless explicitly triggered.

If you need to change the schema:
1. Create a new migration file: `cd knowledgeable/db-migrations && npx knex migrate:make <name>`
2. Write the `up` (apply) and `down` (rollback) functions
3. The migration runs automatically on next deploy



## Branch Naming

type/scope-keywords

- feat/scope-short-keywords
- fix/scope-short-keywords
- chore/scope-short-keywords
- docs/scope-short-keywords
- refactor/scope-short-keywords

* lowercase only
* use - not _
* max 50 chars

Example:
feat/user-registration

---

## Commit Convention

We follow Conventional Commits:

type(scope): past tense verb + short summary

- feat(scope): added new feature
- fix(scope): made bug fix
- chore(scope): added maintenance
- docs: added documentation
- refactor(scope): started code refactor
- test(scope): implemented tests

Example:
feat(user): implemented registration

---

## Issues

- Every change must be linked to an issue
- Use correct issue type (feature, bug, chore, docs, epic)
- Add domain label (database, user, docs, etc.)
- Acceptance criteria required

---

## Pull Requests
- No direct pushes to main
- Link related issue (Closes #12)
- Keep PRs small and focused


