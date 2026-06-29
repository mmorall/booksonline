# Technology Justification

Rationale behind the primary technology selection for the platform, prioritizing reliability, maintainability, and operational simplicity.

### Core Language: Go (Golang)
Go was selected for its high runtime concurrency performance, minimal memory footprint, and fast compilation speeds. Its strict type system and standard library design naturally support Hexagonal Architecture, allowing the implementation of domain business logic that remains completely decoupled from external infrastructure frameworks or database drivers.

### Database: PostgreSQL & CloudNativePG (CNPG)
PostgreSQL provides robust transactional consistency. CloudNativePG was chosen as the operator to manage the database cluster natively within Kubernetes. It abstracts operational complexities by automating backup scheduling, handling seamless database upgrades, and managing automated failovers without human intervention.

### Deployment & Orchestration: k3s & ArgoCD
- **k3s:** A highly optimized, lightweight Kubernetes distribution selected to minimize resource overhead on self-hosted environments while maintaining full CNCF compliance.
- **ArgoCD:** Utilized to enforce a strict GitOps deployment model. It continuously reconciles the cluster's actual running state with the desired state declared in Git, eliminating configuration drift and ensuring reproducible infrastructure definitions.