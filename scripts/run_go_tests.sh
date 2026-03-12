#!/bin/sh
# Run Go tests directly inside knowledgeable/

REPO_ROOT="$(git rev-parse --show-toplevel)"
PROJECT_DIR="$REPO_ROOT/knowledgeable"

echo "Running Go tests in $PROJECT_DIR..."

cd "$PROJECT_DIR" || {
    echo "ERROR: Could not cd into $PROJECT_DIR"
    exit 1
}

if ! go test ./...; then
    echo ""
    echo "❌ Go tests failed. Commit aborted."
    exit 1
fi

echo "✅ Go tests passed."
exit 0