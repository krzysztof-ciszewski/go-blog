# Blog Application

A microservices-based blog application built with Go, implementing CQRS (Command Query Responsibility Segregation) pattern with RabbitMQ for asynchronous message processing.

## Architecture

This project follows **Domain-Driven Design (DDD)** principles and implements the **CQRS pattern**:

- **Commands**: Write operations that modify state (e.g., creating posts, deleting posts, creating users). Commands are sent asynchronously via RabbitMQ.
- **Queries**: Read operations that retrieve data (e.g., fetching posts, filtering by author/text). Queries are handled synchronously via a Query Bus.
- **Events**: Domain events published after state changes for further processing (e.g., notifications, search indexing)
- **Dead Letter Queue**: Failed messages are automatically routed to dead letter queues via a custom topology builder that configures RabbitMQ dead letter exchanges

### Query Bus

The application implements a **Query Bus pattern** for handling read operations synchronously. Unlike commands which are processed asynchronously via RabbitMQ, queries are executed directly through the Query Bus, providing immediate responses to API requests.

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
│   ├── consume.go                 # RabbitMQ consumer service
│   └── migrate.go                 # Database migration runner
├── internal/
│   ├── Application/              # Application layer (CQRS)
│   │   ├── Command/              # Command handlers
│   │   │   ├── Post/            # Post commands (CreatePost, DeletePost)
│   │   │   └── User/            # User commands (CreateUser)
│   │   ├── Query/                # Query handlers (GetPost, FindAll, FindBySlug, etc.)
│   │   └── View/                 # Read models
│   ├── Domain/                   # Domain layer
│   │   ├── Entity/              # Domain entities (Post, User)
│   │   └── Repository/           # Repository interfaces (PostRepository, UserRepository)
│   ├── Infrastructure/           # Infrastructure layer
│   │   ├── Amqp/                # AMQP topology builder for dead letter queues
│   │   ├── DependencyInjection/  # DI container
│   │   ├── QueryBus/            # Query Bus implementation
│   │   └── Repository/           # Repository implementations
│   └── UserInterface/           # Presentation layer
│       └── Api/                 # REST API handlers and middleware
├── db/
│   └── migrations/              # Database migrations
└── docker-compose.yaml          # Docker Compose configuration
```

## Features

- **CQRS Pattern**: Separates read and write operations
  - **Commands**: Asynchronous write operations via RabbitMQ
  - **Queries**: Synchronous read operations via Query Bus
- **Query Bus**: Synchronous query handling for immediate API responses
- **Event-Driven Architecture**: Uses RabbitMQ (AMQP) for asynchronous message processing
- **Dead Letter Queue**: Automatic routing of failed messages to dead letter queues via custom RabbitMQ topology builder
- **Domain-Driven Design**: Clean architecture with clear separation of concerns
- **RESTful API**: HTTP endpoints for blog operations with filtering and search
- **OAuth Authentication**: GitHub OAuth integration for user authentication
- **User Management**: User entity with OAuth token storage
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
   - RabbitMQ message broker
   - Database migration service
   - API server (port 8080)
   - RabbitMQ consumer service

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
   docker-compose up -d postgres rabbitmq
   ```

2. **Set environment variables:**
   ```bash
   export DATABASE_URL="postgres://blog:blogpassword@localhost:5432/blog?sslmode=disable"
   export AMQP_URL="amqp://guest:guest@localhost:5672/"
   export AMQP_DLX_EXCHANGE="my-dlx"
   export AMQP_DLX_QUEUE_SUFFIX="dlq"
   export AMQP_DLX_ROUTING_KEY_SUFFIX="dlq"
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

**Note:** The delete operation is processed asynchronously via RabbitMQ.

### Authentication Endpoints

#### OAuth Login

```bash
GET /auth/:provider
```

**Example:**
```bash
GET /auth/github
```

**Response:** Redirects to OAuth provider for authentication

#### OAuth Callback

```bash
GET /auth/:provider/callback
```

**Example:**
```bash
GET /auth/github/callback
```

**Response:** Redirects to client URL after successful authentication

**Note:** This endpoint automatically creates a user account if one doesn't exist. User creation is processed asynchronously via RabbitMQ using `CreateUserCommand`.

#### Logout

```bash
GET /auth/logout/:provider
```

**Example:**
```bash
GET /auth/logout/github
```

**Response:** Logs out the user and redirects to home

## How It Works

### Command Flow (Write Operations)

1. **Client sends POST/DELETE request** to `/posts` or `/posts/:id` endpoint, OR
   **OAuth callback** triggers user creation via `/auth/:provider/callback`
2. **Server validates** the request and creates a command (`CreatePostCommand`, `DeletePostCommand`, or `CreateUserCommand`)
3. **Command is published** to RabbitMQ queue `commands.{CommandName}`
4. **Consumer service** receives the command from RabbitMQ
5. **Command handler** processes the command and modifies the database
6. **Failed messages** are automatically routed to dead letter queues configured via the custom topology builder
7. **Events can be published** for further processing (e.g., notifications, search indexing)

### Query Flow (Read Operations)

1. **Client sends GET request** to a query endpoint (e.g., `/posts`, `/posts/:id`, `/posts/slug/:slug`)
2. **Server creates a query object** (e.g., `GetPostQuery`, `FindAllQuery`, `FindBySlugQuery`)
3. **Query is executed** synchronously through the Query Bus
4. **Query handler** retrieves data from the repository
5. **Response is returned** immediately to the client

**Key Difference:** Queries are handled synchronously for immediate responses, while commands are processed asynchronously via RabbitMQ for better scalability and decoupling.

### Message Queues

- **Commands**: `commands.{CommandName}` (e.g., `commands.CreatePostCommand`, `commands.CreateUserCommand`)
- **Events**: `events.{EventName}` (e.g., `events.PostCreated`)
- **Dead Letter Queue**: `{QueueName}.{DLQ_SUFFIX}` - Failed messages that cannot be processed are automatically routed here

### Dead Letter Queue

The application uses a **custom topology builder** to configure RabbitMQ dead letter exchanges and queues. When a message is negatively acknowledged (nacked) or cannot be processed, RabbitMQ automatically routes it to the corresponding dead letter queue.

**How it works:**
1. A dead letter exchange (DLX) is declared for each queue topology
2. A dead letter queue (DLQ) is created for each main queue with the naming pattern: `{QueueName}.{DLQ_SUFFIX}`
3. The DLQ is bound to the DLX exchange with a routing key
4. Main queues are configured with `x-dead-letter-exchange` and `x-dead-letter-routing-key` arguments
5. When a message is nacked with `requeue=false`, RabbitMQ routes it to the DLX, which then routes it to the DLQ

**Configuration:**
- Dead letter queues are automatically created when queues are declared
- Each queue gets its own dedicated dead letter queue
- Messages in DLQs can be inspected, republished, or manually handled via RabbitMQ Management UI

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

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT,
    provider VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    provider_user_id VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    access_token TEXT NOT NULL,
    access_token_secret TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    id_token TEXT NOT NULL
);
```

**Schema Details:**
- **id**: UUID type for proper unique identifier handling
- **created_at/updated_at**: TIMESTAMPTZ (timestamp with timezone) for accurate datetime tracking
- **email**: VARCHAR(255) with UNIQUE constraint for user email addresses
- **password**: TEXT type for password storage (nullable, used for local authentication)
- **provider**: VARCHAR(255) for OAuth provider name (e.g., "github", "google")
- **name**: VARCHAR(255) for user's full name (optional)
- **first_name/last_name**: VARCHAR(255) for user's first and last names (optional)
- **provider_user_id**: VARCHAR(255) for the user's ID from the OAuth provider
- **avatar_url**: VARCHAR(500) for user's profile picture URL (optional)
- **access_token/access_token_secret**: TEXT type for OAuth access tokens
- **refresh_token**: TEXT type for OAuth refresh token
- **expires_at**: TIMESTAMPTZ for token expiration time
- **id_token**: TEXT type for OAuth ID token

**Indexes:**
- `idx_users_email`: Index on email for fast lookups
- `idx_users_provider_user_id`: Composite index on (provider, provider_user_id) for OAuth lookups

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
| `AMQP_URL` | RabbitMQ connection string | `amqp://guest:guest@rabbitmq:5672/` (Docker) or `amqp://guest:guest@localhost:5672/` (local) |
| `AMQP_DLX_EXCHANGE` | Dead letter exchange name | `my-dlx` (default if not set) |
| `AMQP_DLX_QUEUE_SUFFIX` | Suffix for dead letter queue names | `dlq` (default if not set) |
| `AMQP_DLX_ROUTING_KEY_SUFFIX` | Suffix for dead letter routing keys | `dlq` (default if not set) |
| `GITHUB_CLIENT_ID` | GitHub OAuth client ID | Required for OAuth authentication |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth client secret | Required for OAuth authentication |
| `SESSION_KEY` | Session encryption key | Required for session management |
| `SESSION_NAME` | Session cookie name | Required for session management |
| `CLIENT_URL` | Frontend client URL for OAuth redirects | Required for OAuth callbacks |

## Dependencies

### Core Libraries

- **Gin**: HTTP web framework
- **Watermill**: Event-driven architecture library
- **Watermill-AMQP**: RabbitMQ (AMQP) integration for Watermill
- **golang-migrate**: Database migration tool
- **PostgreSQL Driver**: Database connectivity

### Architecture Libraries

- **CQRS**: Command Query Responsibility Segregation implementation (Watermill)
- **Query Bus**: Custom synchronous query handling implementation
- **Topology Builder**: Custom RabbitMQ topology builder for dead letter queue configuration
- **UUID**: Unique identifier generation (Google UUID library)
- **Goth**: OAuth authentication library for multiple providers
- **Gorilla Sessions**: Session management for authenticated users

## Docker Services

The `docker-compose.yaml` file defines the following services:

- **postgres**: PostgreSQL 18 database
- **rabbitmq**: RabbitMQ message broker with management UI (ports 5672 and 15672)
- **migrate**: Database migration service (runs once)
- **server**: HTTP API server
- **consume**: RabbitMQ consumer service

## Stopping Services

```bash
docker-compose down
```

To remove volumes (database data):

```bash
docker-compose down -v
```

## Troubleshooting

### RabbitMQ Connection Issues

If the consumer can't connect to RabbitMQ:

1. Verify RabbitMQ is running: `docker-compose ps rabbitmq`
2. Check RabbitMQ logs: `docker-compose logs rabbitmq`
3. Ensure `AMQP_URL` environment variable is correct
4. Access RabbitMQ Management UI at `http://localhost:15672` (default: guest/guest) to inspect queues and messages

### Dead Letter Queue

To inspect failed messages:

1. Access RabbitMQ Management UI at `http://localhost:15672`
2. Navigate to the "Exchanges" tab to see the dead letter exchange (default: `my-dlx`)
3. Navigate to the "Queues" tab
4. Look for queues with the pattern `{QueueName}.{DLQ_SUFFIX}` (e.g., `commands.CreatePostCommand.dlq`)
5. Click on a DLQ to see failed messages
6. Messages in the DLQ can be:
   - Inspected (view message payload and headers)
   - Republished to the original queue
   - Deleted
   - Manually routed to another queue

**Note:** Dead letter queues are automatically created when their corresponding main queues are declared. Each queue has its own dedicated dead letter queue.

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

