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
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL
- **Packages**:
    - `github.com/gin-gonic/gin`: HTTP web framework
    - `gorm.io/gorm`: ORM
    - `gorm.io/driver/postgres`: PostgreSQL driver
    - Custom validation package

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
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   ├── types.go
│   │   └── utils.go
│   ├── models/
│   │   ├── service.go
│   │   └── version.go
│   └── validation/
│       └── service.go
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
psql -U postgres -f db/setup.sql
```

4. Run the application
```bash
go run cmd/api/main.go
```

## Testing

Run all tests:
```bash
go test ./...
```

## Performance Considerations

- Uses optimized SQL queries with JOINs to minimize database calls
- Implements proper indexing on frequently queried columns
- Uses connection pooling with GORM
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