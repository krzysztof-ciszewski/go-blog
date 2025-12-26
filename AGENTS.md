# Agent Guidelines

This document provides guidelines for AI agents working on this codebase. It covers how to run tests, how to add tests, and how to run the application - all using Docker.

## Table of Contents

1. [Running Tests](#running-tests)
2. [Adding Tests](#adding-tests)
3. [Running the Application](#running-the-application)
4. [Docker Commands Reference](#docker-commands-reference)

## Running Tests

### Prerequisites

- Docker and Docker Compose installed
- `.env.test` file configured (or use default values)

### Quick Start

Run all tests using Docker Compose:

```bash
docker-compose build test; docker compose run --remove-orphans test
```

This command will:
- Build the test container image
- Start PostgreSQL and RabbitMQ services
- Run database migrations automatically
- Execute all tests in `./internal/...` with race detection enabled
- Clean up containers after tests complete

### Test Environment Variables

The test service uses `.env.test` file. If not present, you can create one based on the main `.env` file. Required variables:

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
GITHUB_CLIENT_ID=test_client_id
GITHUB_CLIENT_SECRET=test_client_secret
SESSION_SECRET=test_session_secret_32_bytes_or_longer
SESSION_NAME=blog_session
API_URL=http://localhost:8080
CLIENT_URL=http://localhost:3000
```

### Running Specific Test Packages

To run tests for a specific package, you can modify the Dockerfile.test command or run tests locally after starting dependencies:

```bash
# Start dependencies
docker-compose up -d postgres rabbitmq

# Run migrations
docker-compose up migrate

# Run specific test package (requires local Go installation)
go test ./internal/UserInterface/Api/Handler/Post/...
```

### Test Output

The test container will output:
- Test execution results
- Race detector warnings (if any)
- Test coverage information
- Exit code 0 on success, non-zero on failure

### Viewing Test Logs

Test output is displayed in real-time when running tests. To view logs after execution:

```bash
docker compose logs test
```

## Adding Tests

### Test Structure

This project uses:
- **testify/suite**: For organized test suites
- **testify/assert**: For assertions
- **Table-driven tests**: For comprehensive coverage
- **Test DI Container**: Custom dependency injection for tests

### Test File Naming

- Test files must end with `_test.go`
- Place test files in the same package as the code being tested
- Example: `create_post_test.go` tests `create_post.go`

### Test Suite Pattern

Follow this structure for handler tests:

```go
package post

import (
    "bytes"
    "database/sql"
    test "main/internal/Infrastructure/DependencyInjection/Test"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/ThreeDotsLabs/watermill/components/cqrs"
    "github.com/stretchr/testify/suite"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

type YourHandlerTestSuite struct {
    suite.Suite
    CommandBus *cqrs.CommandBus
    QueryBus   query_bus.QueryBus  // If needed
    Ctx        *gin.Context
    W          *httptest.ResponseRecorder
    PubSubDb   *sql.DB
}

func (s *YourHandlerTestSuite) SetupTest() {
    // Get test container
    s.CommandBus = test.GetTestContainer().CommandBus
    s.QueryBus = test.GetTestContainer().QueryBus
    
    // Setup HTTP test context
    s.W = httptest.NewRecorder()
    s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
    gin.SetMode(gin.TestMode)
    
    // Get PubSub database for command verification
    s.PubSubDb = test.GetPubSubDb()
    
    // Clean up command tables for test isolation
    s.PubSubDb.Exec("DELETE FROM `watermill_commands.{commandName}`")
    
    // Setup test data if needed
    test.GetTestContainer().DB.Exec("DELETE FROM posts")
    // Insert test data...
}

func (s *YourHandlerTestSuite) TestYourHandler() {
    // Create request
    s.Ctx.Request = httptest.NewRequest(
        "POST",
        "/api/v1/endpoint",
        bytes.NewBufferString(`{"field":"value"}`),
    )
    s.Ctx.Request.Header.Set("Content-Type", "application/json")
    
    // Call handler
    YourHandler(s.Ctx, s.CommandBus, s.QueryBus)
    
    // Assertions
    assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
    assert.Equal(s.T(), `{"message":"Success"}`, s.W.Body.String())
    
    // Verify command was sent
    count := test.GetCommandCount("{commandName}")
    assert.Equal(s.T(), 1, count)
}

func TestYourHandlerTestSuite(t *testing.T) {
    suite.Run(t, new(YourHandlerTestSuite))
}
```

### Test Guidelines

1. **Test Isolation**: Always clean up command tables and test data in `SetupTest()`
2. **Test Coverage**: Include:
   - Happy path (successful request)
   - Validation errors (each validation rule separately)
   - Edge cases (empty strings, null values)
   - Error scenarios (missing params, invalid IDs)
3. **Naming**: Use descriptive test method names: `Test{HandlerName}{Scenario}`
4. **Assertions**: 
   - Use `assert.Equal()` for status codes and response bodies
   - Use `test.GetCommandCount()` to verify commands were sent
5. **Test Data**: Use realistic but clearly test data (e.g., `"testslug"`, `"testtitle"`)

### Example Test Scenarios

**Successful Request:**
```go
func (s *YourHandlerTestSuite) TestCreatePost() {
    // Valid request
    // Assert StatusAccepted (202) or StatusOK (200)
    // Assert success message
    // Assert command count is 1
}
```

**Validation Error:**
```go
func (s *YourHandlerTestSuite) TestCreatePostInvalidSlug() {
    // Invalid slug in request
    // Assert StatusBadRequest (400)
    // Assert error message
    // Assert command count is 0
}
```

**Not Found:**
```go
func (s *YourHandlerTestSuite) TestGetPostNotFound() {
    // Request with non-existent ID
    // Assert StatusNotFound (404)
}
```

### Reference Files

- `docs/HANDLER_TEST_GUIDELINES.md` - Detailed testing guidelines
- `internal/UserInterface/Api/Handler/Post/create_post_test.go` - Example POST handler test
- `internal/UserInterface/Api/Handler/Post/delete_post_test.go` - Example DELETE handler test
- `internal/UserInterface/Api/Handler/Post/list_posts_test.go` - Example GET handler with pagination tests

## Running the Application

### Prerequisites

- Docker and Docker Compose installed
- `.env` file configured with required environment variables

### Quick Start

Start all services:

```bash
docker-compose up -d
```

This will start:
- **PostgreSQL** (port 5432)
- **RabbitMQ** (ports 5672 and 15672 for management UI)
- **Migration service** (runs once automatically)
- **API Server** (port 8080)
- **Consumer service** (processes commands from RabbitMQ)

### Environment Variables

Create a `.env` file in the project root:

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

### Service Status

Check if services are running:

```bash
docker-compose ps
```

### Viewing Logs

View logs for specific services:

```bash
# API Server logs
docker-compose logs -f server

# Consumer logs
docker-compose logs -f consume

# Database logs
docker-compose logs -f postgres

# RabbitMQ logs
docker-compose logs -f rabbitmq
```

### Accessing Services

- **API Server**: `http://localhost:8080`
- **RabbitMQ Management UI**: `http://localhost:15672` (default: guest/guest)
- **PostgreSQL**: `localhost:5432`

### Stopping Services

Stop all services:

```bash
docker-compose down
```

Stop and remove volumes (database data):

```bash
docker-compose down -v
```

### Restarting Services

Restart a specific service:

```bash
docker-compose restart server
docker-compose restart consume
```

Rebuild and restart (after code changes):

```bash
docker-compose up -d --build server
docker-compose up -d --build consume
```

## Docker Commands Reference

### Test Commands

```bash
# Run all tests
docker-compose build test; docker compose run --remove-orphans test

# View test logs
docker compose logs test

# Clean up test containers
docker compose down
```

### Application Commands

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d postgres rabbitmq

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# View service status
docker-compose ps

# View logs
docker-compose logs -f [service_name]

# Rebuild and restart service
docker-compose up -d --build [service_name]

# Restart service
docker-compose restart [service_name]
```

### Service Names

- `postgres` - PostgreSQL database
- `rabbitmq` - RabbitMQ message broker
- `migrate` - Database migration service
- `server` - HTTP API server
- `consume` - RabbitMQ consumer service
- `test` - Test runner service

### Building Images

Build specific service images:

```bash
docker-compose build server
docker-compose build consume
docker-compose build test
```

Build all images:

```bash
docker-compose build
```

## Troubleshooting

### Tests Fail to Connect to Database

1. Ensure PostgreSQL is running: `docker-compose ps postgres`
2. Check database logs: `docker-compose logs postgres`
3. Verify `DATABASE_URL` in `.env.test`

### Tests Fail to Connect to RabbitMQ

1. Ensure RabbitMQ is running: `docker-compose ps rabbitmq`
2. Check RabbitMQ logs: `docker-compose logs rabbitmq`
3. Verify `AMQP_URL` in `.env.test`

### Application Services Won't Start

1. Check service logs: `docker-compose logs [service_name]`
2. Verify environment variables in `.env`
3. Ensure dependencies are healthy: `docker-compose ps`
4. Check if ports are already in use

### Migration Issues

1. Check migration logs: `docker-compose logs migrate`
2. Ensure migration files are in `db/migrations/`
3. Verify migration file naming convention

### Rebuilding After Code Changes

After modifying code, rebuild the affected service:

```bash
docker-compose up -d --build server
docker-compose up -d --build consume
```

For tests, rebuild and run:

```bash
docker-compose build test; docker compose run --remove-orphans test
```

## Best Practices for Agents

1. **Always use Docker**: Never run tests or the application locally without Docker unless explicitly debugging
2. **Check logs first**: When something fails, check Docker logs before making changes
3. **Test isolation**: Ensure tests clean up after themselves
4. **Follow patterns**: Use existing test files as templates
5. **Verify with Docker**: After adding tests, run them using `docker-compose build test; docker compose run --remove-orphans test`
6. **Environment variables**: Always use environment variables from `.env` or `.env.test` files
7. **Service dependencies**: Be aware of service startup order (postgres/rabbitmq → migrate → server/consume)

