# Service Catalog API

A RESTful API service for managing service catalog information. Built with Go, Gin, and PostgreSQL.

## Features

- List services with pagination, sorting, and filtering
- Get service details with version count
- Get detailed version history for a service
- RESTful API design
- Proper error handling and validation
- SQL query optimization

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: 
  - **Gin Web Framework**. Why Gin? Gin is a high-performance HTTP framework that simplifies building RESTful APIs and 
  handle a large number of concurrent requests with minimal latency. It includes built-in middleware for tasks like logging, error handling 
  and authentication, which reduces the need for additional boilerplate code.
    -  Built-In Features which comes in Gin: 
      - Query Parameters, Path Parameters, and Form Parsing: Automatically parses query strings and form data, making it easier to handle user inputs. 
      - Data Binding and Validation: Bind JSON, XML, or form inputs directly to structs with validation rules. 
      - Rendering: Supports JSON, XML, HTML, and custom data rendering out of the box.
  - **Viper**: Viper is beneficial for configuration management in Go because it:
    - Offers type-safe value binding
    - Enables hierarchical configuration with overrides
    - Supports multiple config formats (JSON, YAML, TOML, ENV)
    - Can read from environment variables
    - Provides live config watching/reloading
    - Database connection string and other sensitive information (e.g., API keys) to environment variables or a configuration file.
- **ORM**: GORM
  - Why GORM?
    GORM abstracts SQL queries, making it easier to interact with the database. 
  - It supports multiple databases (e.g., PostgreSQL, MySQL, SQLite) and provides features like migrations, relationships, and hooks.
- **Database**: PostgreSQL
  - Why PostgreSQL?
    PostgreSQL is a robust, production-grade database that supports advanced features like transactions, full-text search, and JSONB. 
  - It ensures high consistency, scalability, and performance making it suitable for business-critical systems.
- **Packages**:
    - `github.com/gin-gonic/gin`: HTTP web framework
    - `gorm.io/gorm`: ORM
    - `gorm.io/driver/postgres`: PostgreSQL driver
    - `github.com/spf13/viper`:  Viper - configuration solution for Go
    - `github.com/stretchr/testify`: Testify -  toolkit with common assertions and mocks
    - `github.com/sirupsen/logrus`: Logrus 
    - Custom validation package

- **Assumptions**:
  - Database Choice: We have not set up a scalable RDS (Relational Database Service) and are using PostgreSQL for its advanced features and reliability. 
  - Authentication: Authentication and authorization are not implemented in this iteration but will be added in the future.
  - Scalability: The current implementation focuses on core functionality. Scalability improvements (e.g., caching, connection pooling) will be addressed in future iterations.

## Implementation Details
- Database Schema: The database includes tables for services and versions, with a one-to-many relationship between them.

- Soft Deletion: Services and versions support soft deletion, allowing users to restore deleted items if needed.

- Pagination and Filtering: The /services endpoint supports pagination, sorting and filtering to handle large datasets efficiently.

- Error Handling: Errors are handled gracefully, with meaningful messages returned to the user.
- Middleware: Logger middleware to track and log all API requests.
- Inject the database connection into handlers using dependency injection
  - Injects database dependency through constructor
  - Makes testing easier with mock databases
  - Separates router setup for better maintainability
  - Follows SOLID principles
- Unit and Integration Tests:
  - Added unit tests for handlers, models, and validation logic.
  - Integration tests to test API endpoints against a real database `servicecatalog_test`.
  - Use a testing framework like testify for assertions and mocking.

## Database Schema

```sql
CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE versions (
    id SERIAL PRIMARY KEY,
    service_id INTEGER REFERENCES services(id),
    number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## API Documentation

### 1. List Services

GET /services

Query Parameters:
```
page: int (default: 1)
pageSize: int (default: 10)
search: string
sortBy: string (id, name, description)
sortDir: string (asc, desc)
```

Success Response (200 OK):
```json
{
    "services": [
        {
            "id": 1,
            "name": "Authentication Service",
            "description": "Handles authentication",
            "versions": 3
        }
    ],
    "total_count": 10,
    "current_page": 1,
    "page_size": 10
}
```

### 2. Get Service Details

GET /services/:id

Success Response (200 OK):
```json
{
    "id": 1,
    "name": "Authentication Service",
    "description": "Handles authentication",
    "versions": 3
}
```

Error Responses:
```json
// 400 Bad Request
{
    "error": "Invalid service ID: must be a positive integer",
    "details": "strconv.ParseUint: parsing \"abc\": invalid syntax"
}

// 404 Not Found
{
    "error": "Service not found",
    "service_id": 999
}
```

### 3. Get Service Versions

GET /services/:id/versions

Success Response (200 OK):
```json
{
    "service_id": 1,
    "versions": [
        {
            "id": 1,
            "number": "2.0.0",
            "created_at": "2024-01-20T10:00:00Z"
        }
    ]
}
```

## Project Structure

```
servicecatalog/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   ├── config.go
│   └── config.yaml
├── internal/
│   ├── database/
│   │   ├── setup.sql
│   │   ├── test_setup.sql
│   │   └── database.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   ├── service_get.go
│   │   ├── service_list.go
│   │   ├── service_versions.go
│   │   ├── types.go
│   │   └── handlers_test.go
│   ├── models/
│   │   ├── models.go
│   │   ├── models_test.go
│   │   └── version.go
│   └── validation/
│       ├── validation.go
│       └── validation_test.go
└── go.mod
```

## Setup and Installation

1. Clone the repository
```bash
git clone github.com/yourusername/servicecatalog
```

2. Install dependencies
```bash
go mod tidy
```

3. Set up PostgreSQL database
```bash
psql -U postgres -f database/setup.sql
```

4. Run the application
```bash
go run cmd/api/main.go
```

## Testing

Run all tests:
```bash
psql -U postgres -f db/test_setup.sql
go test ./... -v
go test ./...
```

## Performance Considerations

- Uses optimized SQL queries with JOINs to minimize database calls

- Implements proper indexing on frequently queried columns
  ```
  CREATE INDEX idx_services_name ON services(name);  -- For search functionality

  CREATE INDEX idx_versions_service_id ON versions(service_id);  -- For JOIN operations
  ```

  - services(name): Required for search functionality in ListServices 
  - versions(service_id): Required for efficient version counting and retrieval


- Uses connection pooling with GORM. Connection pooling configs like 
  - **SetMaxIdleConns**: Sets the maximum number of idle connections in the pool. Idle connections are connections that are not currently in use but are kept open for reuse. A higher value can improve performance by reducing the overhead of establishing new connections. 
  - **SetMaxOpenConns**: Sets the maximum number of open connections to the database. This limits the total number of concurrent connections that can be established. 
  - **SetConnMaxLifetime**: Sets the maximum amount of time a connection can be reused. After this duration, the connection is closed and a new one is created. This helps to gracefully handle database restarts or network changes.
- Implements pagination to handle large datasets efficiently

## Error Handling

The API uses standard HTTP status codes:
- 200: Success
- 400: Bad Request (validation errors)
- 404: Not Found
- 500: Internal Server Error

All error responses follow the format:
```json
{
    "error": "Error message",
    "details": "Additional error details (optional)"
}
```

## Future Improvements
If given more time, the following enhancements could be made:

- Rate limiting to prevent abuse of the GET services / versions API. We can use this framework [github.com/ulule/limiter]() 
to integrate well with Gin Framework in middleware.
- **Caching**: Integrate Redis to cache frequently accessed data (e.g., service details, version lists) and improve performance.

- **Monitoring**: Use Prometheus and Grafana to monitor API performance, error rates, and resource usage.

- **Search**: Implement full-text search using PostgreSQL's tsvector and tsquery for advanced search capabilities.

- **Asynchronous Processing**: Use a message queue (e.g., RabbitMQ) to handle long-running tasks like notifications or report generation.

- **Event Sourcing**: Track all changes to services and versions for auditability and state reconstruction.

- **Webhooks**: Notify external systems when services or versions are created, updated, or deleted.

- **Multi-Tenancy**: Support multiple organizations with isolated data.