# Handler Test Guidelines

This document provides guidelines for writing handler tests in this codebase, based on the patterns established in `create_post_test.go` and `delete_post_test.go`.

## Table of Contents
1. [Test Suite Structure](#test-suite-structure)
2. [Setup Method](#setup-method)
3. [Test Method Patterns](#test-method-patterns)
4. [Common Test Scenarios](#common-test-scenarios)
5. [Assertions](#assertions)
6. [Best Practices](#best-practices)

## Test Suite Structure

### Test Suite Type
- Use `testify/suite` for test organization
- Create a test suite struct that embeds `suite.Suite`
- Name the struct following the pattern: `{HandlerName}TestSuite` (e.g., `CreatePostTestSuite`, `DeletePostTestSuite`)

### Required Fields
Every test suite should include these standard fields:

```go
type {HandlerName}TestSuite struct {
    suite.Suite
    CommandBus *cqrs.CommandBus  // For command handlers
    QueryBus   *cqrs.QueryBus    // For query handlers (if applicable)
    Ctx        *gin.Context      // Gin context for handler
    W          *httptest.ResponseRecorder  // Response recorder
    PubSubDb   *sql.DB           // Database connection for command verification
}
```

### Additional Fields
Add fields specific to your test needs:
- `PostUuid uuid.UUID` - For tests that need a specific entity ID
- Any other test-specific data structures

## Setup Method

### Standard Setup Pattern
Every test suite must implement `SetupTest()` with the following pattern:

```go
func (s *{HandlerName}TestSuite) SetupTest() {
    // 1. Get CommandBus/QueryBus from test container
    s.CommandBus = test.GetTestContainer().CommandBus
    
    // 2. Create HTTP response recorder
    s.W = httptest.NewRecorder()
    
    // 3. Create Gin test context
    s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
    
    // 4. Set Gin to test mode
    gin.SetMode(gin.TestMode)
    
    // 5. Get PubSubDb for command verification
    s.PubSubDb = test.GetPubSubDb()
    
    // 6. Clean up command table (replace with your command name)
    s.PubSubDb.Exec("DELETE FROM `watermill_commands.{commandName}`")
    
    // 7. Set up any test data (e.g., insert test entities)
    // Example for delete tests:
    // postUuid, err := uuid.NewRandom()
    // if err != nil {
    //     panic(err)
    // }
    // s.PostUuid = postUuid
    // test.GetTestContainer().DB.Exec("INSERT INTO posts ...", postUuid.String())
}
```

### Important Notes
- Always clean up command tables in `SetupTest()` to ensure test isolation
- Use `gin.CreateTestContextOnly()` to create a test context
- Set `gin.TestMode` to avoid unnecessary logging during tests
- For tests requiring existing entities, insert them in `SetupTest()`

## Test Method Patterns

### Basic Test Method Structure

```go
func (s *{HandlerName}TestSuite) Test{ScenarioName}() {
    // 1. Create HTTP request
    s.Ctx.Request = httptest.NewRequest(
        "HTTP_METHOD",  // GET, POST, PUT, DELETE, etc.
        "/api/v1/endpoint",
        requestBody,    // bytes.Buffer or nil
    )
    
    // 2. Set request headers
    s.Ctx.Request.Header.Set("Content-Type", "application/json")
    // Add other headers as needed (e.g., Authorization)
    
    // 3. Set route parameters (for routes with path params)
    s.Ctx.Params = gin.Params{
        gin.Param{
            Key:   "paramName",
            Value: "paramValue",
        },
    }
    
    // 4. Call the handler
    {HandlerFunction}(s.Ctx, s.CommandBus)
    
    // 5. Assert response status code
    assert.Equal(s.T(), http.Status{ExpectedCode}, s.W.Code)
    
    // 6. Assert response body
    assert.Equal(s.T(), `{"expected":"json"}`, s.W.Body.String())
    
    // 7. Verify command was sent (for command handlers)
    count := test.GetCommandCount("{commandName}")
    assert.Equal(s.T(), expectedCount, count)
}
```

### Request Body Examples

**JSON Request Body:**
```go
bytes.NewBufferString(`{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "slug": "testslug",
    "title": "testtitle",
    "content": "testcontent",
    "author": "testauthor"
}`)
```

**No Body (for GET/DELETE):**
```go
nil
```

### Route Parameters Example
```go
s.Ctx.Params = gin.Params{
    gin.Param{
        Key:   "id",
        Value: s.PostUuid.String(),
    },
}
```

## Common Test Scenarios

### 1. Successful Request Test
Test the happy path where all validations pass:

```go
func (s *{HandlerName}TestSuite) Test{HandlerName}() {
    // Create valid request
    // Call handler
    // Assert StatusAccepted (202) or StatusOK (200)
    // Assert success message
    // Assert command count is 1
}
```

### 2. Validation Error Tests
Test each validation rule separately:

```go
func (s *{HandlerName}TestSuite) Test{HandlerName}Invalid{Field}Request() {
    // Create request with invalid field
    // Call handler
    // Assert StatusBadRequest (400)
    // Assert error message matches validation error
    // Assert command count is 0 (command should not be sent)
}
```

**Example validation tests:**
- `TestCreatePostInvalidSlugRequest` - Tests slug validation
- `TestCreatePostInvalidTitleRequest` - Tests title min length
- `TestCreatePostInvalidContentRequest` - Tests content min length

### 3. Missing Parameter Tests
For handlers that require path parameters:

```go
func (s *{HandlerName}TestSuite) Test{HandlerName}MissingId() {
    // Create request without required parameter
    // Call handler
    // Assert appropriate error status
}
```

### 4. Not Found Tests
For handlers that query entities:

```go
func (s *{HandlerName}TestSuite) Test{HandlerName}NotFound() {
    // Create request with non-existent ID
    // Call handler
    // Assert StatusNotFound (404)
}
```

## Assertions

### Status Code Assertions
- Use `assert.Equal(s.T(), http.Status{Code}, s.W.Code)`
- Common status codes:
  - `http.StatusAccepted` (202) - For async command handlers
  - `http.StatusOK` (200) - For query handlers
  - `http.StatusBadRequest` (400) - For validation errors
  - `http.StatusNotFound` (404) - For missing resources
  - `http.StatusUnauthorized` (401) - For auth failures
  - `http.StatusForbidden` (403) - For permission failures

### Response Body Assertions
- Use `assert.Equal(s.T(), expectedJSON, s.W.Body.String())`
- Match exact JSON strings for consistency
- Success responses: `{"message":"..."}`
- Error responses: `{"error":"..."}`

### Command Count Assertions
- Use `test.GetCommandCount("{commandName}")` to verify commands were sent
- Successful commands: `assert.Equal(s.T(), 1, count)`
- Failed/validation errors: `assert.Equal(s.T(), 0, count)`

## Best Practices

### 1. Test Isolation
- Always clean up command tables in `SetupTest()`
- Each test should be independent and not rely on other tests
- Use unique IDs/UUIDs for test data

### 2. Naming Conventions
- Test suite: `{HandlerName}TestSuite`
- Test methods: `Test{ScenarioDescription}()`
- Use descriptive names that explain what is being tested

### 3. Test Coverage
- Test the happy path (successful request)
- Test all validation rules separately
- Test edge cases (empty strings, null values, etc.)
- Test error scenarios (missing params, invalid IDs, etc.)

### 4. Code Organization
- Group related tests together
- Order tests logically (happy path first, then validation errors)
- Keep test methods focused on a single scenario

### 5. Test Data
- Use realistic but clearly test data
- Use UUIDs for IDs: `"123e4567-e89b-12d3-a456-426614174000"`
- Use descriptive test values: `"testslug"`, `"testtitle"`, `"testcontent"`

### 6. Imports
Standard imports for handler tests:
```go
import (
    "database/sql"
    test "main/internal/Infrastructure/DependencyInjection/Test"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/ThreeDotsLabs/watermill/components/cqrs"
    "github.com/stretchr/testify/suite"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    // Add other imports as needed (uuid, bytes, etc.)
)
```

### 7. Test Runner
Always include a test runner function:
```go
func Test{HandlerName}TestSuite(t *testing.T) {
    suite.Run(t, new({HandlerName}TestSuite))
}
```

## Example: Complete Test File Template

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

type {HandlerName}TestSuite struct {
    suite.Suite
    CommandBus *cqrs.CommandBus
    Ctx        *gin.Context
    W          *httptest.ResponseRecorder
    PubSubDb   *sql.DB
}

func (s *{HandlerName}TestSuite) SetupTest() {
    s.CommandBus = test.GetTestContainer().CommandBus
    s.W = httptest.NewRecorder()
    s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
    gin.SetMode(gin.TestMode)
    s.PubSubDb = test.GetPubSubDb()
    s.PubSubDb.Exec("DELETE FROM `watermill_commands.{commandName}`")
}

func (s *{HandlerName}TestSuite) Test{HandlerName}() {
    s.Ctx.Request = httptest.NewRequest(
        "POST",
        "/api/v1/endpoint",
        bytes.NewBufferString(`{"field":"value"}`),
    )
    s.Ctx.Request.Header.Set("Content-Type", "application/json")
    
    {HandlerFunction}(s.Ctx, s.CommandBus)
    
    assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
    assert.Equal(s.T(), `{"message":"Success message"}`, s.W.Body.String())
    count := test.GetCommandCount("{commandName}")
    assert.Equal(s.T(), 1, count)
}

func Test{HandlerName}TestSuite(t *testing.T) {
    suite.Run(t, new({HandlerName}TestSuite))
}
```

## References
- `internal/UserInterface/Api/Handler/Post/create_post_test.go` - Example of POST handler with validation tests
- `internal/UserInterface/Api/Handler/Post/delete_post_test.go` - Example of DELETE handler with path parameters

