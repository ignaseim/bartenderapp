# Bartender App

A production-ready application for bar inventory management, cocktail recipes, and ordering systems.

## Architecture Overview

The Bartender App is built as a microservices architecture with:

- **Frontend**: React 18 SPA (TypeScript, Vite)
- **Backend Services**:
  - `auth-svc`: Authentication, JWT and RBAC management
  - `order-svc`: Order processing and management
  - `inventory-svc`: Ingredient and recipe management
  - `pricing-svc`: Cost calculations and price management
- **Database**: PostgreSQL 17 with primary + read replica
- **Infrastructure**: Kubernetes, Terraform with automatic teardown after 24h
- **Messaging**: NATS for async event handling
- **Observability**: Prometheus + Grafana, structured JSON logs

## Repository Structure

```
bartenderapp/
├── frontend/              # React 18 SPA with TypeScript
├── services/              # Go microservices
│   ├── auth/              # Authentication service
│   ├── inventory/         # Inventory management service
│   ├── order/             # Order processing service
│   ├── pricing/           # Pricing calculation service
│   └── pkg/               # Shared packages
├── infra/                 # Terraform IaC
│   ├── modules/           # Reusable Terraform modules
│   └── environments/      # Environment-specific configurations
├── k8s/                   # Kubernetes configurations
├── pipelines/             # CI/CD Azure DevOps YAML files
└── scripts/               # Utility scripts for development and operations
```

## Local Development

### Prerequisites

- Docker and Docker Compose
- Go 1.22+
- Node.js 18+
- PostgreSQL client tools

### Setup Local Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bartenderapp.git
   cd bartenderapp
   ```

2. Start the local development environment:
   ```bash
   docker-compose up -d
   ```

3. Run database migrations:
   ```bash
   scripts/run-migrations.sh
   ```

4. Optionally seed the database with sample data:
   ```bash
   scripts/seed-database.sh
   ```

5. Frontend development:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

6. Backend services development:
   ```bash
   cd services/[service-name]
   go run cmd/main.go
   ```

## Build and Deploy

### Building Docker Images

```bash
./scripts/build-images.sh
```

### Deploying to Kubernetes

```bash
# Apply Terraform infrastructure
cd infra/environments/dev
terraform init
terraform apply

# Apply Kubernetes configurations
kubectl apply -f k8s/
```

## Teardown and Database Backup

The infrastructure is configured to automatically destroy after 24 hours of uptime. Before teardown:

1. A PostgreSQL database dump is created and exported to object storage
2. The Terraform state file is committed to the repository

To manually trigger teardown:

```bash
cd infra/environments/dev
terraform destroy
```

To restore from a backup:

```bash
./scripts/restore-database.sh [backup-name]
```

## Testing

### Running Tests

```bash
# Backend services
cd services
go test ./...

# Frontend
cd frontend
npm test
```

Current test coverage: >= 80%

## API Documentation

API documentation is available through OpenAPI:

- Auth Service: http://localhost:8081/swagger/
- Inventory Service: http://localhost:8082/swagger/
- Order Service: http://localhost:8083/swagger/
- Pricing Service: http://localhost:8084/swagger/

## User Roles and Permissions

- **Admin**: Full access to system. Can manage ingredients, recipes, users, and view reports.
- **Bartender**: Can indicate which cocktails they can make and process orders.
- **Guest**: Anonymous role that can browse cocktails and place orders. 