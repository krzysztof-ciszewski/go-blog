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
- **FindAllByQuery**: Retrieve posts with advanced filtering and pagination
  - **Pagination**: Supports `page` and `pageSize` parameters
  - **Filtering**: Supports multiple filter combinations:
    - `slug`: Filter by post slug (partial match)
    - `text`: Search in post title and content (partial match)
    - `author`: Filter by author name (partial match)
  - Filters can be combined (e.g., filter by text AND author)
  - Returns paginated results with total count, current page, and page size

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
- **RESTful API**: HTTP endpoints for blog operations with advanced filtering, search, and pagination
- **OAuth Authentication**: GitHub OAuth integration for user authentication
- **User Management**: User entity with OAuth token storage
- **PostgreSQL**: Persistent data storage with proper data types
- **Database Migrations**: Version-controlled schema changes

## Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose
- PostgreSQL 18 (if running manually)
- RabbitMQ 3 (if running manually)
- GitHub OAuth App (for authentication)
- SQLite (for testing - used by Watermill for command/event storage)

## Getting Started

### Using Docker Compose (Recommended)

1. **Create a `.env` file** in the project root with the required environment variables:
   ```bash
   POSTGRES_USER=blog
   POSTGRES_PASSWORD=blogpassword
   POSTGRES_DB=blog
   RABBITMQ_USER=guest
   RABBITMQ_PASSWORD=guest
   DATABASE_URL=postgres://blog:blogpassword@postgres:5432/blog?sslmode=disable
   AMQP_URL=amqp://guest:guest@rabbitmq:5672/
   AMQP_DLX_EXCHANGE=my-dlx
   AMQP_DLX_QUEUE_SUFFIX=dlq
   AMQP_DLX_ROUTING_KEY_SUFFIX=dlq
   GITHUB_CLIENT_ID=your_github_client_id
   GITHUB_CLIENT_SECRET=your_github_client_secret
   SESSION_SECRET=your_32_byte_or_longer_secret_key
   SESSION_NAME=blog_session
   API_URL=http://localhost:8080
   CLIENT_URL=http://localhost:3000
   ```

2. **Start all services:**
   ```bash
   docker compose up -d
   ```

   This will start:
   - PostgreSQL database (port 5432)
   - RabbitMQ message broker (ports 5672 and 15672 for management UI)
   - Database migration service (runs once)
   - API server (port 8080)
   - RabbitMQ consumer service

   **Note:** A test service is also available but runs only when explicitly started with the `test` profile.

3. **Check service status:**
   ```bash
   docker compose ps
   ```

4. **View logs:**
   ```bash
   docker compose logs -f server
   docker compose logs -f consume
   ```

5. **Access RabbitMQ Management UI:**
   - URL: `http://localhost:15672`
   - Default credentials: `guest` / `guest`

### Observability (Optional)

This repo includes an optional local OpenTelemetry “LGTM” stack (Grafana + Tempo + Loki + Mimir) via the `otel-lgtm` service in `docker-compose.yaml`.

Start it with:

```bash
docker compose up -d otel-lgtm
```

Then access:
- Grafana: `http://localhost:3000`

### Manual Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start dependencies:**
   ```bash
   docker compose up -d postgres rabbitmq
   ```

   Or install and run PostgreSQL and RabbitMQ manually.

3. **Set environment variables:**

   Create a `.env` or `.env.local` file with the following variables:
   ```bash
   DATABASE_URL="postgres://blog:blogpassword@localhost:5432/blog?sslmode=disable"
   AMQP_URL="amqp://guest:guest@localhost:5672/"
   AMQP_DLX_EXCHANGE="my-dlx"
   AMQP_DLX_QUEUE_SUFFIX="dlq"
   AMQP_DLX_ROUTING_KEY_SUFFIX="dlq"
   GITHUB_CLIENT_ID="your_github_client_id"
   GITHUB_CLIENT_SECRET="your_github_client_secret"
   SESSION_SECRET="your_32_byte_or_longer_secret_key"
   SESSION_NAME="blog_session"
   API_URL="http://localhost:8080"
   CLIENT_URL="http://localhost:3000"
   ```

   The migration tool automatically loads `.env.local` and `.env` files. For the server and consumer, you may need to export them in your shell:
   ```bash
   export DATABASE_URL="postgres://blog:blogpassword@localhost:5432/blog?sslmode=disable"
   export AMQP_URL="amqp://guest:guest@localhost:5672/"
   export AMQP_DLX_EXCHANGE="my-dlx"
   export AMQP_DLX_QUEUE_SUFFIX="dlq"
   export AMQP_DLX_ROUTING_KEY_SUFFIX="dlq"
   export GITHUB_CLIENT_ID="your_github_client_id"
   export GITHUB_CLIENT_SECRET="your_github_client_secret"
   export SESSION_SECRET="your_32_byte_or_longer_secret_key"
   export SESSION_NAME="blog_session"
   export API_URL="http://localhost:8080"
   export CLIENT_URL="http://localhost:3000"
   ```

4. **Run database migrations:**
   ```bash
   go run cmd/migrate.go
   ```

5. **Start the API server:**
   ```bash
   go run cmd/server.go
   ```

   The server will start on port 8080 by default.

6. **Start the consumer service (in a separate terminal):**
   ```bash
   go run cmd/consume.go
   ```

## API Documentation

The complete API specification is available in OpenAPI 3.0 format at [`docs/openapi.json`](docs/openapi.json).

**Key Points:**
- All API endpoints are prefixed with `/api/v1` and require authentication via session cookies (except OAuth endpoints)
- Authentication is handled through GitHub OAuth, and a session cookie is set after successful login
- Write operations (POST, DELETE) are processed asynchronously via RabbitMQ
- Read operations (GET) are handled synchronously through the Query Bus for immediate responses

You can use tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) or [Postman](https://www.postman.com/) to import and explore the OpenAPI specification.

## How It Works

### Command Flow (Write Operations)

1. **Client sends POST/DELETE request** to `/api/v1/posts` or `/api/v1/posts/:id` endpoint, OR
   **OAuth callback** triggers user creation via `/auth/:provider/callback`
2. **Server validates** the request and creates a command (`CreatePostCommand`, `DeletePostCommand`, or `CreateUserCommand`)
3. **Command is published** to RabbitMQ queue `commands.{CommandName}`
4. **Consumer service** receives the command from RabbitMQ
5. **Command handler** processes the command and modifies the database
6. **Failed messages** are automatically routed to dead letter queues configured via the custom topology builder
7. **Events can be published** for further processing (e.g., notifications, search indexing)

### Query Flow (Read Operations)

1. **Client sends GET request** to a query endpoint (e.g., `/api/v1/posts`, `/api/v1/posts/:id`)
2. **Server creates a query object** (e.g., `GetPostQuery`, `FindAllByQuery`)
3. **Query is executed** synchronously through the Query Bus
4. **Query handler** retrieves data from the repository with optional filtering and pagination
5. **Response is returned** immediately to the client

**Key Difference:** Queries are handled synchronously for immediate responses, while commands are processed asynchronously via RabbitMQ for better scalability and decoupling.

**Filtering and Pagination:**
- The `FindAllByQuery` supports multiple filter parameters (slug, text, author) that can be combined
- Filters use partial matching (LIKE queries) for flexible searching
- Pagination is handled at the repository level with proper offset/limit calculations
- Response includes total count for building pagination UI

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

## Development

### Running Tests

#### Using Docker Compose (Recommended)

Run tests in an isolated Docker container with all dependencies:

```bash
docker compose build test; docker compose run --remove-orphans test
```

This will:
- Start PostgreSQL and RabbitMQ services
- Run database migrations
- Execute all tests in `./internal/...`
- Use SQLite for Watermill command/event storage (test isolation)

#### Running Tests Locally

1. **Start dependencies:**
   ```bash
   docker compose up -d postgres rabbitmq
   ```

2. **Run database migrations:**
   ```bash
   go run cmd/migrate.go
   ```

3. **Run tests:**
   ```bash
   go test ./...
   ```

   Or run tests for a specific package:
   ```bash
   go test ./internal/UserInterface/Api/Handler/Post/...
   ```

#### Test Infrastructure

The project uses:
- **testify/suite**: For organized test suites
- **testify/assert**: For assertions
- **Table-driven tests**: Comprehensive test coverage using table-driven test patterns for pagination and filtering scenarios
- **SQLite**: For Watermill command/event storage in tests (via `watermill-sqlite`)
- **Test DI Container**: Custom dependency injection container for tests (see `internal/Infrastructure/DependencyInjection/Test/`)

**Test Coverage:**
- Handler tests include comprehensive table-driven tests for:
  - Pagination scenarios (different page sizes, page numbers, edge cases)
  - Multi-field filtering (combinations of slug, text, and author filters)
  - Pagination combined with filtering
  - Empty result scenarios
  - Invalid input handling

For detailed testing guidelines, see [`docs/HANDLER_TEST_GUIDELINES.md`](docs/HANDLER_TEST_GUIDELINES.md).

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
| `SESSION_SECRET` | Session encryption key (32+ bytes) | Required for session management |
| `SESSION_NAME` | Session cookie name | Required for session management |
| `API_URL` | Base URL of the API server | Required for OAuth callback URLs |
| `CLIENT_URL` | Frontend client URL for OAuth redirects | Required for OAuth callbacks |
| `POSTGRES_USER` | PostgreSQL database user | `blog` (Docker Compose) |
| `POSTGRES_PASSWORD` | PostgreSQL database password | `blogpassword` (Docker Compose) |
| `POSTGRES_DB` | PostgreSQL database name | `blog` (Docker Compose) |
| `RABBITMQ_USER` | RabbitMQ username | `guest` (Docker Compose) |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `guest` (Docker Compose) |

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
- **Goth**: OAuth authentication library for multiple providers (GitHub)
- **Gorilla Sessions**: Session management for authenticated users
- **Godotenv**: Environment variable management

### Testing Libraries

- **testify/suite**: Test suite organization
- **testify/assert**: Test assertions
- **watermill-sqlite**: SQLite integration for Watermill (used in tests for command/event storage)

## Docker Services

The `docker-compose.yaml` file defines the following services:

- **postgres**: PostgreSQL 18 database
- **rabbitmq**: RabbitMQ message broker with management UI (ports 5672 and 15672)
- **migrate**: Database migration service (runs once)
- **server**: HTTP API server
- **consume**: RabbitMQ consumer service
- **test**: Test runner service (runs with `--profile test` flag)

## Stopping Services

```bash
docker compose down
```

To remove volumes (database data):

```bash
docker compose down -v
```

To stop services including test profile:

```bash
docker compose --profile test down
```

## Troubleshooting

### RabbitMQ Connection Issues

If the consumer can't connect to RabbitMQ:

1. Verify RabbitMQ is running: `docker compose ps rabbitmq`
2. Check RabbitMQ logs: `docker compose logs rabbitmq`
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

1. Verify PostgreSQL is running: `docker compose ps postgres`
2. Check database logs: `docker compose logs postgres`
3. Verify `DATABASE_URL` environment variable

### Migration Issues

1. Check migration logs: `docker compose logs migrate`
2. Ensure migrations are in `db/migrations/` directory
3. Verify migration file naming convention

## License

This project is provided as-is for educational and development purposes.

