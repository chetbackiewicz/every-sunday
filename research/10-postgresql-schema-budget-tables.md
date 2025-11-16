# PostgreSQL Schema - Budget and Transaction Tables Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: Normalized Schema with Foreign Keys and Indexes

---

## Recommendation Summary

**Chosen Approach: Normalized Relational Schema with Financial Best Practices**

Based on PostgreSQL documentation and best practices research, the budget and transaction tables should follow a normalized design with proper constraints, indexes, and financial data types.

**Key Decisions**:
1. ✅ **NUMERIC(12,2) for amounts** - Exact decimal storage for financial data (NO money type, NO floats)
2. ✅ **TIMESTAMPTZ for dates** - Timezone-aware timestamps for accurate temporal data
3. ✅ **Foreign keys with ON DELETE CASCADE** - Referential integrity and automatic cleanup
4. ✅ **Composite indexes** - Optimized queries for user_id + month patterns
5. ✅ **Generated columns** - Auto-calculated difference (projected - actual)
6. ✅ **CHECK constraints** - Data validation at database level
7. ✅ **NOT NULL constraints** - Prevent invalid empty states

**Why This Approach**:
- Prevents rounding errors in financial calculations
- Maintains referential integrity across tables
- Optimizes common query patterns (user budgets by month)
- Validates data at database level (not just application)
- Standard SQL portability (no PostgreSQL-specific types except timestamptz)

---

## Core Requirements

From spike document and existing research:
- **Monthly budgets** - User can create budgets for different months
- **Category hierarchy** - Integration with closure table from 09-postgresql-schema-hierarchical.md
- **CSV file tracking** - Up to 10 CSV files per month
- **Transaction storage** - All CSV rows stored with metadata
- **Cell references** - Link transactions to category "Actual" values
- **Manual overrides** - User can manually edit category actual values
- **User isolation** - Multi-user support with proper data separation

---

## Complete Schema Design

### Table 1: users (Multi-user support)

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  -- bcrypt hash
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for login queries
CREATE INDEX idx_users_email ON users(email);
```

**Key Points**:
- User authentication table
- Email as unique identifier
- Timestamps for audit trail

---

### Table 2: monthly_budgets (Month container)

```sql
CREATE TABLE monthly_budgets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    month VARCHAR(7) NOT NULL,  -- Format: "YYYY-MM"
    name VARCHAR(255),           -- Optional: e.g., "October Budget - House Purchase"
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, month),  -- One budget per user per month
    CHECK (month ~ '^\d{4}-\d{2}$')  -- Validate YYYY-MM format
);

-- Composite index for user's budgets
CREATE INDEX idx_monthly_budgets_user_month ON monthly_budgets(user_id, month);
```

**Key Points**:
- One record per user per month
- VARCHAR(7) for month string (YYYY-MM format)
- CHECK constraint validates format
- Composite unique index prevents duplicates

---

### Table 3: categories (Updated from 09-*)

**From 09-postgresql-schema-hierarchical.md**, update to include monthly_budget_id:

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    projected NUMERIC(12,2) NOT NULL DEFAULT 0,
    manual_actual NUMERIC(12,2) DEFAULT NULL,  -- Override actual calculation
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE
);

-- Index for budget's categories
CREATE INDEX idx_categories_budget ON categories(monthly_budget_id);
```

**Changes from 09-***:
- Replaced user_id + month with monthly_budget_id FK
- Added manual_actual for overrides
- Removed user_id and month columns (accessed via FK)

---

### Table 4: category_paths (Unchanged from 09-*)

```sql
CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Indexes from 09-*
CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);
CREATE INDEX idx_category_paths_depth ON category_paths(depth);
CREATE INDEX idx_category_paths_ancestor_depth ON category_paths(ancestor_id, depth);
```

---

### Table 5: csv_files (File upload tracking)

```sql
CREATE TABLE csv_files (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_size INTEGER,  -- Bytes
    row_count INTEGER NOT NULL DEFAULT 0,
    
    -- Constraints
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE,
    CHECK (row_count >= 0)
);

-- Index for budget's CSV files
CREATE INDEX idx_csv_files_budget ON csv_files(monthly_budget_id);
```

**Key Points**:
- One record per uploaded CSV file
- Tracks original filename and metadata
- row_count for quick stats
- Up to 10 files per budget (enforced in application logic)

---

### Table 6: transactions (CSV row storage)

```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    csv_file_id INTEGER NOT NULL,
    row_number INTEGER NOT NULL,  -- Original row position in CSV
    
    -- CSV columns
    transaction_date DATE,           -- "Transaction Date" from CSV
    post_date DATE,                  -- "Post Date" or "Posting Date" from CSV
    description TEXT NOT NULL,
    category VARCHAR(255),           -- User-editable category name
    type VARCHAR(50),                -- e.g., "Sale", "Payment", "Credit"
    amount NUMERIC(12,2) NOT NULL,  -- Financial amount (exact decimal)
    note TEXT,                       -- User-editable note
    
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    FOREIGN KEY (csv_file_id) REFERENCES csv_files(id) ON DELETE CASCADE,
    UNIQUE(csv_file_id, row_number)  -- Prevent duplicate rows from same CSV
);

-- Indexes for queries
CREATE INDEX idx_transactions_csv_file ON transactions(csv_file_id);
CREATE INDEX idx_transactions_amount ON transactions(amount);
CREATE INDEX idx_transactions_dates ON transactions(transaction_date, post_date);
```

**Key Points**:
- **NUMERIC(12,2)** for amount (NOT money, NOT float) - exact decimal storage
- **TEXT** for description/note - no arbitrary length limits
- row_number preserves original CSV order
- Nullable dates (some CSVs may have missing dates)
- Composite unique constraint prevents duplicate imports

**Why NUMERIC(12,2)**:
- 12 total digits: supports up to $9,999,999,999.99 (10 billion)
- 2 decimal places: cents precision
- Exact storage: no floating-point rounding errors
- Standard across databases: portable

---

### Table 7: cell_references (Transaction → Category links)

```sql
CREATE TABLE cell_references (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    transaction_id INTEGER NOT NULL,
    amount_snapshot NUMERIC(12,2) NOT NULL,  -- Cached amount at time of link
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    UNIQUE(category_id, transaction_id)  -- Prevent duplicate links
);

-- Indexes for queries
CREATE INDEX idx_cell_references_category ON cell_references(category_id);
CREATE INDEX idx_cell_references_transaction ON cell_references(transaction_id);
```

**Key Points**:
- Tracks which transactions contribute to category "Actual" values
- amount_snapshot caches transaction amount (denormalization for performance)
- ON DELETE CASCADE: deleting transaction removes cell reference
- Unique constraint: one transaction can only link to one category once

**Why amount_snapshot**:
- Fast calculation: no need to JOIN transactions table
- Historical accuracy: preserves amount even if transaction edited
- Tradeoff: requires update if transaction amount changes

---

## Complete ER Diagram (ASCII)

```
users
  ├─── monthly_budgets (user_id FK)
       ├─── categories (monthly_budget_id FK)
       │    ├─── category_paths (ancestor_id, descendant_id FKs)
       │    └─── cell_references (category_id FK)
       │
       └─── csv_files (monthly_budget_id FK)
            └─── transactions (csv_file_id FK)
                 └─── cell_references (transaction_id FK)
```

**Relationships**:
- User → Monthly Budgets (1:N)
- Monthly Budget → Categories (1:N)
- Monthly Budget → CSV Files (1:N)
- Category → Category Paths (N:N via closure table)
- Category → Cell References (1:N)
- CSV File → Transactions (1:N)
- Transaction → Cell References (1:N)

---

## Common Queries

### 1. Get User's Budget for Specific Month

```sql
SELECT mb.*, u.email
FROM monthly_budgets mb
JOIN users u ON mb.user_id = u.id
WHERE u.id = $1 AND mb.month = $2;
```

### 2. Get All Categories for a Budget (with totals)

```sql
SELECT 
    c.id,
    c.name,
    c.projected,
    c.manual_actual,
    -- Calculate actual from linked cells
    COALESCE(c.manual_actual, COALESCE(SUM(cr.amount_snapshot), 0)) AS actual,
    c.projected - COALESCE(c.manual_actual, COALESCE(SUM(cr.amount_snapshot), 0)) AS difference
FROM categories c
LEFT JOIN cell_references cr ON c.id = cr.category_id
WHERE c.monthly_budget_id = $1
GROUP BY c.id
ORDER BY c.name;
```

**Logic**:
- If `manual_actual IS NOT NULL`: use manual override
- Else: sum all `cell_references.amount_snapshot`
- Calculate difference: `projected - actual`

### 3. Get All Transactions for a Budget

```sql
SELECT 
    t.*,
    cf.filename,
    EXISTS(SELECT 1 FROM cell_references cr WHERE cr.transaction_id = t.id) AS is_linked
FROM transactions t
JOIN csv_files cf ON t.csv_file_id = cf.id
WHERE cf.monthly_budget_id = $1
ORDER BY t.transaction_date DESC, t.row_number;
```

### 4. Get Cell References for a Category

```sql
SELECT 
    cr.id,
    cr.amount_snapshot,
    t.transaction_date,
    t.description,
    t.amount AS current_amount,  -- Compare with snapshot
    cf.filename
FROM cell_references cr
JOIN transactions t ON cr.transaction_id = t.id
JOIN csv_files cf ON t.csv_file_id = cf.id
WHERE cr.category_id = $1
ORDER BY t.transaction_date DESC;
```

### 5. Get Categories with Hierarchy (using closure table)

```sql
-- Get category with all descendants
SELECT 
    c.*,
    cp.depth
FROM categories c
JOIN category_paths cp ON c.id = cp.descendant_id
WHERE cp.ancestor_id = $1  -- Parent category ID
ORDER BY cp.depth, c.name;
```

### 6. Calculate Budget Summary (Income, Expenses, Remaining)

```sql
WITH category_totals AS (
    SELECT 
        c.id,
        c.name,
        c.monthly_budget_id,
        COALESCE(c.manual_actual, COALESCE(SUM(cr.amount_snapshot), 0)) AS actual
    FROM categories c
    LEFT JOIN cell_references cr ON c.id = cr.category_id
    WHERE c.monthly_budget_id = $1
    GROUP BY c.id
)
SELECT 
    SUM(CASE WHEN name = 'Income' THEN actual ELSE 0 END) AS total_income,
    SUM(CASE WHEN name = 'Savings (Pre-Tax)' THEN actual ELSE 0 END) AS savings_pretax,
    SUM(CASE WHEN name = 'Savings (Post-Tax)' THEN actual ELSE 0 END) AS savings_posttax,
    SUM(CASE WHEN name = 'Expenses' THEN actual ELSE 0 END) AS total_expenses,
    SUM(CASE WHEN name = 'Income' THEN actual ELSE 0 END) 
        - SUM(CASE WHEN name LIKE 'Savings%' THEN actual ELSE 0 END)
        - SUM(CASE WHEN name = 'Expenses' THEN actual ELSE 0 END) AS remaining
FROM category_totals;
```

---

## CRUD Operations

### Create New Budget

```sql
-- Step 1: Create monthly budget
INSERT INTO monthly_budgets (user_id, month, name)
VALUES ($1, $2, $3)
RETURNING id;

-- Step 2: Create default categories (Income, Savings, Expenses)
INSERT INTO categories (monthly_budget_id, name, projected)
VALUES 
    ($budget_id, 'Income', 0),
    ($budget_id, 'Savings (Pre-Tax)', 0),
    ($budget_id, 'Savings (Post-Tax)', 0),
    ($budget_id, 'Expenses', 0)
RETURNING id;

-- Step 3: Create category self-referencing paths
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
VALUES 
    ($income_id, $income_id, 0),
    ($savings_pretax_id, $savings_pretax_id, 0),
    ($savings_posttax_id, $savings_posttax_id, 0),
    ($expenses_id, $expenses_id, 0);
```

### Upload CSV File and Import Transactions

```sql
BEGIN;

-- Step 1: Create CSV file record
INSERT INTO csv_files (monthly_budget_id, filename, file_size, row_count)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- Step 2: Bulk insert transactions
INSERT INTO transactions (csv_file_id, row_number, transaction_date, post_date, description, category, type, amount, note)
VALUES 
    ($csv_file_id, 1, '2025-10-15', '2025-10-16', 'GITHUB INC PAYROLL', '', 'ACH_CREDIT', 2840.48, ''),
    ($csv_file_id, 2, '2025-10-14', '2025-10-14', 'Payment to Chase', '', 'LOAN_PMT', -325.68, ''),
    -- ... more rows
ON CONFLICT (csv_file_id, row_number) DO NOTHING;  -- Prevent duplicate imports

COMMIT;
```

### Link Transaction to Category

```sql
-- User clicks category "Actual" cell, then clicks transaction
INSERT INTO cell_references (category_id, transaction_id, amount_snapshot)
SELECT $category_id, $transaction_id, t.amount
FROM transactions t
WHERE t.id = $transaction_id
ON CONFLICT (category_id, transaction_id) DO NOTHING;  -- Already linked
```

### Delete Transaction (with automatic cell reference cleanup)

```sql
-- ON DELETE CASCADE automatically removes cell_references
DELETE FROM transactions WHERE id = $1;

-- If amount needs to be re-calculated, fetch updated totals:
SELECT 
    c.id,
    COALESCE(c.manual_actual, COALESCE(SUM(cr.amount_snapshot), 0)) AS new_actual
FROM categories c
LEFT JOIN cell_references cr ON c.id = cr.category_id
WHERE c.id = $category_id
GROUP BY c.id;
```

### Override Category Actual Value

```sql
-- User double-clicks "Actual" cell and types manual value
UPDATE categories
SET manual_actual = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $category_id;

-- To revert to calculated value, set manual_actual to NULL
UPDATE categories
SET manual_actual = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = $category_id;
```

### Update Transaction Amount (sync cell_references)

```sql
BEGIN;

-- Update transaction
UPDATE transactions
SET amount = $new_amount, updated_at = CURRENT_TIMESTAMP
WHERE id = $transaction_id;

-- Update all cell references that use this transaction
UPDATE cell_references
SET amount_snapshot = $new_amount
WHERE transaction_id = $transaction_id;

COMMIT;
```

---

## Data Type Decisions

### Financial Amounts: NUMERIC(12,2)

**✅ Use NUMERIC(12,2)**:
- Exact decimal storage (no rounding errors)
- Supports up to $9,999,999,999.99
- Standard SQL, portable across databases
- Arbitrary precision calculations

**❌ Avoid money type**:
- PostgreSQL-specific (not portable)
- Fixed to database's lc_monetary setting
- Changing lc_monetary corrupts all money values
- No fractional cent support

**❌ Avoid FLOAT/DOUBLE PRECISION**:
- Inexact representation (0.1 + 0.2 ≠ 0.3)
- Rounding errors accumulate in calculations
- Unsuitable for financial data

### Dates: DATE vs TIMESTAMPTZ

**Use DATE for**:
- transaction_date, post_date (no time component needed)
- month string VARCHAR(7) for simplicity

**Use TIMESTAMPTZ for**:
- created_at, updated_at, upload_date (audit trail)
- Timezone-aware, stores UTC internally
- Automatically adjusts to user's timezone on retrieval

**❌ Avoid TIMESTAMP (without time zone)**:
- No timezone information
- Ambiguous for international users
- Difficult date math across timezones

### Text: VARCHAR(N) vs TEXT

**Use VARCHAR(N) for**:
- Known-length fields: email (255), month (7), type (50)
- Database-level length validation
- Slightly better query planner estimates

**Use TEXT for**:
- Unknown-length fields: description, note
- No arbitrary limits
- Same storage efficiency as VARCHAR in PostgreSQL

---

## Constraints and Validation

### Foreign Key Constraints with CASCADE

```sql
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
```

**Benefits**:
- Referential integrity enforced by database
- Automatic cleanup: deleting user removes all budgets
- Prevents orphaned records
- No application logic needed

### CHECK Constraints

```sql
CHECK (month ~ '^\d{4}-\d{2}$')  -- Validate YYYY-MM format
CHECK (row_count >= 0)            -- Prevent negative counts
```

**Benefits**:
- Data validation at database level
- Prevents invalid inserts/updates
- Self-documenting schema

### UNIQUE Constraints

```sql
UNIQUE(user_id, month)                    -- One budget per user per month
UNIQUE(csv_file_id, row_number)          -- No duplicate CSV rows
UNIQUE(category_id, transaction_id)      -- No duplicate cell links
```

**Benefits**:
- Business rule enforcement
- Prevents data anomalies
- Composite unique indexes optimize queries

### NOT NULL Constraints

```sql
user_id INTEGER NOT NULL
month VARCHAR(7) NOT NULL
amount NUMERIC(12,2) NOT NULL
```

**Benefits**:
- Prevents invalid empty states
- Makes schema self-documenting
- Simplifies query logic (no NULL handling)

---

## Index Strategy

### Composite Indexes for Common Queries

```sql
CREATE INDEX idx_monthly_budgets_user_month ON monthly_budgets(user_id, month);
```

**Benefits**:
- Optimizes: "Get user's budget for October 2025"
- Covers both columns in single index
- PostgreSQL can use for user_id-only queries too

### Foreign Key Indexes

```sql
CREATE INDEX idx_categories_budget ON categories(monthly_budget_id);
CREATE INDEX idx_transactions_csv_file ON transactions(csv_file_id);
```

**Benefits**:
- Speeds up JOIN operations
- Required for efficient ON DELETE CASCADE
- Prevents full table scans

### Functional Indexes (Advanced)

```sql
-- Index for case-insensitive email lookup
CREATE INDEX idx_users_email_lower ON users(LOWER(email));

-- Query using index
SELECT * FROM users WHERE LOWER(email) = LOWER($1);
```

---

## Performance Considerations

### Denormalization: amount_snapshot in cell_references

**Tradeoff**:
- **Pro**: Fast calculation - no JOIN to transactions table
- **Pro**: Historical accuracy - preserves amount even if transaction edited
- **Con**: Requires UPDATE if transaction amount changes
- **Con**: Data duplication (amount stored twice)

**Recommendation**: Keep amount_snapshot for performance
- Budget app needs fast totals calculation
- Transaction edits are infrequent
- Application logic handles sync updates

### Materialized Views for Complex Aggregations (Optional)

```sql
CREATE MATERIALIZED VIEW budget_summaries AS
SELECT 
    mb.id AS budget_id,
    mb.user_id,
    mb.month,
    SUM(CASE WHEN c.name = 'Income' THEN COALESCE(c.manual_actual, cr_sums.total) ELSE 0 END) AS total_income,
    SUM(CASE WHEN c.name = 'Expenses' THEN COALESCE(c.manual_actual, cr_sums.total) ELSE 0 END) AS total_expenses
FROM monthly_budgets mb
LEFT JOIN categories c ON mb.id = c.monthly_budget_id
LEFT JOIN (
    SELECT category_id, SUM(amount_snapshot) AS total
    FROM cell_references
    GROUP BY category_id
) cr_sums ON c.id = cr_sums.category_id
GROUP BY mb.id;

-- Refresh after data changes
REFRESH MATERIALIZED VIEW budget_summaries;
```

**When to use**:
- Complex queries run frequently
- Data doesn't change constantly
- Dashboard/summary views

---

## Migration Strategy

### Migration File Example (using golang-migrate)

**000001_create_initial_schema.up.sql**:

```sql
BEGIN;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Create monthly_budgets table
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

-- ... continue with other tables

COMMIT;
```

**000001_create_initial_schema.down.sql**:

```sql
BEGIN;

DROP TABLE IF EXISTS cell_references CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS csv_files CASCADE;
DROP TABLE IF EXISTS category_paths CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS monthly_budgets CASCADE;
DROP TABLE IF EXISTS users CASCADE;

COMMIT;
```

### Migration Best Practices

1. **Wrap in transactions (BEGIN/COMMIT)** - Atomic migrations
2. **Idempotent operations** - Use IF EXISTS/IF NOT EXISTS
3. **Sequential versioning** - Timestamp or sequential numbers
4. **Reversible migrations** - Always provide .down.sql
5. **Test rollback** - Verify down migrations work
6. **Separate DDL and data** - Schema changes in one migration, data changes in another

---

## Pros of This Schema Design

- ✅ **Referential integrity** - Foreign keys prevent orphaned records
- ✅ **Data validation** - CHECK constraints enforce business rules
- ✅ **Performance** - Composite indexes optimize common queries
- ✅ **Financial accuracy** - NUMERIC type prevents rounding errors
- ✅ **Audit trail** - created_at/updated_at timestamps on all tables
- ✅ **Flexibility** - Manual overrides supported via manual_actual column
- ✅ **Cascading deletes** - Automatic cleanup when deleting users/budgets/transactions
- ✅ **Multi-user support** - User isolation via user_id + unique constraints
- ✅ **Scalability** - Indexes on all foreign keys and common query patterns
- ✅ **Standard SQL** - Portable to other PostgreSQL-compatible databases

---

## Cons/Challenges

- ⚠️ **Amount snapshot sync** - Requires application logic to update cell_references when transaction amount changes
- ⚠️ **Manual actual tracking** - Need UI to show when category uses override vs calculated
- ⚠️ **10-file limit** - Enforced in application, not database (could add trigger)
- ⚠️ **Complex queries** - Calculating category totals requires LEFT JOIN + COALESCE logic
- ⚠️ **Migration complexity** - Changing schema requires careful versioning

---

## Next Steps

1. **Create initial migration** (000001_create_initial_schema.up/down.sql)
2. **Add seed data migration** (000002_seed_default_categories.up.sql)
3. **Test migration rollback** (ensure .down.sql works)
4. **Implement Golang repository layer** (database access patterns)
5. **Add database connection pooling** (pgx pool configuration)
6. **Create query optimization tests** (EXPLAIN ANALYZE for slow queries)

---

## References

- [PostgreSQL NUMERIC Type](https://www.postgresql.org/docs/current/datatype-numeric.html)
- [PostgreSQL Don't Do This - Best Practices](https://wiki.postgresql.org/wiki/Don%27t_Do_This)
- [golang-migrate Examples](https://github.com/golang-migrate/migrate/tree/main/database/postgres/examples/migrations)
- [PostgreSQL Foreign Keys](https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-FK)
- [PostgreSQL Indexes](https://www.postgresql.org/docs/current/indexes.html)

---

**Next Research**: golang-migrate deep dive (11-postgresql-migrations.md), then API endpoint design (12-api-design.md)
