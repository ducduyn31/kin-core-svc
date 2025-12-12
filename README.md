# Kin - Stay Connected

Kin is an instant messaging app designed to keep you connected with the people who matter most, no matter the distance.

## The Problem

When you're separated from loved ones by timezones and busy schedules, finding the right moment to connect is hard. You don't want to call when they're sleeping, working, or commuting.

## The Solution

Kin helps you know when your people are available. Share as much or as little as you're comfortable with - from a simple "free to talk" status to real-time location updates.

### Features

- **Instant Messaging**: Text, photos, videos, voice notes, files
- **Kin Circles**: Groups designed for your closest relationships
- **Smart Availability**: Know when your loved ones are free to talk
- **Timezone Aware**: Never accidentally call at 3am again
- **Privacy First**: You control exactly what you share with each circle

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- [buf](https://buf.build/) (for proto generation)
- Auth0 account

### Quick Start

1. Clone the repository
```bash
git clone https://github.com/danielng/kin-core-svc.git
cd kin-core-svc
```

2. Copy `.env.example` to `.env` and configure
```bash
cp .env.example .env
# Edit .env with your Auth0 credentials
```

3. Run infrastructure services
```bash
task docker:up
```

4. Generate code (protobuf, gRPC, Bob models)
```bash
task generate
```

5. Run database migrations
```bash
task migrate:up
```

6. Run the application
```bash
task run
```

### Configuration

See `.env.example` for required environment variables.

| Variable | Description |
|----------|-------------|
| `DATABASE_WRITE_URL` | PostgreSQL write connection string |
| `DATABASE_READ_URL` | PostgreSQL read connection string |
| `REDIS_URL` | Redis connection string |
| `AUTH0_DOMAIN` | Auth0 tenant domain |
| `AUTH0_AUDIENCE` | Auth0 API audience |
| `S3_ENDPOINT` | S3-compatible storage endpoint |
| `S3_ACCESS_KEY` | S3 access key |
| `S3_SECRET_KEY` | S3 secret key |
| `S3_BUCKET` | S3 bucket name |
| `GRPC_PORT` | gRPC server port (default: 50051) |
| `GRPC_GATEWAY_PORT` | REST gateway port (default: 8080) |

## Development

### Available Commands

```bash
task build           # Build the application
task run             # Run the application
task test            # Run tests
task lint            # Run linter
task generate        # Generate all code (proto + Bob)
task proto:generate  # Generate protobuf/gRPC code
task bob:generate    # Generate Bob ORM models
task docker:up       # Start Docker containers
task docker:down     # Stop Docker containers
task migrate:up      # Run all migrations
```

### Project Structure

```
kin/
├── cmd/api/              # Application entry point
├── proto/                # Protocol buffer definitions
│   └── kin/v1/           # API version 1 protos
├── gen/                  # Generated protobuf/gRPC code
│   ├── proto/            # Go protobuf and gRPC stubs
│   └── openapi/          # OpenAPI specs
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Business logic (DDD aggregates)
│   ├── application/      # Use cases and services
│   ├── infrastructure/   # External implementations
│   └── interfaces/grpc/  # gRPC handlers and gateway
├── pkg/                  # Shared utilities
└── config/               # Configuration files
```

## API

The service exposes both gRPC and REST APIs:

| Service | Port | Description |
|---------|------|-------------|
| gRPC | 50051 | Native gRPC API |
| REST | 8080 | gRPC-Gateway (HTTP/JSON) |

### Health Endpoints

```bash
# Health check (checks DB, Redis)
curl http://localhost:8080/health

# Readiness check (version info)
curl http://localhost:8080/ready
```

### gRPC (with grpcurl)

```bash
# List services (requires reflection enabled)
grpcurl -plaintext localhost:50051 list

# Get current user
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  localhost:50051 kin.v1.UserService/GetMe

# Create circle
grpcurl -plaintext \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Family"}' \
  localhost:50051 kin.v1.CircleService/CreateCircle
```

### REST (via gRPC-Gateway)

```bash
# Get current user
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/me

# Create circle
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Family"}' \
  http://localhost:8080/api/v1/circles
```

### OpenAPI

Generated OpenAPI specs are available in `gen/openapi/` after running `task generate`.

## License

[License details]
