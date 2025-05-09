#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running database migrations for Bartender App...${NC}"

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-bartender}
DB_PASSWORD=${DB_PASSWORD:-bartenderpass}
DB_NAME=${DB_NAME:-bartenderdb}

# Wait for the PostgreSQL database to be ready
echo "Waiting for PostgreSQL to be ready..."
export PGPASSWORD=$DB_PASSWORD
until psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing migrations"

# Run the migrations
# Here we use the golang-migrate/migrate binary, but this could be replaced by goose or other migration tools
migrate -path ./services/pkg/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Migrations completed successfully!${NC}"
else
  echo -e "${RED}Migration failed!${NC}"
  exit 1
fi 