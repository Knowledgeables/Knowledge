#!/bin/sh
# Runs golangci-lint inside the knowledgeable Go project

REQUIRED_VERSION="v2.11.4"
REPO_ROOT="$(git rev-parse --show-toplevel)"
PROJECT_DIR="$REPO_ROOT/knowledgeable"

echo "Running golangci-lint in $PROJECT_DIR ..."

if ! command -v golangci-lint > /dev/null 2>&1; then
  echo "ERROR: golangci-lint is not installed or not in PATH."
  echo "Install v${REQUIRED_VERSION#v}:"
  echo "   go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4  "
  exit 1
fi

INSTALLED_VERSION="v$(golangci-lint --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)"

if [ "$INSTALLED_VERSION" != "$REQUIRED_VERSION" ]; then
  echo "ERROR: golangci-lint $REQUIRED_VERSION required, but $INSTALLED_VERSION is installed."
  echo "Install the correct version:"
  echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b \$(go env GOPATH)/bin $REQUIRED_VERSION"
  exit 1
fi

cd "$PROJECT_DIR" || {
  echo "ERROR: Could not cd into $PROJECT_DIR"
  exit 1
}

golangci-lint run --allow-parallel-runners ./...

EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
  echo ""
  echo "Lint failed. Please fix the issues above before committing."
  exit 1
fi

echo "Lint passed."
exit 0
