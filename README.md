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

4. Run database migrations
```bash
task migrate:up
```

5. Run the application
```bash
task run
```

### Configuration

See `.env.example` for required environment variables.

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `REDIS_URL` | Redis connection string |
| `AUTH0_DOMAIN` | Auth0 tenant domain |
| `AUTH0_AUDIENCE` | Auth0 API audience |
| `S3_ENDPOINT` | S3-compatible storage endpoint |
| `S3_ACCESS_KEY` | S3 access key |
| `S3_SECRET_KEY` | S3 secret key |
| `S3_BUCKET` | S3 bucket name |

## Development

### Available Commands

```bash
task build           # Build the application
task run             # Run the application
task test            # Run tests
task lint            # Run linter
task docker:up       # Start Docker containers
task docker:down     # Stop Docker containers
task migrate:up      # Run all migrations
task bob:generate    # Generate Bob ORM models
```

### Project Structure

```
kin/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Business logic (DDD aggregates)
│   ├── application/      # Use cases and services
│   ├── infrastructure/   # External implementations
│   └── interfaces/       # HTTP handlers
├── pkg/                  # Shared utilities
└── config/               # Configuration files
```

## API Documentation

API documentation is available at `/swagger` when running in development mode.

## License

[License details]
