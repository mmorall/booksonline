# BooksOnline API

A modular monolith backend service built in Go, implementing Hexagonal Architecture and DDD principles. The service provides functionalities for managing a product catalog and processing customer checkouts with automated digital token fulfillment.

## Core Features
- Hexagonal Architecture with complete separation between core domain and technical adapters.
- Multi-stage containerized build configuration with BuildKit caching optimizations.
- Embedded database migrations executed natively on application startup.
- High-performance horizontal pod autoscaling with custom instant scale-up behavior.

## API Endpoints

### Catalog Module
- `GET /products` - Retrieves the entire product list including books, digital licenses, and vouchers.
- `GET /products/{id}` - Retrieves detailed information for a specific product by its UUID.

### Orders Module
- `POST /orders` - Submits a new customer purchase order. For digital items (licenses and vouchers), cryptographic asset tokens are automatically generated upon successful creation.
- **Expected Request Payload:**
    ```json
    {
      "customer_email": "customer@example.com",
      "items": {
        "11111111-1111-1111-1111-111111111111": 2,
        "22222222-2222-2222-2222-222222222222": 1
      }
    }
    ```
- `GET /orders` - Lists all historic client orders (Administrative view).
- `GET /orders/{id}` - Retrieves a single order audit contract by its UUID, including all inner purchased line items and generated assets.
- `GET /health` - Returns the service operational availability status. Used by Kubernetes liveness and readiness probes.

---

## Quick Start Guide

### Prerequisites
- Go 1.22 or higher
- Docker and Docker Compose

### 1. Local Infrastructure Configuration
Create a local environment file by copying the template:
```bash
cp .env.example .env
```
Start the local decoupled PostgreSQL engine:
```bash
make db-up
```

### 2. Execution Run
Launch the application binary natively on your host machine. The runtime will automatically connect to the database container and execute all pending schema modifications:
```bash
make dev
```
The server listens on port 8080. You can test accessibility via standard HTTP clients:
```bash
curl http://localhost:8080/products
```

### 3. Verification Commands
Execute the testing suite with race condition detection and structural code coverage metrics:
```bash
make test
```
Execute static code analysis and formatting compliance policies:
```bash
make lint
```