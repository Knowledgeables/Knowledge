#!/bin/sh
# setup-hooks.sh
# Configures the repo to use the version-controlled githooks directory

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
HOOKS_DIR="$REPO_ROOT/githooks"

echo "Configuring Git hooks..."

# Tell git to use the githooks directory
git config core.hooksPath githooks

# Ensure hooks are executable
chmod +x "$HOOKS_DIR"/*

echo ""
echo "Hooks enabled from: $HOOKS_DIR"
echo ""
echo "Available hooks:"
ls "$HOOKS_DIR"