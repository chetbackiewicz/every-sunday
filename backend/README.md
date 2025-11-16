# Every Sunday Backend

Backend API server for the Every Sunday budget tracking application.

## Technology Stack

- **Go 1.21+** - Programming language
- **Gin** - Web framework
- **sqlx** - SQL database access
- **pgx/v5** - PostgreSQL driver
- **golang-migrate** - Database migrations
- **golang-jwt** - JWT authentication
- **bcrypt** - Password hashing

## Project Structure

```
backend/
├── cmd/server/          # Application entry point
│   └── main.go
├── internal/
│   ├── auth/            # JWT and password hashing
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Gin middleware
│   ├── models/          # Data models
│   └── repositories/    # Database access layer
├── db/migrations/       # SQL migration files
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- golang-migrate CLI (for running migrations)

### Installation

1. Install dependencies:
```bash
go mod download
```

2. Install golang-migrate:
```bash
# macOS
brew install golang-migrate

# Or download from https://github.com/golang-migrate/migrate
```

3. Create PostgreSQL database:
```bash
createdb every_sunday_dev
```

4. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=every_sunday_dev
export DB_SSLMODE=disable
export PORT=8080
export JWT_ACCESS_SECRET=your-access-secret
export JWT_REFRESH_SECRET=your-refresh-secret
```

5. Run migrations:
```bash
migrate -path db/migrations -database "postgres://localhost/every_sunday_dev?sslmode=disable" up
```

6. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on http://localhost:8080

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user (protected)

### Budgets
- `GET /api/v1/budgets` - List all budgets (protected)
- `GET /api/v1/budgets/:month` - Get budget by month (protected)
- `POST /api/v1/budgets` - Create new budget (protected)
- `PUT /api/v1/budgets/:month` - Update budget (protected)
- `DELETE /api/v1/budgets/:month` - Delete budget (protected)

### Categories
- `GET /api/v1/budgets/:month/categories` - Get category tree (protected)
- `POST /api/v1/budgets/:month/categories` - Create category (protected)
- `PUT /api/v1/budgets/:month/categories/:id` - Update category (protected)
- `DELETE /api/v1/budgets/:month/categories/:id` - Delete category (protected)

### Files & Transactions
- `POST /api/v1/budgets/:month/files` - Upload CSV file (protected)
- `GET /api/v1/budgets/:month/transactions` - List transactions (protected)
- `PUT /api/v1/budgets/:month/transactions/:id` - Update transaction (protected)
- `DELETE /api/v1/budgets/:month/transactions/:id` - Delete transaction (protected)

### Cell References
- `POST /api/v1/budgets/:month/references` - Create cell reference (protected)
- `DELETE /api/v1/budgets/:month/references/:id` - Delete cell reference (protected)

## Database Migrations

Create new migration:
```bash
migrate create -ext sql -dir db/migrations -seq migration_name
```

Run migrations:
```bash
migrate -path db/migrations -database "postgres://localhost/every_sunday_dev?sslmode=disable" up
```

Rollback migration:
```bash
migrate -path db/migrations -database "postgres://localhost/every_sunday_dev?sslmode=disable" down 1
```

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/server cmd/server/main.go
```

## Research Documentation

All implementation decisions are backed by comprehensive research documentation in the `/research` directory:

- **01-csv-parsing.md** - CSV file handling
- **02-editable-tables.md** - Frontend table requirements
- **03-category-tree.md** - Hierarchical categories
- **04-cell-linking.md** - Transaction linking
- **08-golang-framework.md** - Gin framework selection
- **09-postgresql-schema-hierarchical.md** - Closure table pattern
- **10-postgresql-schema-budget-tables.md** - Database schema
- **11-postgresql-migrations.md** - Migration strategy
- **12-api-design.md** - API design patterns
- **13-database-integration.md** - sqlx and connection pooling
- **14-auth.md** - JWT and bcrypt authentication

## TODO

- [ ] Implement authentication handlers
- [ ] Implement repository methods
- [ ] Add CSV file parsing
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add API documentation (Swagger)
- [ ] Add logging
- [ ] Add request validation
- [ ] Add error handling improvements
