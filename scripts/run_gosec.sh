#!/bin/sh
# Runs gosec security scanner inside the knowledgeable Go project

REPO_ROOT="$(git rev-parse --show-toplevel)"
PROJECT_DIR="$REPO_ROOT/knowledgeable"

echo "Running gosec in $PROJECT_DIR ..."

if ! command -v gosec > /dev/null 2>&1; then
  echo "ERROR: gosec is not installed or not in PATH."
  echo "Install it from: https://github.com/securego/gosec"
  echo "  go install github.com/securego/gosec/v2/cmd/gosec@latest"
  exit 1
fi

cd "$PROJECT_DIR" || {
  echo "ERROR: Could not cd into $PROJECT_DIR"
  exit 1
}

gosec -exclude=G706 ./...

EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
  echo ""
  echo "Gosec found security issues. Please fix them before committing."
  exit 1
fi

echo "Gosec passed."
exit 0
