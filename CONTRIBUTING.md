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


