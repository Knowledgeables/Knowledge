# Legacy System Analysis

## Purpose
This document describes identified issues in the legacy system and the rewrite plan.



Legacy Issues Identified
## SQL Injection Risk
User input is written directly into SQL queries.
This makes the system vulnerable to SQL injection attacks.

# Rewrite strategy:
Use parameterized queries or prepared statements when implementing Go.

## Insecure Password Hashing
Passwords are hashed using MD5.
MD5 is cryptographically broken and vulnerable to attacks.

# Rewrite strategy:
Use bcrypt


## Poor Error Handling (Database failure)
Application terminates if database connection fails.
No structured logging or shutdown.

# Rewrite strategy:
Implement proper error handling and structured logging.

## Hardcoded Admin User
Admin user is inserted directly in the code.
Security risk.
# Rewrite strategy:
Create  initialization logic with if statements in a config file using environmental values.

## Hardcoded SECRET_KEY
Secret stored directly in source code.

# Rewrite strategy:
Load secrets from environmental variables.

## Test Coverage Gaps
Integration tests only validate response body.
HTTP status codes are not tested.

## Code Quality Issues
Unused imports
Datetime is imported but never used.
Per_page=30 is never used.
Unused variables
Backup files committed
Makefile uses Python 2
README outdated.
