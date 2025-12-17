# Blog Application

A microservices-based blog application built with Go, implementing CQRS (Command Query Responsibility Segregation) pattern with Kafka for asynchronous message processing.

## Architecture

This project follows **Domain-Driven Design (DDD)** principles and implements the **CQRS pattern**:

- **Commands**: Write operations that modify state (e.g., creating, deleting posts). Commands are sent asynchronously via Kafka.
- **Queries**: Read operations that retrieve data (e.g., fetching posts, filtering by author/text). Queries are handled synchronously via a Query Bus.
- **Events**: Domain events published after state changes for further processing (e.g., notifications, search indexing)

### Query Bus

The application implements a **Query Bus pattern** for handling read operations synchronously. Unlike commands which are processed asynchronously via Kafka, queries are executed directly through the Query Bus, providing immediate responses to API requests.

The Query Bus supports:
- **GetPostQuery**: Retrieve a single post by UUID
- **FindBySlugQuery**: Retrieve a post by its unique slug
- **FindAllQuery**: Retrieve all posts
- **FindAllByTextQuery**: Search posts by content text
- **FindAllByAuthorQuery**: Filter posts by author name

### Project Structure

```
blog/
├── cmd/                          # Application entry points
│   ├── server.go                  # HTTP API server
│   ├── consume.go                 # Kafka consumer service
│   └── migrate.go                 # Database migration runner
├── internal/
│   ├── Application/              # Application layer (CQRS)
│   │   ├── Command/              # Command handlers (CreatePost, DeletePost)
│   │   ├── Query/                # Query handlers (GetPost, FindAll, FindBySlug, etc.)
│   │   └── View/                 # Read models
│   ├── Domain/                   # Domain layer
│   │   ├── Entity/              # Domain entities
│   │   └── Repository/           # Repository interfaces
│   ├── Infrastructure/           # Infrastructure layer
│   │   ├── DependencyInjection/  # DI container
│   │   ├── QueryBus/            # Query Bus implementation
│   │   └── Repository/           # Repository implementations
│   └── UserInterface/           # Presentation layer
│       └── Api/                 # REST API handlers
├── db/
│   └── migrations/              # Database migrations
└── docker-compose.yaml          # Docker Compose configuration
```

## Features

- **CQRS Pattern**: Separates read and write operations
  - **Commands**: Asynchronous write operations via Kafka
  - **Queries**: Synchronous read operations via Query Bus
- **Query Bus**: Synchronous query handling for immediate API responses
- **Event-Driven Architecture**: Uses Kafka for asynchronous message processing
- **Domain-Driven Design**: Clean architecture with clear separation of concerns
- **RESTful API**: HTTP endpoints for blog operations with filtering and search
- **PostgreSQL**: Persistent data storage with proper data types
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

### Get Post by ID

```bash
GET /posts/:id
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-12-17T19:16:32Z",
  "updated_at": "2025-12-17T19:16:32Z",
  "slug": "my-first-post",
  "title": "My First Post",
  "content": "This is the content of my first blog post.",
  "author": "John Doe"
}
```

### Get Post by Slug

```bash
GET /posts/slug/:slug
```

**Example:**
```bash
GET /posts/slug/my-first-post
```

**Response:** Same as Get Post by ID

### List Posts

```bash
GET /posts
GET /posts?text=search+term
GET /posts?author=John+Doe
```

**Query Parameters:**
- `text` (optional): Search posts by content text
- `author` (optional): Filter posts by author name

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-12-17T19:16:32Z",
    "updated_at": "2025-12-17T19:16:32Z",
    "slug": "my-first-post",
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "author": "John Doe"
  }
]
```

### Delete Post

```bash
DELETE /posts/:id
```

**Response:**
```json
{
  "message": "Post deleted"
}
```

**Note:** The delete operation is processed asynchronously via Kafka.

## How It Works

### Command Flow (Write Operations)

1. **Client sends POST/DELETE request** to `/posts` or `/posts/:id` endpoint
2. **Server validates** the request and creates a command (`CreatePostCommand` or `DeletePostCommand`)
3. **Command is published** to Kafka topic `commands.{CommandName}`
4. **Consumer service** receives the command from Kafka
5. **Command handler** processes the command and modifies the database
6. **Events can be published** for further processing (e.g., notifications, search indexing)

### Query Flow (Read Operations)

1. **Client sends GET request** to a query endpoint (e.g., `/posts`, `/posts/:id`, `/posts/slug/:slug`)
2. **Server creates a query object** (e.g., `GetPostQuery`, `FindAllQuery`, `FindBySlugQuery`)
3. **Query is executed** synchronously through the Query Bus
4. **Query handler** retrieves data from the repository
5. **Response is returned** immediately to the client

**Key Difference:** Queries are handled synchronously for immediate responses, while commands are processed asynchronously via Kafka for better scalability and decoupling.

### Message Topics

- **Commands**: `commands.{CommandName}` (e.g., `commands.CreatePostCommand`)
- **Events**: `events.{EventName}` (e.g., `events.PostCreated`)

## Database Schema

### Posts Table

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL
);
```

**Schema Details:**
- **id**: UUID type for proper unique identifier handling
- **created_at/updated_at**: TIMESTAMPTZ (timestamp with timezone) for accurate datetime tracking
- **slug**: VARCHAR(255) with UNIQUE constraint for SEO-friendly URLs
- **title**: VARCHAR(255) for post titles
- **content**: TEXT type for unlimited post content
- **author**: VARCHAR(255) for author names

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

- **CQRS**: Command Query Responsibility Segregation implementation (Watermill)
- **Query Bus**: Custom synchronous query handling implementation
- **UUID**: Unique identifier generation (Google UUID library)

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

