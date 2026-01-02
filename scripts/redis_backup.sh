#!/bin/bash

# Redis Backup Script
# Usage: ./redis_backup.sh

# --- Configuration ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/env.toml"

get_config_value() {
  local key=$1
  if [ -f "$ENV_FILE" ]; then
    grep "^${key}[[:space:]]*=" "$ENV_FILE" | cut -d'=' -f2- | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//"
  fi
}

REDIS_HOST=$(get_config_value "redis_host")
REDIS_PORT=$(get_config_value "redis_port")
REDIS_PASSWORD=$(get_config_value "redis_password")
REDIS_USERNAME=$(get_config_value "redis_username")

# Set defaults if env.toml is missing or values are empty
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"

# Directory where Redis saves its dump.rdb (Check your redis.conf 'dir' setting)
REDIS_DATA_DIR="${REDIS_DATA_DIR:-/var/lib/redis}"
# Directory where you want to store backups
BACKUP_DIR="./redis_backups"
RETENTION_DAYS=7
LOG_FILE="$BACKUP_DIR/backup.log"

# --- Setup ---
mkdir -p "$BACKUP_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILENAME="redis_dump_${TIMESTAMP}.rdb"

# --- Execution ---
echo "[$(date)] Starting Redis backup..." | tee -a "$LOG_FILE"

# Prepare redis-cli command
CLI_CMD="redis-cli -h $REDIS_HOST -p $REDIS_PORT"
if [ -n "$REDIS_USERNAME" ]; then
  CLI_CMD="$CLI_CMD --user $REDIS_USERNAME"
fi

if [ -n "$REDIS_PASSWORD" ]; then
  # Use REDISCLI_AUTH environment variable to avoid warning about password on CLI
  export REDISCLI_AUTH="$REDIS_PASSWORD"
fi

# 1. Trigger a synchronous SAVE
# NOTE: 'SAVE' blocks the Redis server until the dump is created.
# For high-traffic production systems, consider using 'BGSAVE' and waiting for it to complete.
echo "[$(date)] Sending SAVE command to Redis..." | tee -a "$LOG_FILE"
$CLI_CMD SAVE >>"$LOG_FILE" 2>&1
SAVE_EXIT_CODE=$?

if [ $SAVE_EXIT_CODE -ne 0 ]; then
  echo "[$(date)] ERROR: Redis SAVE command failed. Check host/password." | tee -a "$LOG_FILE"
  exit 1
fi

# 2. Copy the dump.rdb file
# Note: You might need sudo permissions to read from /var/lib/redis depending on your user
if [ -f "$REDIS_DATA_DIR/dump.rdb" ]; then
  echo "[$(date)] Copying dump.rdb..." | tee -a "$LOG_FILE"
  cp "$REDIS_DATA_DIR/dump.rdb" "$BACKUP_DIR/$BACKUP_FILENAME"

  if [ $? -eq 0 ]; then
    echo "[$(date)] Backup copied to $BACKUP_DIR/$BACKUP_FILENAME" | tee -a "$LOG_FILE"
  else
    echo "[$(date)] ERROR: Failed to copy dump.rdb. Check permissions." | tee -a "$LOG_FILE"
    exit 1
  fi
else
  echo "[$(date)] ERROR: dump.rdb not found in $REDIS_DATA_DIR. Check REDIS_DATA_DIR config." | tee -a "$LOG_FILE"
  exit 1
fi

# # 3. Compress the backup
# gzip "$BACKUP_DIR/$BACKUP_FILENAME"
# echo "[$(date)] Backup compressed to $BACKUP_DIR/$BACKUP_FILENAME.gz" | tee -a "$LOG_FILE"

# 4. Cleanup old backups
# echo "[$(date)] Cleaning up backups older than $RETENTION_DAYS days..." | tee -a "$LOG_FILE"
# find "$BACKUP_DIR" -name "redis_dump_*.rdb.gz" -mtime +$RETENTION_DAYS -delete
#
echo "[$(date)] Backup completed successfully." | tee -a "$LOG_FILE"
