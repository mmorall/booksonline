# BooksOnline API

A modular monolith backend built in Go using Hexagonal Architecture and DDD principles, handling product browsing, order processing, and automatic digital asset generation for licenses and vouchers.

## Quick Start

### Prerequisites
- Go 1.22+
- Docker and Docker Compose

### Run Locally

```bash
cp .env.example .env
make db-up
make dev
```

The server starts on port `8080`. API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

### Test & Lint

```bash
make test   # runs tests with race detection and coverage
make lint   # runs golangci-lint
```

## API Reference

Full interactive documentation is available via Swagger UI at `/swagger/index.html`.

For admin endpoints (`GET /orders`, `GET /orders/{id}`), use HTTP Basic Auth with the credentials configured in your `.env` file.

## Architecture

The service is structured as a modular monolith with strict boundaries between domains. Each module (`catalog`, `orders`) owns its own domain model, ports, and adapters. Modules communicate through interfaces — never by direct import.

```
cmd/server/        → entry point, dependency wiring
internal/
  catalog/         → product domain
  orders/          → order domain
  shared/          → logger, telemetry, event bus
config/            → environment-based configuration
migrations/        → embedded SQL migrations, run on startup
k8s/               → Kubernetes manifests (ArgoCD, CNPG, HPA)
```

## Deployment

The application is deployed to a K3s cluster via ArgoCD. Pushes to `main` are automatically synced to the cluster.

The Horizontal Pod Autoscaler is configured for instant scale-up to handle peak load without manual intervention.