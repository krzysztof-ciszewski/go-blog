# Blog Application

A microservices-based blog application built with Go, implementing CQRS (Command Query Responsibility Segregation) pattern with Kafka for asynchronous message processing.

## Architecture

This project follows **Domain-Driven Design (DDD)** principles and implements the **CQRS pattern**:

- **Commands**: Write operations that modify state (e.g., creating posts)
- **Queries**: Read operations that retrieve data (e.g., fetching posts)
- **Events**: Domain events published after state changes

### Project Structure

```
blog/
├── cmd/                          # Application entry points
│   ├── server.go                  # HTTP API server
│   ├── consume.go                 # Kafka consumer service
│   └── migrate.go                 # Database migration runner
├── internal/
│   ├── Application/              # Application layer (CQRS)
│   │   ├── Command/              # Command handlers
│   │   ├── Query/                # Query handlers
│   │   └── View/                 # Read models
│   ├── Domain/                   # Domain layer
│   │   ├── Entity/              # Domain entities
│   │   └── Repository/           # Repository interfaces
│   ├── Infrastructure/           # Infrastructure layer
│   │   ├── DependencyInjection/  # DI container
│   │   └── Repository/           # Repository implementations
│   └── UserInterface/           # Presentation layer
│       └── Api/                 # REST API handlers
├── db/
│   └── migrations/              # Database migrations
└── docker-compose.yaml          # Docker Compose configuration
```

## Features

- **CQRS Pattern**: Separates read and write operations
- **Event-Driven Architecture**: Uses Kafka for asynchronous message processing
- **Domain-Driven Design**: Clean architecture with clear separation of concerns
- **RESTful API**: HTTP endpoints for blog operations
- **PostgreSQL**: Persistent data storage
- **Database Migrations**: Version-controlled schema changes

## Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose
- Make (optional, for convenience commands)

## Getting Started

### Using Docker Compose (Recommended)

1. **Start all services:**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - PostgreSQL database
   - Zookeeper
   - Kafka
   - Database migration service
   - API server (port 8080)
   - Kafka consumer service

2. **Check service status:**
   ```bash
   docker-compose ps
   ```

3. **View logs:**
   ```bash
   docker-compose logs -f server
   docker-compose logs -f consume
   ```

### Manual Setup

1. **Start dependencies:**
   ```bash
   docker-compose up -d postgres zookeeper kafka
   ```

2. **Set environment variables:**
   ```bash
   export DATABASE_URL="postgres://blog:blogpassword@localhost:5432/blog?sslmode=disable"
   export KAFKA_BROKER="localhost:9092"
   ```

3. **Run database migrations:**
   ```bash
   go run cmd/migrate.go
   ```

4. **Start the API server:**
   ```bash
   go run cmd/server.go
   ```

5. **Start the consumer service (in a separate terminal):**
   ```bash
   go run cmd/consume.go
   ```

## API Endpoints

### Health Check

```bash
GET /ping
```

**Response:**
```json
{
  "message": "pong"
}
```

### Create Post

```bash
POST /posts
Content-Type: application/json

{
  "Id": "550e8400-e29b-41d4-a716-446655440000",
  "Slug": "my-first-post",
  "Title": "My First Post",
  "Content": "This is the content of my first blog post.",
  "Author": "John Doe"
}
```

**Response:**
```json
{
  "message": "Post created"
}
```

**Validation Rules:**
- `Id`: Required, must be a valid UUID
- `Slug`: Required, 3-255 characters, alphanumeric
- `Title`: Required, 3-255 characters
- `Content`: Required, 10-10000 characters
- `Author`: Required, 3-255 characters

### Get Post

```bash
GET /posts/:id
```

## How It Works

### Command Flow

1. **Client sends POST request** to `/posts` endpoint
2. **Server validates** the request and creates a `CreatePostCommand`
3. **Command is published** to Kafka topic `commands.CreatePostCommand`
4. **Consumer service** receives the command from Kafka
5. **Command handler** processes the command and saves the post to the database
6. **Events can be published** for further processing (e.g., notifications, search indexing)

### Message Topics

- **Commands**: `commands.{CommandName}` (e.g., `commands.CreatePostCommand`)
- **Events**: `events.{EventName}` (e.g., `events.PostCreated`)

## Database Schema

### Posts Table

```sql
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author TEXT NOT NULL
);
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build server
go build -o bin/server ./cmd/server.go

# Build consumer
go build -o bin/consume ./cmd/consume.go

# Build migration tool
go build -o bin/migrate ./cmd/migrate.go
```

### Database Migrations

Migrations are located in `db/migrations/`. The migration service runs automatically when using Docker Compose.

To create a new migration:

1. Create `{version}_migration_name.up.sql` for the migration
2. Create `{version}_migration_name.down.sql` for the rollback

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://blog:blogpassword@postgres:5432/blog?sslmode=disable` |
| `KAFKA_BROKER` | Kafka broker address | `kafka:9093` (Docker) or `localhost:9092` (local) |

## Dependencies

### Core Libraries

- **Gin**: HTTP web framework
- **Watermill**: Event-driven architecture library
- **Watermill-Kafka**: Kafka integration for Watermill
- **golang-migrate**: Database migration tool
- **PostgreSQL Driver**: Database connectivity

### Architecture Libraries

- **CQRS**: Command Query Responsibility Segregation implementation
- **UUID**: Unique identifier generation

## Docker Services

The `docker-compose.yaml` file defines the following services:

- **postgres**: PostgreSQL 16 database
- **zookeeper**: Zookeeper for Kafka coordination
- **kafka**: Apache Kafka message broker
- **migrate**: Database migration service (runs once)
- **server**: HTTP API server
- **consume**: Kafka consumer service

## Stopping Services

```bash
docker-compose down
```

To remove volumes (database data):

```bash
docker-compose down -v
```

## Troubleshooting

### Kafka Connection Issues

If the consumer can't connect to Kafka:

1. Verify Kafka is running: `docker-compose ps kafka`
2. Check Kafka logs: `docker-compose logs kafka`
3. Ensure `KAFKA_BROKER` environment variable is correct

### Database Connection Issues

1. Verify PostgreSQL is running: `docker-compose ps postgres`
2. Check database logs: `docker-compose logs postgres`
3. Verify `DATABASE_URL` environment variable

### Migration Issues

1. Check migration logs: `docker-compose logs migrate`
2. Ensure migrations are in `db/migrations/` directory
3. Verify migration file naming convention

## License

This project is provided as-is for educational and development purposes.

