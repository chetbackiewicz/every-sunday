# Database Integration Patterns Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: sqlx with Repository Pattern + Connection Pooling

---

## Recommendation Summary

**Chosen Approach: sqlx + Repository Pattern**

After evaluating GORM, sqlx, and sqlc, **sqlx with the repository pattern** is recommended for the budget tracking application.

**Key Decisions**:
1. ✅ **sqlx for database access** - Minimal abstraction, type-safe, explicit SQL
2. ✅ **Repository pattern** - Clean architecture with testable interfaces
3. ✅ **pgx/v5 connection pooling** - Efficient connection management
4. ✅ **Transaction support** - Manual BEGIN/COMMIT for complex operations
5. ✅ **Named parameters** - Struct binding with `:field` syntax
6. ✅ **StructScan** - Direct scanning into structs
7. ✅ **golang-migrate** - Schema migrations (from doc 11)

**Why sqlx over GORM and sqlc**:
- **vs GORM**: More explicit control, less magic, better performance, simpler debugging
- **vs sqlc**: No code generation step, more flexibility, easier iterative development

---

## ORM Comparison

### Option 1: GORM (Full-Featured ORM)

**Overview**: Most popular Go ORM with 37k+ GitHub stars, full ActiveRecord-like ORM.

**Pros**:
- Rich feature set: associations, hooks, soft deletes, preloading
- Automatic migrations: `db.AutoMigrate(&Model{})`
- Association handling: `db.Preload("Orders").Find(&users)`
- Query builder: `db.Where("name = ?", "john").Find(&user)`
- Built-in connection pooling
- Large community and ecosystem

**Cons**:
- **Performance overhead**: Reflection-heavy, slower than raw SQL
- **Magic behavior**: Implicit conventions can surprise developers
- **Debugging difficulty**: Complex error traces, hard to see actual SQL
- **Over-engineering risk**: Features we don't need (polymorphism, STI, etc.)
- **Breaking changes**: v2 migration was painful for many projects

**Code Example**:
```go
// GORM model
type User struct {
    gorm.Model
    Email    string `gorm:"unique;not null"`
    Password string
}

// Query
var user User
result := db.Where("email = ?", email).First(&user)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // Handle not found
}

// Association handling
db.Preload("Categories").Find(&budgets)
```

**Verdict**: ❌ **Too heavyweight** for our needs. We have simple CRUD operations and explicit SQL is better for financial calculations.

---

### Option 2: sqlx (Minimal Extension)

**Overview**: Extensions on top of `database/sql` with named parameters and struct scanning (16k+ stars).

**Pros**:
- **Minimal abstraction**: Thin wrapper over database/sql
- **Explicit SQL**: Write raw SQL, no query builder magic
- **Type safety**: StructScan maps columns to struct fields
- **Named parameters**: `:field` syntax with struct/map binding
- **Excellent performance**: Near-zero overhead
- **No code generation**: Simple iterative development
- **Context support**: QueryContext, ExecContext built-in
- **Transaction support**: db.Beginx(), tx.Commit()

**Cons**:
- **Manual SQL writing**: No query builder (but this is also a pro)
- **No associations**: Must handle foreign keys manually
- **No migrations**: Need separate tool (we have golang-migrate)
- **More boilerplate**: Repository pattern required for clean code

**Code Example**:
```go
// sqlx usage
type User struct {
    ID       int    `db:"id"`
    Email    string `db:"email"`
    Password string `db:"password"`
}

// Named query
var user User
err := db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
if errors.Is(err, sql.ErrNoRows) {
    // Handle not found
}

// Named parameters with struct
_, err = db.NamedExec(`
    INSERT INTO users (email, password) 
    VALUES (:email, :password)
`, user)

// Select multiple
var users []User
err = db.Select(&users, "SELECT * FROM users WHERE created_at > $1", date)
```

**Verdict**: ✅ **Best fit** for our application. Explicit, performant, and flexible.

---

### Option 3: sqlc (Code Generation)

**Overview**: Generates type-safe Go code from SQL queries (16k+ stars).

**Pros**:
- **Type safety**: Compile-time checking of SQL queries
- **No reflection**: Generated code is pure Go
- **Excellent performance**: Near-raw SQL performance
- **SQL-first**: Write SQL, get Go interfaces
- **PostgreSQL support**: Full PostgreSQL feature support

**Cons**:
- **Code generation step**: Must run sqlc after SQL changes
- **Build complexity**: Extra step in CI/CD pipeline
- **Less flexibility**: Generated code is rigid
- **Iterative development**: Slower feedback loop (write SQL → generate → test)
- **Custom queries**: Complex queries require manual editing

**Code Example**:
```sql
-- queries/users.sql
-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;
```

Generated Go code:
```go
// Automatically generated
type Queries struct {
    db *sql.DB
}

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
    // Generated implementation
}
```

**Verdict**: ⚠️ **Good but unnecessary** for our project. Extra build step adds complexity without clear benefit over sqlx.

---

## Decision Matrix

| Feature              | GORM      | sqlx         | sqlc         |
| -------------------- | --------- | ------------ | ------------ |
| **Explicit SQL**     | ❌         | ✅            | ✅            |
| **Performance**      | ⚠️ Medium  | ✅ Fast       | ✅ Fast       |
| **Type Safety**      | ⚠️ Runtime | ✅ StructScan | ✅ Generated  |
| **Learning Curve**   | ⚠️ Steep   | ✅ Easy       | ⚠️ Medium     |
| **Debugging**        | ❌ Hard    | ✅ Easy       | ✅ Easy       |
| **Flexibility**      | ⚠️ Limited | ✅ Full       | ⚠️ Limited    |
| **Build Complexity** | ✅ Simple  | ✅ Simple     | ❌ Generation |
| **Iterative Dev**    | ✅ Fast    | ✅ Fast       | ⚠️ Slow       |
| **Query Builder**    | ✅ Yes     | ❌ No         | ❌ No         |
| **Associations**     | ✅ Auto    | ❌ Manual     | ❌ Manual     |
| **Community**        | ✅ Large   | ✅ Large      | ⚠️ Medium     |

**Winner: sqlx** - Best balance of simplicity, performance, and flexibility for financial application.

---

## Repository Pattern

### Why Repository Pattern?

1. **Testability**: Mock repositories for unit tests
2. **Separation of Concerns**: Business logic separate from SQL
3. **DRY**: Reuse common queries across handlers
4. **Consistency**: Standardized error handling
5. **Maintainability**: Single place to update queries

### Repository Interface

```go
// repositories/user_repository.go
package repositories

import (
    "context"
    "database/sql"
)

type User struct {
    ID           int       `db:"id"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
}

type UserRepository interface {
    Create(ctx context.Context, email, passwordHash string) (*User, error)
    GetByID(ctx context.Context, id int) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, email, passwordHash string) (*User, error) {
    var user User
    query := `
        INSERT INTO users (email, password_hash, created_at, updated_at)
        VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, email, password_hash, created_at, updated_at
    `
    err := r.db.GetContext(ctx, &user, query, email, passwordHash)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    query := `SELECT * FROM users WHERE email = $1`
    err := r.db.GetContext(ctx, &user, query, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found
        }
        return nil, err
    }
    return &user, nil
}
```

### Budget Repository Example

```go
// repositories/budget_repository.go
package repositories

import (
    "context"
    "database/sql"
)

type MonthlyBudget struct {
    ID        int       `db:"id"`
    UserID    int       `db:"user_id"`
    Month     string    `db:"month"` // YYYY-MM
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type BudgetRepository interface {
    Create(ctx context.Context, userID int, month, name string) (*MonthlyBudget, error)
    GetByUserAndMonth(ctx context.Context, userID int, month string) (*MonthlyBudget, error)
    ListByUser(ctx context.Context, userID int) ([]MonthlyBudget, error)
    Update(ctx context.Context, budget *MonthlyBudget) error
    Delete(ctx context.Context, id int) error
}

type budgetRepository struct {
    db *sqlx.DB
}

func NewBudgetRepository(db *sqlx.DB) BudgetRepository {
    return &budgetRepository{db: db}
}

func (r *budgetRepository) GetByUserAndMonth(ctx context.Context, userID int, month string) (*MonthlyBudget, error) {
    var budget MonthlyBudget
    query := `
        SELECT id, user_id, month, name, created_at, updated_at
        FROM monthly_budgets
        WHERE user_id = $1 AND month = $2
    `
    err := r.db.GetContext(ctx, &budget, query, userID, month)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &budget, nil
}

func (r *budgetRepository) ListByUser(ctx context.Context, userID int) ([]MonthlyBudget, error) {
    var budgets []MonthlyBudget
    query := `
        SELECT id, user_id, month, name, created_at, updated_at
        FROM monthly_budgets
        WHERE user_id = $1
        ORDER BY month DESC
    `
    err := r.db.SelectContext(ctx, &budgets, query, userID)
    if err != nil {
        return nil, err
    }
    return budgets, nil
}
```

---

## Connection Pooling

### pgx/v5 Connection Pool Configuration

```go
// database/postgres.go
package database

import (
    "context"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string

    // Connection pool settings
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Connection pool configuration
    db.SetMaxOpenConns(cfg.MaxOpenConns)       // Max number of open connections
    db.SetMaxIdleConns(cfg.MaxIdleConns)       // Max number of idle connections
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime) // Max lifetime of a connection
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // Max idle time before closing

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}

// Recommended settings for budget app
func DefaultConfig() Config {
    return Config{
        Host:            "localhost",
        Port:            5432,
        User:            "postgres",
        Password:        "password",
        DBName:          "budget_tracker",
        SSLMode:         "disable", // Use "require" in production

        MaxOpenConns:    25,                   // Max 25 concurrent connections
        MaxIdleConns:    5,                    // Keep 5 idle connections ready
        ConnMaxLifetime: 30 * time.Minute,     // Recycle connections after 30min
        ConnMaxIdleTime: 5 * time.Minute,      // Close idle conns after 5min
    }
}
```

**Pool Settings Rationale**:
- `MaxOpenConns=25`: Budget app is low-traffic, 25 connections sufficient
- `MaxIdleConns=5`: Keep 5 ready for quick response
- `ConnMaxLifetime=30min`: Recycle to handle DB restarts gracefully
- `ConnMaxIdleTime=5min`: Release idle connections to save resources

---

## Transaction Management

### Manual Transactions for Complex Operations

```go
// repositories/transaction.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/jmoiron/sqlx"
)

// WithTransaction executes fn within a database transaction.
// If fn returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func WithTransaction(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            // Rollback on panic
            tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("rollback error: %w (original error: %v)", rbErr, err)
        }
        return err
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}

// Example: Create budget with default categories
func (r *budgetRepository) CreateWithCategories(ctx context.Context, userID int, month, name string) (*MonthlyBudget, error) {
    var budget *MonthlyBudget

    err := WithTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
        // Create budget
        var b MonthlyBudget
        budgetQuery := `
            INSERT INTO monthly_budgets (user_id, month, name, created_at, updated_at)
            VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
            RETURNING id, user_id, month, name, created_at, updated_at
        `
        if err := tx.GetContext(ctx, &b, budgetQuery, userID, month, name); err != nil {
            return fmt.Errorf("insert budget: %w", err)
        }
        budget = &b

        // Create default categories
        categoryQuery := `
            INSERT INTO categories (monthly_budget_id, name, projected, created_at, updated_at)
            VALUES 
                ($1, 'Income', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
                ($1, 'Savings (Pre-Tax)', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
                ($1, 'Savings (Post-Tax)', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
                ($1, 'Expenses', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        `
        if _, err := tx.ExecContext(ctx, categoryQuery, budget.ID); err != nil {
            return fmt.Errorf("insert categories: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return budget, nil
}
```

### Transaction Isolation Levels

```go
// For critical financial operations
func (r *repository) UpdateWithIsolation(ctx context.Context) error {
    // Serializable isolation for strict consistency
    tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly:  false,
    })
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // ... do work ...

    return tx.Commit()
}
```

---

## Query Patterns

### Named Parameters with Structs

```go
// Use struct tags for named parameters
type CreateCategoryParams struct {
    BudgetID  int     `db:"monthly_budget_id"`
    Name      string  `db:"name"`
    Projected float64 `db:"projected"`
    ParentID  *int    `db:"parent_id"` // Nullable
}

func (r *categoryRepository) Create(ctx context.Context, params CreateCategoryParams) (*Category, error) {
    var category Category
    query := `
        INSERT INTO categories (monthly_budget_id, name, projected, parent_id, created_at, updated_at)
        VALUES (:monthly_budget_id, :name, :projected, :parent_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, monthly_budget_id, name, projected, manual_actual, created_at, updated_at
    `
    stmt, err := r.db.PrepareNamedContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    err = stmt.GetContext(ctx, &category, params)
    if err != nil {
        return nil, err
    }
    return &category, nil
}
```

### Bulk Insert with Named Parameters

```go
// Insert multiple transactions from CSV
func (r *transactionRepository) BulkInsert(ctx context.Context, csvFileID int, transactions []Transaction) error {
    query := `
        INSERT INTO transactions (csv_file_id, row_number, transaction_date, post_date, description, category, type, amount, note, created_at, updated_at)
        VALUES (:csv_file_id, :row_number, :transaction_date, :post_date, :description, :category, :type, :amount, :note, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `
    _, err := r.db.NamedExecContext(ctx, query, transactions)
    return err
}
```

### Complex Queries with JOINs

```go
// Get budget summary with category totals
func (r *budgetRepository) GetSummaryWithCategories(ctx context.Context, userID int, month string) (*BudgetSummary, error) {
    query := `
        SELECT 
            mb.id AS budget_id,
            mb.month,
            mb.name,
            c.id AS category_id,
            c.name AS category_name,
            c.projected,
            COALESCE(c.manual_actual, COALESCE(SUM(cr.amount_snapshot), 0)) AS actual,
            COUNT(cr.id) AS linked_count
        FROM monthly_budgets mb
        LEFT JOIN categories c ON mb.id = c.monthly_budget_id
        LEFT JOIN cell_references cr ON c.id = cr.category_id
        WHERE mb.user_id = $1 AND mb.month = $2
        GROUP BY mb.id, c.id
        ORDER BY c.name
    `
    
    var results []CategoryRow
    err := r.db.SelectContext(ctx, &results, query, userID, month)
    if err != nil {
        return nil, err
    }

    // Build summary from results
    summary := buildSummary(results)
    return summary, nil
}
```

---

## Error Handling

### Database Error Types

```go
// errors/database.go
package errors

import (
    "database/sql"
    "errors"
    "fmt"

    "github.com/lib/pq" // PostgreSQL driver errors
)

var (
    ErrNotFound      = errors.New("resource not found")
    ErrAlreadyExists = errors.New("resource already exists")
    ErrForeignKey    = errors.New("foreign key constraint violation")
    ErrCheckFailed   = errors.New("check constraint failed")
)

// MapPostgresError converts PostgreSQL errors to application errors
func MapPostgresError(err error) error {
    if err == nil {
        return nil
    }

    // Check for sql.ErrNoRows
    if errors.Is(err, sql.ErrNoRows) {
        return ErrNotFound
    }

    // Check for PostgreSQL-specific errors
    var pqErr *pq.Error
    if errors.As(err, &pqErr) {
        switch pqErr.Code {
        case "23505": // unique_violation
            return fmt.Errorf("%w: %s", ErrAlreadyExists, pqErr.Detail)
        case "23503": // foreign_key_violation
            return fmt.Errorf("%w: %s", ErrForeignKey, pqErr.Detail)
        case "23514": // check_violation
            return fmt.Errorf("%w: %s", ErrCheckFailed, pqErr.Detail)
        }
    }

    return err
}

// Example usage in repository
func (r *userRepository) Create(ctx context.Context, email, passwordHash string) (*User, error) {
    var user User
    query := `...`
    err := r.db.GetContext(ctx, &user, query, email, passwordHash)
    if err != nil {
        return nil, MapPostgresError(err)
    }
    return &user, nil
}
```

---

## Testing Strategy

### Mock Repository

```go
// repositories/mocks/user_repository_mock.go
package mocks

import (
    "context"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Usage in handler test
func TestRegisterUser(t *testing.T) {
    mockRepo := new(mocks.MockUserRepository)
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, "test@example.com", mock.AnythingOfType("string")).
        Return(&User{ID: 1, Email: "test@example.com"}, nil)

    service := NewUserService(mockRepo)
    user, err := service.Register(context.Background(), "test@example.com", "password123")

    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests with Test Database

```go
// repositories/user_repository_test.go
package repositories_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
    db, err := sqlx.Connect("postgres", "postgres://test:test@localhost/test_budget_tracker?sslmode=disable")
    require.NoError(t, err)

    // Run migrations
    // ... migration logic ...

    t.Cleanup(func() {
        db.Exec("TRUNCATE TABLE users CASCADE")
        db.Close()
    })

    return db
}

func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    user, err := repo.Create(context.Background(), "test@example.com", "hashedpassword")
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "test@example.com", user.Email)
}
```

---

## Migration Integration

### Run Migrations on Startup

```go
// database/migrate.go
package database

import (
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB, migrationsPath string) error {
    driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("create postgres driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", migrationsPath),
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("run migrations: %w", err)
    }

    return nil
}
```

---

## Application Startup

### Complete Database Setup

```go
// main.go
package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    
    "your-app/database"
    "your-app/repositories"
    "your-app/handlers"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Database configuration
    dbConfig := database.Config{
        Host:            os.Getenv("DB_HOST"),
        Port:            5432,
        User:            os.Getenv("DB_USER"),
        Password:        os.Getenv("DB_PASSWORD"),
        DBName:          os.Getenv("DB_NAME"),
        SSLMode:         "disable",
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 30 * time.Minute,
        ConnMaxIdleTime: 5 * time.Minute,
    }

    // Connect to database
    db, err := database.NewPostgresDB(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Run migrations
    if err := database.RunMigrations(db, "./db/migrations"); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    budgetRepo := repositories.NewBudgetRepository(db)
    categoryRepo := repositories.NewCategoryRepository(db)

    // Initialize handlers
    userHandler := handlers.NewUserHandler(userRepo)
    budgetHandler := handlers.NewBudgetHandler(budgetRepo, categoryRepo)

    // Setup Gin router
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // Routes
    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", userHandler.Register)
            auth.POST("/login", userHandler.Login)
        }

        budgets := api.Group("/budgets")
        budgets.Use(AuthMiddleware())
        {
            budgets.GET("", budgetHandler.List)
            budgets.POST("", budgetHandler.Create)
            budgets.GET("/:month", budgetHandler.Get)
        }
    }

    // Start server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

---

## Pros of sqlx + Repository Pattern

- ✅ **Explicit SQL**: Full control over queries, easier debugging
- ✅ **Performance**: Near-zero overhead, optimal for financial calculations
- ✅ **Type Safety**: StructScan provides compile-time safety
- ✅ **Testability**: Easy to mock repositories
- ✅ **Flexibility**: Can write any SQL query needed
- ✅ **Simplicity**: No code generation, no ORM magic
- ✅ **Learning Curve**: Familiar SQL for developers
- ✅ **Context Support**: Built-in context.Context support
- ✅ **Transaction Control**: Manual transaction management when needed

---

## Cons/Challenges

- ⚠️ **Boilerplate**: Repository pattern requires interface definitions
- ⚠️ **Manual SQL**: Must write all queries (but also a pro for control)
- ⚠️ **No Associations**: Must handle foreign keys and joins manually
- ⚠️ **No Migrations**: Need golang-migrate (already decided in doc 11)
- ⚠️ **Error Handling**: Must map database errors to application errors

---

## Next Steps

1. **Create database package** (database/postgres.go) with connection pooling
2. **Define repository interfaces** (repositories/interfaces.go)
3. **Implement repositories** (repositories/*_repository.go)
4. **Add error mapping** (errors/database.go)
5. **Write repository tests** (repositories/*_test.go)
6. **Integrate with Gin handlers** (handlers/*.go using repositories)
7. **Add transaction helper** (database/transaction.go)

---

## References

- [sqlx GitHub Repository](https://github.com/jmoiron/sqlx)
- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [GORM Documentation](https://gorm.io/)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [PostgreSQL Connection Pooling](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [Go database/sql Package](https://pkg.go.dev/database/sql)

---

**Next Research**: Authentication (14-auth.md) - JWT tokens, password hashing, login/register flows
