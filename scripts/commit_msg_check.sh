#!/bin/sh
# commit-msg hook: validates commit message follows the project's commit convention.
#
# Format:  <type>(<scope>): <past-tense verb> + short summary
# Scope:   required for feat, fix, chore, refactor, test — optional for docs
#
# Examples:
#   feat(user): implement registration
#   fix(auth): resolve token expiry bug
#   docs: update API documentation
#   chore(deps): bump golangci-lint to v1.57

COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Allowed types (and whether scope is required)
TYPES_WITH_SCOPE="feat|fix|chore|refactor|test"
TYPES_WITHOUT_SCOPE="docs"
ALL_TYPES="$TYPES_WITH_SCOPE|$TYPES_WITHOUT_SCOPE"

# Past-tense verb starters — description must begin with one of these
PAST_TENSE_VERBS="added|fixed|removed|updated|implemented|refactored|resolved|improved|changed|renamed|moved|replaced|bumped|started|created|deleted|migrated|extracted|merged|reverted"

# Strip lines starting with '#' (git comment lines) before validating
STRIPPED_MSG=$(echo "$COMMIT_MSG" | grep -v '^#')

# Get the first non-empty line as the subject
SUBJECT=$(echo "$STRIPPED_MSG" | awk 'NF{print; exit}')

if [ -z "$SUBJECT" ]; then
  echo ""
  echo "ERROR: Commit message is empty."
  echo ""
  exit 1
fi

# Rule 1: types that require a scope — type(scope): description
PATTERN_WITH_SCOPE="^($TYPES_WITH_SCOPE)\([a-zA-Z0-9_-]+\): ($PAST_TENSE_VERBS) .+[^\.]$"

# Rule 2: docs — scope is optional
PATTERN_DOCS="^docs(\([a-zA-Z0-9_-]+\))?: ($PAST_TENSE_VERBS) .+[^\.]$"

if ! echo "$SUBJECT" | grep -qE "$PATTERN_WITH_SCOPE" && \
   ! echo "$SUBJECT" | grep -qE "$PATTERN_DOCS"; then

  # Identify which part failed for a more helpful error message
  TYPE=$(echo "$SUBJECT" | sed -n 's/^\([a-zA-Z]*\).*/\1/p')

  echo ""
  echo "ERROR: Commit message does not follow the project convention."
  echo ""
  echo "  Your message:    $SUBJECT"
  echo ""
  echo "  Format:          <type>(<scope>): <past-tense verb> <summary>"
  echo "  Scope:           required for feat, fix, chore, refactor, test"
  echo "                   optional for docs"
  echo ""
  echo "  Allowed types:   feat, fix, chore, docs, refactor, test"
  echo "  Verb must be past tense, e.g.:"
  echo "    added / fixed / removed / updated / implemented / refactored"
  echo "    resolved / improved / changed / renamed / moved / bumped"
  echo ""
  echo "  Examples:"
  echo "    feat(user): implemented registration"
  echo "    fix(auth): resolved token expiry bug"
  echo "    chore(deps): bumped golangci-lint to v1.57"
  echo "    docs: updated API documentation"
  echo "    refactor(payment): extracted billing service"
  echo "    test(user): added unit tests for registration"
  echo ""
  exit 1
fi

echo "Commit message OK."
exit 0
