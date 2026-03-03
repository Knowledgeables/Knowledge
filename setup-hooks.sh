#!/bin/sh
# setup-hooks.sh
# Run this script once from the repo root to install the golangci-lint pre-commit hook.
# Usage: sh setup-hooks.sh

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"
HOOK_FILE="$HOOKS_DIR/pre-commit"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SOURCE_HOOK="$SCRIPT_DIR/pre-commit"

# Verify we're inside a git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo "ERROR: Not inside a git repository. Run this from your project root."
  exit 1
fi

# Verify the knowledgeable project folder exists
if [ ! -d "$REPO_ROOT/knowledgeable" ]; then
  echo "ERROR: Could not find the knowledgeable project folder at $REPO_ROOT/knowledgeable"
  exit 1
fi

# Warn if golangci-lint is not installed
if ! command -v golangci-lint > /dev/null 2>&1; then
  echo "WARNING: golangci-lint is not installed."
  echo ""
  echo "Install it before committing:"
  echo "  https://golangci-lint.run/usage/install/"
  echo ""
  echo "Quick install (Linux/macOS):"
  echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
  echo ""
  echo "Windows (via Chocolatey):"
  echo "  choco install golangci-lint"
  echo ""
fi

# Copy the pre-commit hook into .git/hooks/
if [ ! -f "$SOURCE_HOOK" ]; then
  echo "ERROR: pre-commit hook file not found at $SOURCE_HOOK"
  echo "Make sure 'pre-commit' is in the same directory as this script."
  exit 1
fi

cp "$SOURCE_HOOK" "$HOOK_FILE"
chmod +x "$HOOK_FILE"

echo "Pre-commit hook installed at $HOOK_FILE"
echo "golangci-lint will run in ./knowledgeable on every commit."
