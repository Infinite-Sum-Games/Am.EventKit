#!/bin/bash

PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

DB_NAME="postgres"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"
PGPASSWORD="password"

BASE_DIR="path/to/dir"
BACKUP_DIR="$BASE_DIR/backups"
LOG_FILE="$BACKUP_DIR/backup.log"

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_backup_$TIMESTAMP.dump"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] Starting PostgreSQL backup..." >> "$LOG_FILE"

/usr/bin/pg_dump \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -F c \
  -b \
  -v \
  -f "$BACKUP_FILE" \
  "$DB_NAME" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    echo "[$(date)] Backup successful: $BACKUP_FILE" >> "$LOG_FILE"
else
    echo "[$(date)] Backup FAILED!" >> "$LOG_FILE"
    exit 1
fi
