# Contributing


## Set up local work environment with Docker compose

When cloning the project, in the knowledge directory, you will see two docker files and two docker compose in the main directory. These are used to seperate **development** and **production** environments.

To start the development environment (requires docker desktop), run:

'''bash
docker compose -f docker-compose-dev.yml up -d
'''

Now you have access to the development environment used by all the developers on the team.

## Before You Start Comitting

Before making your first commit, install the Git hooks by running this from the repo root:

```sh
sh setup-hooks.sh
```

This sets up a pre-commit hook that runs `golangci-lint` on the `knowledgeable` project automatically. Your commit will be blocked if there are any lint errors — fix them before committing.

> **Requires `golangci-lint` to be installed.**
> Install it from: https://golangci-lint.run/usage/install/
> Windows: `choco install golangci-lint`
> Mac: `brew install golangci-lint`
---

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

type(scope): present verb + short summary

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


