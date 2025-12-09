# Kin - Development Guide

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.23+ |
| HTTP Framework | Gin |
| Configuration | Viper (YAML + environment variables) |
| Database | PostgreSQL with pgx/v5 |
| SQL Builder | Bob (github.com/stephenafamo/bob) |
| Migrations | pgroll (zero-downtime migrations) |
| Authentication | Auth0 (JWT validation) |
| Cache/Presence | Redis (go-redis/v9) |
| Object Storage | S3-compatible |
| Logging | slog (stdlib) |
| Validation | go-playground/validator/v10 |

## Architecture

This project follows Domain-Driven Design (DDD) with clean architecture principles.

### Layer Overview
```
interfaces → application → domain ← infrastructure
```

- **Domain**: Core business logic, no external dependencies
- **Application**: Use cases, orchestrates domain objects
- **Infrastructure**: External implementations (database, cache, auth)
- **Interfaces**: Entry points (HTTP handlers)

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

### Bob Code Generation

After migrations, regenerate Bob models:
```bash
task bob:generate
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

### Auth Middleware Flow
```
Request → Extract JWT → Validate with Auth0 JWKS → Get/Create User → Set in Context
```

### Getting Current User
```go
func (h *Handler) SomeEndpoint(c *gin.Context) {
    user := c.MustGet(ctxkey.User).(*domain.User)
    userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
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

# Generate Bob models
task bob:generate

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
- Use `pkg/apperror` for HTTP error responses

### Naming

- Use singular for domain packages: `user`, `circle`, `message`
- Use `_repo.go` suffix for repository implementations
- Use `_handler.go` suffix for HTTP handlers

### Testing

- Unit tests alongside source files: `user_test.go`
- Integration tests in `_test` package
- Use testcontainers for database tests

## API Conventions

### Response Format
```json
{
  "data": {},
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100
  }
}
```

### Error Format
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": []
  }
}
```

### Pagination

Use cursor-based pagination for lists:
```
GET /circles/{id}/messages?cursor=xxx&limit=50
```
