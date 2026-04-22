#!/bin/bash
set -e

echo "[TEARDOWN] Stopping staging environment"
docker compose -p staging -f ~/app/docker-compose-staging.yml down -v
echo "[TEARDOWN] Staging environment stopped"