# PostgreSQL Migrations with golang-migrate Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: golang-migrate with Sequential Versioning

---

## Recommendation Summary

**Chosen Tool: golang-migrate**  
**Chosen Versioning: Sequential integers (000001, 000002, 000003...)**

Based on comprehensive research of the golang-migrate repository, this is the recommended migration strategy for the budget tracking application.

**Why golang-migrate**:
1. ✅ **Standard Go migration tool** - Most popular migration library (15k+ stars)
2. ✅ **PostgreSQL driver support** - Multiple drivers (pgx/v4, pgx/v5, postgres)
3. ✅ **CLI and library modes** - Can run from CLI or embed in Go application
4. ✅ **Atomic migrations** - Transactions wrap migrations for safety
5. ✅ **Dirty state detection** - Prevents running migrations on failed state
6. ✅ **Locking mechanism** - Advisory locks prevent concurrent migrations
7. ✅ **Separate up/down files** - Clear reversibility, no custom syntax
8. ✅ **Active maintenance** - Regular updates and bug fixes

**Why Sequential Versioning over Timestamps**:
- Easier to track migration order in development
- Simpler to coordinate between team members
- Less merge conflict prone
- Still supports 2.1 billion migrations (plenty!)

---

## Core Requirements

From research and spike document:
- Create and track database schema changes
- Support up/down migrations for reversibility
- Prevent concurrent migrations from multiple instances
- Handle failed migrations gracefully (dirty state)
- Support both CLI usage and programmatic Go usage
- Integrate with PostgreSQL schema from research docs 09-* and 10-*

---

## Migration File Structure

### Filename Format

```
{version}_{title}.{direction}.{extension}

Examples:
000001_create_initial_schema.up.sql
000001_create_initial_schema.down.sql
000002_seed_default_categories.up.sql
000002_seed_default_categories.down.sql
000003_add_manual_actual_column.up.sql
000003_add_manual_actual_column.down.sql
```

**Components**:
- **version**: Sequential integer (000001, 000002, ...) with leading zeros for sorting
- **title**: Descriptive name (unused by tool, for human readability)
- **direction**: `up` (apply migration) or `down` (revert migration)
- **extension**: `.sql` for SQL migrations

**Key Points from golang-migrate FAQ**:
- Version must be distinct 64-bit unsigned integer
- Versions applied in ascending order (up) or descending order (down)
- Title is for readability only (can use underscores or hyphens)
- Extension not validated by library (use `.sql` for SQL databases)

---

## Migration Directory Structure

```
project/
├── backend/
│   ├── main.go
│   ├── db/
│   │   └── migrations/
│   │       ├── 000001_create_initial_schema.up.sql
│   │       ├── 000001_create_initial_schema.down.sql
│   │       ├── 000002_seed_default_categories.up.sql
│   │       ├── 000002_seed_default_categories.down.sql
│   │       └── 000003_add_manual_actual_column.up.sql
│   │       └── 000003_add_manual_actual_column.down.sql
└── frontend/
```

**Recommended**: Store migrations in `backend/db/migrations/` directory

---

## Creating Migrations

### Using CLI

```bash
# Install golang-migrate CLI
brew install golang-migrate

# Or via Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration
migrate create -ext sql -dir db/migrations -seq create_initial_schema

# Output:
# db/migrations/000001_create_initial_schema.up.sql
# db/migrations/000001_create_initial_schema.down.sql
```

**Flags**:
- `-ext sql`: File extension (for SQL migrations)
- `-dir db/migrations`: Directory path
- `-seq`: Sequential versioning (000001, 000002, ...)
- `-timestamp`: Alternative - timestamp versioning (1500360784)

### Programmatically in Go

```go
// Not needed - use CLI for creating migration files
// Only use Go code for running migrations
```

---

## Migration Content Best Practices

### Wrapping in Transactions

**From golang-migrate docs**:
> In PostgreSQL running multiple SQL statements in one `Exec` executes them inside a transaction.

```sql
-- 000001_create_initial_schema.up.sql
BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

COMMIT;
```

**Benefits**:
- Atomic migration: all-or-nothing execution
- Failed statement rolls back entire migration
- Database left in consistent state

**Exceptions**: Some PostgreSQL operations cannot run in transactions:
- `CREATE INDEX CONCURRENTLY` - must be in separate migration without transaction
- `ALTER TYPE ... ADD VALUE` (enums) - not transactional in PostgreSQL < 12

### Idempotent Migrations

```sql
-- Good: Idempotent migration (can run multiple times safely)
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

COMMIT;
```

**Benefits**:
- Safe to re-run if migration partially failed
- Prevents errors if schema already exists
- Easier development workflow

### Reversible Down Migrations

```sql
-- 000001_create_initial_schema.down.sql
BEGIN;

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

COMMIT;
```

**Best Practices**:
- Always provide down migration (reversibility)
- Use `IF EXISTS` for safety
- Drop in reverse order of creation (indexes before tables)
- Drop foreign key constraints before referenced tables

### Empty Migrations (No-op)

From golang-migrate FAQ:
> It is recommended to still include both migration files by making the whole migration file consist of a comment.

```sql
-- 000005_no_op_migration.up.sql
-- This migration intentionally left empty
-- Data backfill handled in separate script
```

**When to Use**:
- Migration is no-op (placeholder for version)
- Migration is irreversible (e.g., data deletion)
- Separating DDL and data migrations

---

## Example Migrations for Budget App

### 000001_create_initial_schema.up.sql

```sql
BEGIN;

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Monthly budgets table
CREATE TABLE monthly_budgets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    month VARCHAR(7) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, month),
    CHECK (month ~ '^\d{4}-\d{2}$')
);

CREATE INDEX idx_monthly_budgets_user_month ON monthly_budgets(user_id, month);

-- Categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    projected NUMERIC(12,2) NOT NULL DEFAULT 0,
    manual_actual NUMERIC(12,2) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE
);

CREATE INDEX idx_categories_budget ON categories(monthly_budget_id);

-- Category paths (closure table)
CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);
CREATE INDEX idx_category_paths_depth ON category_paths(depth);

-- CSV files table
CREATE TABLE csv_files (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_size INTEGER,
    row_count INTEGER NOT NULL DEFAULT 0,
    
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE,
    CHECK (row_count >= 0)
);

CREATE INDEX idx_csv_files_budget ON csv_files(monthly_budget_id);

-- Transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    csv_file_id INTEGER NOT NULL,
    row_number INTEGER NOT NULL,
    
    transaction_date DATE,
    post_date DATE,
    description TEXT NOT NULL,
    category VARCHAR(255),
    type VARCHAR(50),
    amount NUMERIC(12,2) NOT NULL,
    note TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (csv_file_id) REFERENCES csv_files(id) ON DELETE CASCADE,
    UNIQUE(csv_file_id, row_number)
);

CREATE INDEX idx_transactions_csv_file ON transactions(csv_file_id);
CREATE INDEX idx_transactions_amount ON transactions(amount);

-- Cell references table
CREATE TABLE cell_references (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    transaction_id INTEGER NOT NULL,
    amount_snapshot NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    UNIQUE(category_id, transaction_id)
);

CREATE INDEX idx_cell_references_category ON cell_references(category_id);
CREATE INDEX idx_cell_references_transaction ON cell_references(transaction_id);

COMMIT;
```

### 000001_create_initial_schema.down.sql

```sql
BEGIN;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS cell_references CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS csv_files CASCADE;
DROP TABLE IF EXISTS category_paths CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS monthly_budgets CASCADE;
DROP TABLE IF EXISTS users CASCADE;

COMMIT;
```

### 000002_seed_default_categories.up.sql

```sql
-- This migration seeds default categories for new budgets
-- Note: This is intentionally NOT wrapped in transaction
-- because it will be run via application code, not as DDL

-- Application code will handle seeding default categories:
-- - Income
-- - Savings (Pre-Tax)
-- - Savings (Post-Tax)
-- - Expenses

-- This file serves as documentation of data migration intent
```

### 000002_seed_default_categories.down.sql

```sql
-- Default categories are created per-budget by application
-- No down migration needed (deleting budget cascades to categories)
```

---

## Running Migrations

### Via CLI (Development)

```bash
# Set database URL
export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/budget_app?sslmode=disable'

# Run all pending migrations
migrate -database ${POSTGRESQL_URL} -path db/migrations up

# Rollback one migration
migrate -database ${POSTGRESQL_URL} -path db/migrations down 1

# Migrate to specific version
migrate -database ${POSTGRESQL_URL} -path db/migrations goto 2

# Check current version
migrate -database ${POSTGRESQL_URL} -path db/migrations version

# Force version (after fixing dirty state)
migrate -database ${POSTGRESQL_URL} -path db/migrations force 1
```

### Via Go Code (Production)

```go
package main

import (
    "database/sql"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    // Open database connection
    db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/budget_app?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create migrate instance
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Run all pending migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migrations applied successfully")
}
```

**Best Practice**: Run migrations at application startup in production

```go
func initDatabase() error {
    // ... create migrate instance

    // Apply migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}

func main() {
    if err := initDatabase(); err != nil {
        log.Fatalf("Database initialization failed: %v", err)
    }

    // Start HTTP server
    startServer()
}
```

---

## Migration Versioning Table

golang-migrate automatically creates `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);
```

**Columns**:
- **version**: Current migration version number
- **dirty**: TRUE if migration failed, FALSE if successful

**Query Current Version**:
```sql
SELECT version, dirty FROM schema_migrations;
```

**Example Output**:
```
 version | dirty
---------+-------
       3 | false
```

---

## Dirty State Handling

### What is Dirty State?

From golang-migrate FAQ:
> Before a migration runs, each database sets a dirty flag. Execution stops if a migration fails and the dirty state persists, which prevents attempts to run more migrations on top of a failed migration.

**Example Scenario**:
1. Migration 000003 starts (version=3, dirty=true)
2. Migration fails midway (syntax error in SQL)
3. Database left in dirty state (version=3, dirty=true)
4. Future migration attempts blocked

### Fixing Dirty State

**Step 1: Diagnose the problem**
```bash
# Check current version
migrate -database ${POSTGRESQL_URL} -path db/migrations version
# Output: 3/dirty
```

**Step 2: Manually fix database**
```sql
-- Connect to database and fix the issue
-- Example: Complete the failed migration manually
CREATE TABLE foo (id SERIAL PRIMARY KEY);
```

**Step 3: Force clean version**
```bash
# Force version to 3 (mark as clean)
migrate -database ${POSTGRESQL_URL} -path db/migrations force 3
```

**Step 4: Continue migrations**
```bash
# Now can run pending migrations
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

**Prevention**: Always wrap migrations in transactions (BEGIN/COMMIT)

---

## Locking Mechanism

From golang-migrate FAQ:
> Database-specific locking features are used by some database drivers to prevent multiple instances of migrate from running migrations on the same database at the same time.

### PostgreSQL Advisory Locks

golang-migrate uses `pg_advisory_lock` function:

```sql
-- Migration acquires lock
SELECT pg_advisory_lock(hashtext('schema_migrations'));

-- Run migration SQL statements

-- Release lock
SELECT pg_advisory_unlock(hashtext('schema_migrations'));
```

**Benefits**:
- Prevents concurrent migrations from multiple app instances
- Automatic lock release on connection close
- No manual cleanup needed

**Configuration** (postgres driver):
```
postgres://user:pass@host:port/db?x-lock-strategy=advisory
```

---

## Multi-Statement Mode

From golang-migrate docs:
> In PostgreSQL running multiple SQL statements in one `Exec` executes them inside a transaction. Sometimes this behavior is not desirable because some statements can be only run outside of transaction (e.g. `CREATE INDEX CONCURRENTLY`).

### Enabling Multi-Statement Mode

**Via URL Parameter**:
```
postgres://user:pass@host:port/db?x-multi-statement=true
```

**Via Go Config**:
```go
driver, err := postgres.WithInstance(db, &postgres.Config{
    MultiStatementEnabled: true,
    MultiStatementMaxSize: 10 * 1024 * 1024, // 10MB
})
```

**Use Case**: Running `CREATE INDEX CONCURRENTLY` (cannot be in transaction)

```sql
-- 000004_add_concurrent_index.up.sql
-- NO BEGIN/COMMIT - runs outside transaction
CREATE INDEX CONCURRENTLY idx_transactions_description ON transactions(description);
```

---

## Best Practices Summary

### 1. Always Provide Down Migrations
✅ **Do**:
```sql
-- 000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

❌ **Don't**: Skip down migrations (makes rollback impossible)

---

### 2. Wrap Migrations in Transactions

✅ **Do**:
```sql
BEGIN;
CREATE TABLE foo (...);
COMMIT;
```

❌ **Don't**: Run bare SQL (leaves database in inconsistent state on error)

**Exception**: Operations that cannot run in transactions (CREATE INDEX CONCURRENTLY)

---

### 3. Use Idempotent Operations

✅ **Do**:
```sql
CREATE TABLE IF NOT EXISTS users (...);
DROP TABLE IF EXISTS users;
```

❌ **Don't**: Assume clean slate (migration may partially succeed)

---

### 4. Sequential Version Numbers

✅ **Do**: 000001, 000002, 000003 (with leading zeros)

❌ **Don't**: Timestamps (harder to track, merge conflicts in teams)

---

### 5. Separate DDL and Data Migrations

✅ **Do**:
```
000001_create_schema.up.sql     (DDL: CREATE TABLE)
000002_seed_categories.up.sql   (Data: INSERT INTO)
```

❌ **Don't**: Mix schema changes and data changes in one migration

---

### 6. Test Rollback

✅ **Do**:
```bash
migrate up    # Apply migration
migrate down  # Test rollback
migrate up    # Re-apply
```

❌ **Don't**: Deploy without testing down migration

---

### 7. Keep Migrations Small

✅ **Do**: One logical change per migration

❌ **Don't**: Combine unrelated changes (harder to rollback specific change)

---

## Production Deployment Workflow

### Pre-Deployment Checklist

- [ ] All migrations have up and down files
- [ ] Migrations wrapped in transactions (or documented why not)
- [ ] Down migrations tested locally
- [ ] Database backup created
- [ ] Migration version noted (for rollback)

### Deployment Steps

**1. Create Database Backup**
```bash
pg_dump -h localhost -U postgres budget_app > backup_$(date +%Y%m%d_%H%M%S).sql
```

**2. Check Current Version**
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations version
# Output: 3
```

**3. Apply Migrations**
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

**4. Verify Success**
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations version
# Output: 5 (if 2 new migrations applied)

# Check dirty flag is false
psql -h localhost -U postgres -d budget_app -c "SELECT version, dirty FROM schema_migrations;"
```

**5. If Migration Fails**

```bash
# Check error message
migrate -database ${POSTGRESQL_URL} -path db/migrations version
# Output: 4/dirty

# Option A: Fix and force version
# 1. Manually fix database issue
# 2. Force clean version
migrate -database ${POSTGRESQL_URL} -path db/migrations force 4

# Option B: Rollback to previous version
# 1. Restore from backup
pg_restore -h localhost -U postgres -d budget_app backup_20251019_143000.sql
# 2. Force previous version
migrate -database ${POSTGRESQL_URL} -path db/migrations force 3
```

---

## Pros of golang-migrate

- ✅ **Standard tool** - Most popular Go migration library
- ✅ **Separate up/down files** - Clear, no custom syntax
- ✅ **Transaction support** - Atomic migrations
- ✅ **Dirty state detection** - Prevents cascading failures
- ✅ **Locking** - Prevents concurrent migrations
- ✅ **CLI and library** - Flexible usage (development vs production)
- ✅ **Multiple database support** - PostgreSQL, MySQL, SQLite, etc.
- ✅ **Embedded migrations** - Can embed in Go binary using go:embed

---

## Cons/Challenges

- ⚠️ **Manual dirty state fixes** - Requires understanding of failed migration
- ⚠️ **No automatic rollback** - Must manually run down migration
- ⚠️ **Version conflicts** - Team coordination needed for sequential versions
- ⚠️ **No migration dependencies** - Must manually order migrations
- ⚠️ **Limited rollback testing** - Down migrations not automatically tested in CI

---

## Alternative Tools (Not Selected)

### goose
- Similar features to golang-migrate
- Supports Go migrations (not just SQL)
- Less popular, smaller community

### sql-migrate
- Embeddable migrations
- gorp ORM integration
- Less active maintenance

### atlas
- Declarative schema approach (define desired state)
- Modern tool, but different paradigm
- More complex learning curve

**Recommendation**: Stick with golang-migrate (industry standard)

---

## Next Steps

1. **Create initial migration** (000001_create_initial_schema.up/down.sql)
2. **Test migration locally** (up, down, up again)
3. **Add to version control** (Git commit migrations)
4. **Document migration workflow** (README for team)
5. **Set up CI/CD** (run migrations in deployment pipeline)
6. **Create backup strategy** (automate pre-migration backups)

---

## References

- [golang-migrate GitHub](https://github.com/golang-migrate/migrate)
- [golang-migrate CLI Documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [PostgreSQL Driver README](https://github.com/golang-migrate/migrate/tree/master/database/postgres)
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
- [FAQ - Dirty State](https://github.com/golang-migrate/migrate/blob/master/FAQ.md)
- [PostgreSQL Tutorial](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)

---

**Next Research**: API endpoint design (12-api-design.md), Database integration (ORM evaluation)
