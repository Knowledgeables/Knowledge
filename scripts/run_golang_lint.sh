#!/bin/sh
# Runs golangci-lint inside the knowledgeable Go project

REPO_ROOT="$(git rev-parse --show-toplevel)"
PROJECT_DIR="$REPO_ROOT/knowledgeable"

echo "Running golangci-lint in $PROJECT_DIR ..."

if ! command -v golangci-lint > /dev/null 2>&1; then
  echo "ERROR: golangci-lint is not installed or not in PATH."
  echo "Install it from: https://golangci-lint.run/usage/install/"
  exit 1
fi

cd "$PROJECT_DIR" || {
  echo "ERROR: Could not cd into $PROJECT_DIR"
  exit 1
}

golangci-lint run ./...

EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
  echo ""
  echo "Lint failed. Please fix the issues above before committing."
  exit 1
fi

echo "Lint passed."
exit 0