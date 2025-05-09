#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Seeding database for Bartender App...${NC}"

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-bartender}
DB_PASSWORD=${DB_PASSWORD:-bartenderpass}
DB_NAME=${DB_NAME:-bartenderdb}

# Check if the database is ready
export PGPASSWORD=$DB_PASSWORD
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
  echo -e "${RED}PostgreSQL is not available. Please make sure the database is running.${NC}"
  exit 1
fi

echo "PostgreSQL is available - seeding database"

# Run the seed migration
migrate -path ./services/pkg/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up 2

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Database seeding completed successfully!${NC}"
else
  echo -e "${RED}Database seeding failed!${NC}"
  exit 1
fi 