#!/bin/bash
set -euo pipefail
CONTAINER="app-db-1"
DB="knowledge"
USER="user"
BACKUP_DIR="/opt/db_backups"
LOG_FILE="$BACKUP_DIR/db_backup.log"
AZURE_BASE_URL="https://whoknowsdbbackup.blob.core.windows.net/db-backups"
SAS_TOKEN="{{ azure_sas_token }}"
mkdir -p "$BACKUP_DIR"
exec > >(tee -a "$LOG_FILE") 2>&1
trap 'echo "BACKUP_FAILED ts=$(date -Is)"' ERR
echo "BACKUP_START ts=$(date -Is)"
FILE="$BACKUP_DIR/knowledge_$(date +%F_%H-%M).sql.gz"
BLOB_NAME="$(basename "$FILE")"
echo "Running pg_dump"
PGPASSWORD=$(docker exec "$CONTAINER" printenv POSTGRES_PASSWORD)
docker exec -e PGPASSWORD="$PGPASSWORD" "$CONTAINER" \
  pg_dump -U "$USER" "$DB" | gzip > "$FILE"
UPLOAD_URL="${AZURE_BASE_URL}/${BLOB_NAME}?${SAS_TOKEN}"
echo "Uploading to Azure: $BLOB_NAME"
for i in 1 2 3; do
  if azcopy copy "$FILE" "$UPLOAD_URL"; then
    echo "Upload success"
    break
  else
    echo "Upload failed (attempt $i)"
    sleep 5
  fi
  if [ "$i" -eq 3 ]; then
    echo "Upload failed after retries"
    exit 1
  fi
done
echo "Cleaning old backups"
find "$BACKUP_DIR" -type f -name "*.sql.gz" -mtime +3 -delete
echo "BACKUP_SUCCESS ts=$(date -Is) file=$BLOB_NAME"