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
Or online at:
```
https://api-booksonline.miguelmoral.com/swagger/index.html
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

## Resilience and Fault Tolerance

The production infrastructure is hardened against arbitrary node and runtime crashes using continuous automated fault injection via Chaos Mesh. 

### Chaos Experiment Profile
A background engine (`PodChaos`) runs a continuous loop that executes a disruptive `pod-kill` action against active application instances at a predefined regular interval (`@every 2m`).

### Automated Recovery Mechanism
1. **ReplicaSet Enforcement**: The Kubernetes Control Plane continuously evaluates cluster drift. When an instance is terminated by an infrastructure fault, a replacement instance is scheduled within milliseconds.
2. **Safe Traffic Routing**: Traffic is shielded during recovery via `readinessProbe` lifecycle endpoints (`GET /health`). The Traefik ingress proxy bypasses newly scheduled instances until they successfully initialize dependencies and establish an operational database connection pool, guaranteeing zero dropped HTTP requests during an active crash cycle.

## Deployment

The application is deployed to a K3s cluster via ArgoCD. Pushes to `main` are automatically synced to the cluster.

The Horizontal Pod Autoscaler is configured for instant scale-up to handle peak load without manual intervention.