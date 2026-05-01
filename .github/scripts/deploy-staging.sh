#!/bin/bash
set -e
echo "[DEPLOY] Starting staging deploy"
cd ~/app
echo "[DEPLOY] Logging into registry"
echo "$CR_PAT" | docker login ghcr.io -u "$DOCKER_USERNAME" --password-stdin
echo "[DEPLOY] Pulling images"
docker compose -p staging -f docker-compose-staging.yml pull
echo "[DEPLOY] Start DB"
docker compose -p staging -f docker-compose-staging.yml up -d db
echo "[DEPLOY] Wait for DB"
DB_READY=0
for i in {1..30}; do
  if docker compose -p staging -f docker-compose-staging.yml exec -T db pg_isready -U postgres -d knowledgeable; then
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
echo "[DEPLOY] Run migrations + seed"
docker compose -p staging -f docker-compose-staging.yml \
  --profile tools run --rm migrations sh -c "knex migrate:latest && knex seed:run --specific=staging_seed.js"
  
echo "[DEPLOY] Start app"
docker compose -p staging -f docker-compose-staging.yml up -d --remove-orphans app
echo "[DEPLOY] Checking containers"
docker compose -p staging -f docker-compose-staging.yml ps
echo "[DEPLOY] Staging deploy completed"