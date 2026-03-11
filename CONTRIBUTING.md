# Contributing

## Set up local work environment with Docker compose

After cloning the repository, navigate to the project root.

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

## Before You Start Comitting

Before making your first commit, install the framework prek, to run automated checks before commits

```
brew install prek
pip install prek
```
See the direct documentation from the developers in case of confusion
https://github.com/j178/prek?tab=readme-ov-file#quick-start

This downloads the prek framework and allows you to run 

```
make setup-hooks
```


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
feat(user): implement registration

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


