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
- Go version 1.24 installed
- Make commands must run in a bash-compatible terminal

Now you have access to the development environment used by all the developers on the team.

## Before You Start Committing

Before making your first commit, install the framework prek, to run automated checks before commits

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

This commands installs the 3 hooks for you inside /.git/hooks and then checks 3 things. 
- Checks if the commit is over 400 lines of code and warns you if it is
- Runs golang-lint before committing, to make sure the code you provided is bulletproof
- Checks the commit message, to make sure you are following the provided conventions

## Dev Flow

Start the development environment from the `knowledgeable/` directory:

    make dev-up

This builds from `Dockerfile.dev`, mounts the entire source tree into the
container, and starts the app on http://localhost:8080.

The database is created automatically at `./data/dev.db` on first run.
With `APP_ENV=dev`, the seed data from `seed-dev.sql` is also applied.

To stop:

    make dev-down



## Prod Flow

Production runs from pre-built images published to the GitHub Container Registry.
No build step happens on the server.

On the server, from the `knowledgeable/` directory:

    docker compose -f docker-compose-prod.yml pull
    docker compose -f docker-compose-prod.yml up -d

This starts two containers:
- **app** — the Go server on internal port 8080
- **nginx** — reverse proxy, exposed on port 80

The database file lives at `./data/prod.db` on the host, mounted into
the container. It persists across deploys and restarts.

To deploy a new version, run the same two commands again. The old container
is replaced; the data volume is untouched.


## Docker Files

| File                      | Purpose                                      |
|---------------------------|----------------------------------------------|
| `Dockerfile.dev`          | Dev image — mounts source, supports live reload |
| `Dockerfile.prod`         | Multi-stage prod build — outputs minimal image  |
| `docker-compose-dev.yml`  | Local development                            |
| `docker-compose-prod.yml` | Production — uses published images           |
| `docker-compose-build.yml`| Builds and tags images for GHCR              |

To build and push new images (CI does this automatically on merge to main):

    docker compose -f docker-compose-build.yml build
    docker compose -f docker-compose-build.yml push


## Database Migrations

Schema initialization runs automatically on startup in **both dev and prod**.

When the server starts, it runs `knowledge.sql` against the database,
creating all tables if they do not exist. No manual step is required —
a fresh database is fully initialized on first boot.

In dev (`APP_ENV=dev`), `seed-dev.sql` is also run on startup to populate
test data. This does not run in production.

If you change the schema:
1. Update `knowledge.sql` with the new table/column definitions
2. For existing databases, write a manual `ALTER TABLE` or coordinate
   a data migration with the team (a migration library is tracked in #81)



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


