#!/bin/bash
set -e
echo "[DEPLOY] Starting staging deploy"
cd ~/app
echo "[DEPLOY] Logging into registry"
echo "[DEPLOY] Pulling images"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml --profile tools pull
echo "[DEPLOY] Start DB"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml up -d db
echo "[DEPLOY] Wait for DB"
DB_READY=0
for i in {1..30}; do
  if docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml exec -T db pg_isready -U user -d knowledge_staging; then
    echo "[DEPLOY] DB ready"
    DB_READY=1
    break
  fi
  echo "[DEPLOY] Waiting for DB... attempt $i/30"
  sleep 5
done
if [ "$DB_READY" -eq 0 ]; then
  echo "[DEPLOY] ERROR: DB never became ready"
  exit 1
fi
echo "[DEPLOY] Run migrations"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml --profile tools run --rm migrations migrate:latest

echo "[DEPLOY] Run seed"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml --profile tools run --rm migrations seed:run --specific=staging_seed.js
echo "[DEPLOY] Start app"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml up -d --remove-orphans app
echo "[DEPLOY] Checking containers"
docker compose --env-file .env.staging -p staging -f docker-compose-staging.yml ps
echo "[DEPLOY] Staging deploy completed"