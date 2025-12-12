# Kin - Development Guide

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.23+ |
| API | gRPC + gRPC-Gateway (REST) |
| Configuration | Viper (YAML + environment variables) |
| Database | PostgreSQL with pgx/v5 |
| SQL Builder | Bob (github.com/stephenafamo/bob) |
| Migrations | pgroll (zero-downtime migrations) |
| Authentication | Auth0 (JWT validation) |
| Cache/Presence | Redis (go-redis/v9) |
| Object Storage | S3-compatible |
| Logging | slog (stdlib) |
| Proto Generation | buf (buf.build) |

## Architecture

This project follows Domain-Driven Design (DDD) with clean architecture principles.

### Layer Overview
```
interfaces → application → domain ← infrastructure
```

- **Domain**: Core business logic, no external dependencies
- **Application**: Use cases, orchestrates domain objects
- **Infrastructure**: External implementations (database, cache, auth)
- **Interfaces**: Entry points (gRPC handlers)

### Dependency Rule

Dependencies point inward. Domain layer has zero external dependencies.

### Domain Structure

Each domain contains:
- `{entity}.go` - Aggregate root and entities
- `repository.go` - Repository interface
- `errors.go` - Domain-specific errors
- `events.go` - Domain events (if applicable)

### Key Domains

| Domain | Responsibility |
|--------|----------------|
| user | Profile, Auth0 link, preferences |
| contact | Contact list management |
| circle | Kin circles, members, sharing preferences |
| conversation | Chat threads, participants |
| messaging | Messages, receipts, reactions |
| availability | Availability windows, status |
| location | Location sharing, places, check-ins |
| presence | Online status, activity |
| media | File uploads |
| notification | Push notifications |

## Database

### Migrations (pgroll)

Migrations are JSON files in `internal/infrastructure/postgres/migrations/`.
```bash
# Start a migration (creates versioned schema)
task migrate:start FILE=0001_create_users.json

# Complete migration (drops old schema)
task migrate:complete

# Rollback
task migrate:rollback

# Check status
task migrate:status
```

### Code Generation

Generate all code (protobuf, gRPC, Bob models):
```bash
task generate
```

Or run individually:
```bash
task proto:generate  # Generate protobuf and gRPC code
task bob:generate    # Generate Bob database models
```

### Bob Usage
```go
// Select
query := psql.Select(
    sm.Columns("id", "display_name"),
    sm.From("users"),
    sm.Where(psql.Quote("id").EQ(psql.Arg(userID))),
)

// Insert
query := psql.Insert(
    im.Into("users", "id", "auth0_sub", "display_name"),
    im.Values(psql.Arg(id, auth0Sub, name)),
    im.OnConflict("auth0_sub").DoUpdate(
        im.SetExcluded("display_name"),
    ),
)

// Update
query := psql.Update(
    um.Table("users"),
    um.Set("display_name").ToArg(name),
    um.Where(psql.Quote("id").EQ(psql.Arg(id))),
)
```

## Authentication

Auth0 handles authentication. The backend:
1. Validates JWT tokens from Auth0
2. Extracts `sub` claim to identify user
3. Creates/retrieves local user record linked by `auth0_sub`

### Auth Flow (gRPC)
```
Request → Extract JWT from metadata → Validate with Auth0 JWKS → Get/Create User → Set in Context
```

### Getting Current User (gRPC handlers)
```go
func (h *Handler) SomeMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    user := ctx.Value(interceptors.UserKey).(*domain.User)
    userID := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
}
```

## Configuration

Viper loads config from:
1. `config/config.yaml` (base)
2. `config/config.{env}.yaml` (environment override)
3. Environment variables (highest priority)

Environment variables use underscore notation: `DATABASE_URL`, `AUTH0_DOMAIN`

## Development Commands
```bash
# Run locally
task run

# Run tests
task test

# Run linter
task lint

# Generate all code (proto, gRPC, Bob)
task generate

# Start all migrations
task migrate:up

# Docker development environment
task docker:up
task docker:down

# Build binary
task build
```

## Project Conventions

### Error Handling

- Define domain errors in each domain's `errors.go`
- Wrap infrastructure errors with context
- Use `pkg/apperror` for application errors (mapped to gRPC status codes)

### Naming

- Use singular for domain packages: `user`, `circle`, `message`
- Use `_repo.go` suffix for repository implementations
- Use `_handler.go` suffix for gRPC handlers

### Testing

- Unit tests alongside source files: `user_test.go`
- Integration tests in `_test` package
- Use testcontainers for database tests

## API

The service exposes both gRPC and REST APIs:
- **gRPC**: Native protocol buffer communication on port 50051
- **REST**: gRPC-Gateway reverse proxy on port 8080

### Ports

| Service | Port | Purpose |
|---------|------|---------|
| gRPC | 50051 | Native gRPC API |
| REST | 8080 | gRPC-Gateway (HTTP/JSON) |

### Health Endpoints (REST)

```bash
# Health check (checks DB, Redis)
curl http://localhost:8080/health

# Readiness check (version info)
curl http://localhost:8080/ready
```

### Configuration

```yaml
grpc:
  port: 50051
  enable_reflection: true    # For grpcurl debugging
  gateway_port: 8080
```

Environment variables:
- `GRPC_PORT` - gRPC server port (default: 50051)
- `GRPC_ENABLE_REFLECTION` - Enable reflection service for debugging
- `GRPC_GATEWAY_PORT` - REST gateway port (default: 8080)

### Services

- `UserService` - User profile and preferences
- `CircleService` - Circle management, members, invitations

### Testing with grpcurl

```bash
# List services (requires reflection)
grpcurl -plaintext localhost:50051 list

# Call GetMe (with auth token)
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  localhost:50051 kin.v1.UserService/GetMe

# Create circle
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Family", "description": "My family circle"}' \
  localhost:50051 kin.v1.CircleService/CreateCircle
```

### REST API (via gRPC-Gateway)

The gateway exposes REST endpoints that proxy to gRPC:

```bash
# Get current user
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/me

# Create circle
curl -X POST -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Family", "description": "My family circle"}' \
  http://localhost:8080/api/v1/circles
```

### Proto Files

Location: `proto/kin/v1/`

Generate code:
```bash
task proto:generate
```

Generated files:
- `gen/proto/kin/v1/*.pb.go` - Protobuf messages
- `gen/proto/kin/v1/*_grpc.pb.go` - gRPC service interfaces
- `gen/proto/kin/v1/*.pb.gw.go` - gRPC-Gateway handlers
- `gen/openapi/kin/v1/*.swagger.json` - OpenAPI specs
