#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Root directory of the project
PROJECT_ROOT=$(dirname "$(dirname "$(readlink -f "$0")")")
cd "$PROJECT_ROOT"

echo -e "${BLUE}Building Docker images for Bartender App...${NC}"

# Build frontend
echo -e "${BLUE}Building frontend image...${NC}"
docker build -t bartenderapp/frontend:latest ./frontend

# Build backend services
SERVICES=("auth" "inventory" "order" "pricing")

for service in "${SERVICES[@]}"; do
  echo -e "${BLUE}Building ${service} service image...${NC}"
  docker build -t "bartenderapp/${service}-service:latest" -f "./services/${service}/Dockerfile" ./services
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}Successfully built ${service} service image.${NC}"
  else
    echo -e "${RED}Failed to build ${service} service image.${NC}"
    exit 1
  fi
done

echo -e "${GREEN}All Docker images built successfully!${NC}"

# List the built images
echo -e "${BLUE}Built images:${NC}"
docker images | grep bartenderapp 