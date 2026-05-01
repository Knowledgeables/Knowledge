#!/bin/bash
set -e

cd ~/app
echo "[TEARDOWN] Stopping staging environment"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml down -v
echo "[TEARDOWN] Staging environment stopped"