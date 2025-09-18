# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

The `multicloud-operators-application` is a Kubernetes operator that provides application wrappers for selecting involved subscriptions and deployables within the Open Cluster Management (OCM) ecosystem. It implements a controller that manages Application resources that can discover and aggregate related subscription and deployable components.

## Key Architecture Components

### Core Components
- **Application Controller** (`pkg/controller/application/`): Main reconciliation logic for Application CRs
- **Webhook Validator** (`webhook/`): Admission webhook for Application validation
- **Utils Package** (`utils/`): Common utilities for event handling, exploration, and application management
- **Manager** (`cmd/manager/`): Entry point and operator initialization

### Custom Resource Definitions
- **Application** (app.k8s.io/v1beta1): Primary resource managed by this operator
- **Deployable** (internal API): Represents deployable resources
- **Subscription** (from multicloud-operators-subscription): External dependency for subscription management

### Key Relationships
- Applications use selectors to discover and aggregate Subscriptions and Deployables
- The operator watches for Application CR changes and updates status based on discovered resources
- Integration with Open Cluster Management's subscription and channel operators

## Development Commands

### Building and Testing
```bash
# Build the operator binary
make build

# Build local binary for macOS development
make local

# Run tests
make test
go test ./...

# Build container image
make build-images
```

### Code Quality
```bash
# Run all linting (includes go, yaml, markdown, etc.)
make lint

# Format code
make fmt

# Run specific linters
make lint-go
make lint-yaml
```

### Development Setup
```bash
# Clone and set up development environment
export GITHUB_USER=<github_user>
export GITHUB_TOKEN=<github_token>
make

# Apply CRDs for standalone development
kubectl apply -f deploy/crds/standalone

# Run operator locally
export POD_NAMESPACE=<namespace>
./build/_output/bin/multicluster-operators-application --application-crd-file deploy/crds/app.k8s.io_applications_crd_v1.yaml
```

### Deployment
```bash
# Deploy to cluster
kubectl apply -f deploy/crds/standalone
kubectl apply -f deploy/crds
kubectl apply -f deploy
```

## Important Files and Directories

- `cmd/manager/main.go`: Application entry point
- `pkg/controller/application/application_controller.go`: Core controller logic
- `utils/application.go`: Application-specific utility functions
- `webhook/application_validator.go`: Admission webhook implementation
- `deploy/`: Kubernetes manifests for deployment
- `Makefile`: Build system with common development tasks
- `go.mod`: Go module dependencies (Go 1.23+)

## Testing Framework

The project uses:
- **Ginkgo v2** for BDD-style testing
- **Gomega** for assertions
- Controller-runtime test framework for Kubernetes controller testing
- Test files follow `*_test.go` naming convention
- Test suites use `*_suite_test.go` pattern

## Build System

The project uses a custom build harness system:
- Main Makefile includes `common/Makefile.common.mk`
- Build scripts in `common/scripts/`
- Supports Travis CI integration
- Container builds use `quay.io/stolostron` registry by default