#!/bin/bash
set -e

echo "[TEARDOWN] Stopping staging environment"
docker compose --env-file .env.staging -p staging -f ~/app/docker-compose-staging.yml down -v
echo "[TEARDOWN] Staging environment stopped"