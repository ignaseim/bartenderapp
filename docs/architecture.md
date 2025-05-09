# Bartender App Architecture

This document outlines the architecture of the Bartender App, a microservices-based application for managing bar operations.

## System Architecture

The Bartender App is built using a microservices architecture with the following components:

```mermaid
graph TD
    Client[Frontend SPA] --> Auth[Auth Service]
    Client --> Inventory[Inventory Service]
    Client --> Order[Order Service] 
    Client --> Pricing[Pricing Service]
    
    Auth --> DB[(PostgreSQL)]
    Inventory --> DB
    Order --> DB
    Pricing --> DB
    
    Auth -.->|Events| MQ[NATS Message Queue]
    Inventory -.->|Events| MQ
    Order -.->|Events| MQ
    Pricing -.->|Events| MQ
    
    subgraph Kubernetes Cluster
        Auth
        Inventory
        Order
        Pricing
        MQ
    end
    
    subgraph Database
        DB
    end
```

## Service Descriptions

1. **Auth Service**: Handles user authentication, authorization, and user management. Provides JWT tokens for secure API access.

2. **Inventory Service**: Manages ingredients, recipes, and stock levels. Tracks inventory movements and alerts for low stock.

3. **Order Service**: Processes customer orders, tracks order status, and manages the order lifecycle.

4. **Pricing Service**: Calculates cocktail costs and prices based on ingredient costs and pour-cost margins.

## Sequence Diagrams

### Authentication Flow

```mermaid
sequenceDiagram
    Client->>Auth Service: Login Request
    Auth Service->>Database: Validate Credentials
    Database-->>Auth Service: User Data
    Auth Service->>Auth Service: Generate JWT
    Auth Service-->>Client: Return JWT Token
    
    Client->>Protected Service: API Request + JWT
    Protected Service->>Auth Service: Validate Token
    Auth Service-->>Protected Service: Token Valid
    Protected Service->>Database: Process Request
    Database-->>Protected Service: Return Data
    Protected Service-->>Client: API Response
```

### Order Processing Flow

```mermaid
sequenceDiagram
    Client->>Order Service: Place Order
    Order Service->>Pricing Service: Calculate Prices
    Pricing Service->>Inventory Service: Get Ingredient Costs
    Inventory Service-->>Pricing Service: Ingredient Costs
    Pricing Service-->>Order Service: Total Price
    Order Service->>Database: Save Order
    Order Service->>NATS: Publish Order Event
    Order Service-->>Client: Order Confirmation
    
    Inventory Service->>NATS: Subscribe to Order Event
    NATS-->>Inventory Service: Order Event
    Inventory Service->>Database: Update Inventory Levels
```

## Deployment Architecture

The application is deployed on Kubernetes with the following infrastructure:

```mermaid
graph TB
    subgraph Cloud Provider
        subgraph Kubernetes Cluster
            subgraph Services
                Auth[Auth Service]
                Inventory[Inventory Service]
                Order[Order Service]
                Pricing[Pricing Service]
            end
            
            subgraph Messaging
                NATS[NATS JetStream]
            end
            
            subgraph Observability
                Prometheus[Prometheus]
                Grafana[Grafana]
            end
            
            subgraph Ingress
                Gateway[API Gateway]
            end
        end
        
        subgraph Database
            PG_Primary[PostgreSQL Primary]
            PG_Replica[PostgreSQL Replica]
            PG_Primary --> PG_Replica
        end
        
        subgraph Storage
            ObjectStorage[Object Storage]
        end
    end
    
    Client[Customers] --> Gateway
    Gateway --> Services
    Services --> PG_Primary
    Services -.-> NATS
    Prometheus --> Services
    Grafana --> Prometheus
```

## Data Model

The core data model is represented in the following entity-relationship diagram:

```mermaid
erDiagram
    USERS {
        int user_id PK
        string username
        string email
        string password_hash
        string role
        timestamp created_at
        timestamp updated_at
    }
    
    INGREDIENTS {
        int ingredient_id PK
        string name
        string category
        decimal package_size_ml
        int package_cost_cents
        timestamp created_at
        timestamp updated_at
    }
    
    INGREDIENT_STOCK {
        int ingredient_id PK,FK
        decimal qty_ml
        timestamp updated_at
    }
    
    RECIPES {
        int recipe_id PK
        string name
        string method
        string glass
        string garnish
        string instructions
        int created_by FK
        timestamp created_at
        timestamp updated_at
    }
    
    RECIPE_ITEMS {
        int recipe_id PK,FK
        int ingredient_id PK,FK
        decimal amount_ml
    }
    
    ORDERS {
        int order_id PK
        int customer_id FK
        int bartender_id FK
        string status
        timestamp created_at
        timestamp updated_at
    }
    
    ORDER_ITEMS {
        int order_id PK,FK
        int recipe_id PK,FK
        int quantity
        int price_cents
        string status
    }
    
    BARTENDER_SKILLS {
        int user_id PK,FK
        int recipe_id PK,FK
    }
    
    INVENTORY_TRANSACTIONS {
        int transaction_id PK
        int ingredient_id FK
        decimal quantity_ml
        string transaction_type
        int reference_id
        int created_by FK
        timestamp created_at
    }
    
    USERS ||--o{ RECIPES : creates
    USERS ||--o{ BARTENDER_SKILLS : has_skill
    USERS ||--o{ ORDERS : places
    USERS ||--o{ ORDERS : prepares
    USERS ||--o{ INVENTORY_TRANSACTIONS : records
    
    INGREDIENTS ||--|| INGREDIENT_STOCK : has_stock
    INGREDIENTS ||--o{ RECIPE_ITEMS : used_in
    INGREDIENTS ||--o{ INVENTORY_TRANSACTIONS : affects
    
    RECIPES ||--o{ RECIPE_ITEMS : contains
    RECIPES ||--o{ ORDER_ITEMS : ordered_as
    RECIPES ||--o{ BARTENDER_SKILLS : required_for
    
    ORDERS ||--o{ ORDER_ITEMS : includes
```

## Auto-Teardown Mechanism

The infrastructure is designed to automatically tear down after 24 hours of operation. The process works as follows:

```mermaid
sequenceDiagram
    participant TF as Terraform
    participant K8S as Kubernetes
    participant DB as Database
    participant S3 as Object Storage
    participant GIT as Git Repository
    
    TF->>K8S: Apply Infrastructure
    K8S->>DB: Create Database
    
    Note over TF,GIT: Infrastructure running for 24 hours
    
    TF->>DB: Trigger Export
    DB-->>S3: Export Database Dump
    TF->>TF: Generate State Snapshot
    TF-->>GIT: Commit State File
    TF->>K8S: Destroy Infrastructure
```

## Security Architecture

The application implements multiple security layers:

1. **Authentication**: JWT-based authentication with refresh tokens
2. **Authorization**: Role-based access control (RBAC)
3. **API Security**: TLS everywhere, rate limiting, input validation
4. **Data Security**: Encrypted database connections, password hashing
5. **Infrastructure Security**: Network policies, secrets management

## Conclusion

This architecture provides a scalable, maintainable, and secure foundation for the Bartender App. The microservices design allows for independent scaling and development of each component, while the infrastructure automation ensures efficient resource usage. 