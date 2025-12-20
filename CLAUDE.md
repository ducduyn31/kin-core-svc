# Kin - Development Guide

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.23+ |
| API | Connect RPC (HTTP/1.1, HTTP/2, gRPC compatible) |
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
- **Interfaces**: Entry points (Connect RPC handlers)

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

Generate all code (protobuf, Connect, Bob models):
```bash
task generate
```

Or run individually:
```bash
task proto:generate  # Generate protobuf and Connect code
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

### Auth Flow (Connect RPC)
```
Request → Extract JWT from Authorization header → Validate with Auth0 JWKS → Get/Create User → Set in Context
```

### Getting Current User (Connect handlers)
```go
func (h *Handler) SomeMethod(ctx context.Context, req *connect.Request[pb.SomeRequest]) (*connect.Response[pb.SomeResponse], error) {
    user := ctx.Value(interceptors.UserKey).(*domain.User)
    userID := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
}
```

## Configuration

Viper loads config from:
1. `config/config.yaml` (base, committed)
2. `config/config.{env}.yaml` (environment override, gitignored for local dev)
3. Environment variables (highest priority)

### Local Development Config

```bash
# Initialize local config from template
task config:init
# Or manually:
cp config/config.development.yaml.example config/config.development.yaml
```

The `config.development.yaml` file is gitignored - each developer can customize without affecting others.

Environment variables use underscore notation: `DATABASE_URL`, `AUTH0_DOMAIN`

## Development Commands
```bash
# Run locally
task run

# Run tests
task test

# Run linter
task lint

# Generate all code (proto, Connect, Bob)
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
- Use `pkg/apperror` for application errors (mapped to Connect error codes)

### Naming

- Use singular for domain packages: `user`, `circle`, `message`
- Use `_repo.go` suffix for repository implementations
- Use `_handler.go` suffix for Connect handlers

### Testing

- Unit tests alongside source files: `user_test.go`
- Integration tests in `_test` package
- Use testcontainers for database tests

## API

The service uses Connect RPC which provides a single HTTP server supporting multiple protocols:
- **Connect protocol**: Simple HTTP/JSON or HTTP/Protobuf
- **gRPC protocol**: Full gRPC compatibility (HTTP/2)
- **gRPC-Web protocol**: Browser-compatible gRPC

All protocols are served on a single port (default: 8080).

### Port

| Service | Port | Purpose |
|---------|------|---------|
| Connect | 8080 | HTTP server (Connect, gRPC, gRPC-Web) |

### Health Endpoints

```bash
# Health check (checks DB, Redis)
curl http://localhost:8080/health

# Readiness check (version info)
curl http://localhost:8080/ready
```

### Configuration

```yaml
server:
  port: 8080
```

Environment variables:
- `PORT` - Server port (default: 8080)

### Services

- `UserService` - User profile and preferences
- `CircleService` - Circle management, members, invitations

### Testing with curl (Connect Protocol)

Connect RPC supports simple JSON over HTTP:

```bash
# Get current user (POST with JSON)
curl -X POST http://localhost:8080/kin.v1.UserService/GetMe \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{}'

# Create circle
curl -X POST http://localhost:8080/kin.v1.CircleService/CreateCircle \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Family", "description": "My family circle"}'

# List circles
curl -X POST http://localhost:8080/kin.v1.CircleService/ListCircles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"limit": 20}'
```

### Testing with grpcurl (gRPC Protocol)

Connect also supports native gRPC protocol:

```bash
# Call GetMe (with auth token)
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  -d '{}' \
  localhost:8080 kin.v1.UserService/GetMe

# Create circle
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Family", "description": "My family circle"}' \
  localhost:8080 kin.v1.CircleService/CreateCircle
```

### Proto Files

Location: `proto/kin/v1/`

Generate code:
```bash
task proto:generate
```

Generated files:
- `gen/proto/kin/v1/*.pb.go` - Protobuf messages
- `gen/proto/kin/v1/kinv1connect/*.connect.go` - Connect service handlers and clients

## Bruno API Collection

The Bruno collection for testing Connect endpoints is located in `kin-core-api/`.

### Structure
```text
kin-core-api/
├── bruno.json              # Collection config
├── environments/
│   └── local.bru           # Environment variables (base_url, auth_token)
├── connect/                # Connect protocol endpoints (HTTP/JSON)
│   ├── UserService/
│   │   ├── GetMe.bru
│   │   ├── UpdateProfile.bru
│   │   ├── UpdateTimezone.bru
│   │   ├── GetPreferences.bru
│   │   └── UpdatePreferences.bru
│   └── CircleService/
│       ├── CreateCircle.bru
│       ├── ListCircles.bru
│       ├── GetCircle.bru
│       ├── UpdateCircle.bru
│       ├── DeleteCircle.bru
│       ├── LeaveCircle.bru
│       ├── ListMembers.bru
│       ├── AddMember.bru
│       ├── RemoveMember.bru
│       ├── GetSharingPreference.bru
│       ├── UpdateSharingPreference.bru
│       ├── CreateInvitation.bru
│       └── AcceptInvitation.bru
└── grpc/                   # gRPC protocol endpoints
    ├── UserService/
    │   ├── GetMe.bru
    │   ├── UpdateProfile.bru
    │   ├── UpdateTimezone.bru
    │   ├── GetPreferences.bru
    │   └── UpdatePreferences.bru
    └── CircleService/
        ├── CreateCircle.bru
        ├── ListCircles.bru
        ├── GetCircle.bru
        ├── UpdateCircle.bru
        ├── DeleteCircle.bru
        ├── LeaveCircle.bru
        ├── ListMembers.bru
        ├── AddMember.bru
        ├── RemoveMember.bru
        ├── GetSharingPreference.bru
        ├── UpdateSharingPreference.bru
        ├── CreateInvitation.bru
        └── AcceptInvitation.bru
```

### Environment Variables

Set `auth_token` in Bruno's environment settings (stored as secret, not committed).

| Variable | Description | Default |
|----------|-------------|---------|
| `base_url` | Connect API base URL | `http://localhost:8080` |
| `auth_token` | JWT bearer token (secret) | - |

### Maintaining the Collection

**IMPORTANT**: When modifying proto files in `proto/kin/v1/`, update the Bruno collection accordingly:

1. **Adding new RPC endpoints**: Create `.bru` files in both `connect/{Service}/` and `grpc/{Service}/`
2. **Removing RPC endpoints**: Delete the corresponding `.bru` files from both folders
3. **Modifying request/response fields**: Update `body:json` (Connect) and `body:grpc` (gRPC) sections
4. **Adding new services**: Create new folders under both `connect/` and `grpc/`

### Sequence Numbering Convention

The `meta.seq` field controls sort order in Bruno's UI. **Seq must be unique within each service folder**:
- `connect/UserService/` - seq 1-5
- `connect/CircleService/` - seq 1-13
- `grpc/UserService/` - seq 1-5
- `grpc/CircleService/` - seq 1-13

When adding new endpoints, use the next available seq number within the target service folder.

### Bruno File Templates

**Connect Endpoint (HTTP/JSON):**
```bru
meta {
  name: EndpointName
  type: http
  seq: 1
}

post {
  url: {{base_url}}/kin.v1.ServiceName/MethodName
  body: json
  auth: bearer
}

auth:bearer {
  token: {{auth_token}}
}

body:json {
  {
    "field": "value"
  }
}
```

**gRPC Endpoint:**
```bru
meta {
  name: EndpointName
  type: grpc
  seq: 1
}

grpc {
  url: {{base_url}}
  method: /kin.v1.ServiceName/MethodName
  body: grpc
  auth: bearer
  methodType: unary
}

auth:bearer {
  token: {{auth_token}}
}

body:grpc {
  name: message 1
  content: '''
    {
      "field": "value"
    }
  '''
}
```
