#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Default values
BACKUP_NAME=$1
BACKUP_DIR=${BACKUP_DIR:-"./backups"}

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-bartender}
DB_PASSWORD=${DB_PASSWORD:-bartenderpass}
DB_NAME=${DB_NAME:-bartenderdb}

# Check if backup name provided
if [ -z "$BACKUP_NAME" ]; then
  echo -e "${RED}Error: Backup name is required.${NC}"
  echo -e "Usage: $0 <backup-name>"
  
  # List available backups
  if [ -d "$BACKUP_DIR" ] && [ "$(ls -A $BACKUP_DIR)" ]; then
    echo -e "\n${YELLOW}Available backups:${NC}"
    ls -1 $BACKUP_DIR
  else
    echo -e "\n${YELLOW}No backups found in $BACKUP_DIR${NC}"
  fi
  
  exit 1
fi

# Check if backup file exists
BACKUP_FILE="${BACKUP_DIR}/${BACKUP_NAME}"
if [ ! -f "$BACKUP_FILE" ]; then
  echo -e "${RED}Error: Backup file '$BACKUP_FILE' not found.${NC}"
  exit 1
fi

echo -e "${GREEN}Restoring database from backup: $BACKUP_NAME${NC}"

# Wait for the PostgreSQL database to be ready
echo "Waiting for PostgreSQL to be ready..."
export PGPASSWORD=$DB_PASSWORD
until psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing database restore"

# Drop and recreate the database
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"

# Restore the database from backup
pg_restore -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -v "$BACKUP_FILE"

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Database restore completed successfully!${NC}"
else
  echo -e "${YELLOW}Database restore completed with warnings or errors.${NC}"
  exit 1
fi 