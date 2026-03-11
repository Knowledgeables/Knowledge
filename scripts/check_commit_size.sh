#!/bin/sh

MAX_LINES=400

echo "Checking commit size..."

LINES_CHANGED=$(git diff --cached --numstat | awk '{added+=$1; removed+=$2} END {print added+removed}')

WARN=false

if [ "$LINES_CHANGED" -gt "$MAX_LINES" ]; then
  echo ""
  echo "⚠️ Large commit detected."
  echo "You are about to commit $LINES_CHANGED lines."
  echo "Recommended maximum is $MAX_LINES."
  WARN=true
fi

if [ "$WARN" = true ]; then
  echo ""
  printf "Continue commit? (y/N): "
  read confirm

  case "$confirm" in
    y|Y) echo "Continuing..." ;;
    *)
      echo "❌ Commit aborted."
      exit 1
      ;;
  esac
fi