# Budget Tracking Application - Copilot Instructions

## Project Overview

A full-stack web application for personal budget tracking across multiple months with CSV bank statement import, hierarchical category management, and cell-based transaction linking.

**Architecture Status**: ✅ Research Complete (14 docs) - Ready for Implementation  
**Research Location**: `/research/` directory - Refer to these docs for all implementation guidance

---

## Product Requirements Summary

### Core Features
1. **CSV Upload** - Drag & drop up to 10 bank statement CSV files per month
2. **Editable Transaction Table** - Display and edit all CSV rows (Transaction Date, Post Date, Description, Category, Type, Amount, Note)
3. **Hierarchical Categories** - 3-level tree structure (Parent → Child → Grandchild)
   - Example: Expenses → House → Mortgage
   - Depth limit: 2 levels below parent (no + button on leaf nodes)
4. **Cell Linking** - Click category "Actual" cell → click CSV transaction cells → auto-sum amounts
5. **Referential Integrity** - Delete CSV row → remove from all linked category cells
6. **Manual Overrides** - User can type directly into any cell including "Actual" values
7. **Auto-Calculations** - Real-time updates for Projected/Actual/Difference and Remaining formula
8. **Multi-Month Management** - Dropdown selector to switch between monthly budgets

### Budget Categories
- Income (with Projected, Actual, Difference columns)
- Savings (Pre-Tax)
- Savings (Post-Tax)
- Expenses
- Remaining = Income - Savings (Pre-Tax) - Savings (Post-Tax) - Expenses

---

## Technology Stack (Research-Validated)

### Frontend (React + TypeScript)
**Philosophy**: Headless UI libraries for maximum design control

| Technology | Version | Purpose | Research Doc |
|------------|---------|---------|--------------|
| **React** | 18+ | UI framework | All docs |
| **TypeScript** | 5+ | Type safety | All docs |
| **PapaParse** | v5 | CSV parsing with streaming | 01-csv-parsing.md |
| **react-dropzone** | Latest | File upload with 10-file limit | 01-csv-parsing.md |
| **TanStack Table** | v8 | Editable table (headless) | 02-editable-tables.md |
| **react-arborist** | Latest | Category tree with drag & drop | 03-category-tree.md |

**State Management**: useReducer at app root for complex state coordination (04-cell-linking.md)

### Backend (Golang)
**Philosophy**: Simple, explicit, lightweight

| Technology | Version | Purpose | Research Doc |
|------------|---------|---------|--------------|
| **Gin** | Latest | Web framework (Express-like API) | 08-golang-framework.md |
| **sqlx** | Latest | Database queries (explicit SQL) | 13-database-integration.md |
| **pgx/v5** | v5 | PostgreSQL driver with connection pooling | 13-database-integration.md |
| **golang-migrate** | Latest | Schema migrations | 11-postgresql-migrations.md |
| **golang-jwt/jwt** | v5 | JWT authentication | 14-auth.md |
| **bcrypt** | Latest | Password hashing (cost 12) | 14-auth.md |

### Database (PostgreSQL)
**Design**: Normalized schema with closure table for hierarchy

| Feature | Implementation | Research Doc |
|---------|---------------|--------------|
| **Hierarchical Categories** | Closure table pattern | 09-postgresql-schema-hierarchical.md |
| **7 Core Tables** | users, monthly_budgets, categories, category_paths, csv_files, transactions, cell_references | 10-postgresql-schema-budget-tables.md |
| **Financial Precision** | NUMERIC(12,2) for all amounts (NOT float) | 10-postgresql-schema-budget-tables.md |
| **Referential Integrity** | Foreign keys with ON DELETE CASCADE | 10-postgresql-schema-budget-tables.md |

---

## Implementation Guidelines

### When Writing Frontend Code

#### CSV Upload Component
```typescript
// ✅ CORRECT: Use react-dropzone + PapaParse
import { useDropzone } from 'react-dropzone';
import Papa from 'papaparse';

const { getRootProps, getInputProps } = useDropzone({
  accept: { 'text/csv': ['.csv'] },
  maxFiles: 10,  // Enforce 10-file limit
  multiple: true,
  onDrop: (files) => {
    files.forEach(file => {
      Papa.parse(file, {
        header: true,           // Map columns to object keys
        dynamicTyping: true,    // Auto-convert numbers
        skipEmptyLines: true,
        complete: (results) => {
          dispatch({ type: 'CSV_UPLOADED', data: results.data });
        }
      });
    });
  }
});
```
**Reference**: research/01-csv-parsing.md

#### Editable Table Component
```typescript
// ✅ CORRECT: Use TanStack Table with custom cell editing
import { useReactTable, getCoreRowModel } from '@tanstack/react-table';

const table = useReactTable({
  data: transactions,
  columns,
  getCoreRowModel: getCoreRowModel(),
  meta: {
    updateData: (rowIndex, columnId, value) => {
      dispatch({ 
        type: 'CSV_CELL_EDITED', 
        csvIndex, 
        rowIndex, 
        field: columnId, 
        value 
      });
    }
  }
});

// Column definition with editable cell
const defaultColumn = {
  cell: ({ getValue, row, column, table }) => {
    const [value, setValue] = useState(getValue());
    return (
      <input
        value={value}
        onChange={e => setValue(e.target.value)}
        onBlur={() => table.options.meta?.updateData(row.index, column.id, value)}
      />
    );
  }
};
```
**Reference**: research/02-editable-tables.md

#### Category Tree Component
```typescript
// ✅ CORRECT: Use react-arborist with depth limiting
import { Tree } from 'react-arborist';

<Tree
  data={categories}
  disableDrop={(args) => args.parentNode.level >= 2}  // Max 3 levels
  onCreate={({ parentNode }) => {
    if (parentNode && parentNode.level >= 2) {
      return null;  // Prevent adding beyond grandchild level
    }
    return { id: newId, name: '' };
  }}
>
  {({ node }) => (
    <div>
      <span>{node.data.name}</span>
      {node.level < 2 && (  // Only show + button if depth allows
        <button onClick={() => tree.onCreate({ parentId: node.id })}>
          + Add
        </button>
      )}
    </div>
  )}
</Tree>
```
**Reference**: research/03-category-tree.md

#### State Management Pattern
```typescript
// ✅ CORRECT: useReducer at app root for complex state
interface AppState {
  csvFiles: { filename: string; data: Transaction[] }[];
  categories: CategoryNode[];
  activeCategoryId: string | null;
  currentMonth: string;
}

type AppAction =
  | { type: 'CSV_UPLOADED'; csvIndex: number; data: Transaction[] }
  | { type: 'CATEGORY_ACTUAL_CLICKED'; categoryId: string }
  | { type: 'CSV_CELL_CLICKED'; csvIndex: number; rowIndex: number }
  | { type: 'CSV_ROW_DELETED'; csvIndex: number; rowIndex: number };

function appReducer(state: AppState, action: AppAction): AppState {
  switch (action.type) {
    case 'CSV_ROW_DELETED':
      // Remove from CSV data
      const updatedCsvFiles = state.csvFiles.map((file, i) => 
        i === action.csvIndex 
          ? { ...file, data: file.data.filter((_, j) => j !== action.rowIndex) }
          : file
      );
      // Remove from all linked categories
      const updatedCategories = removeCellReferences(
        state.categories, 
        action.csvIndex, 
        action.rowIndex
      );
      return { ...state, csvFiles: updatedCsvFiles, categories: updatedCategories };
    // ... other cases
  }
}
```
**Reference**: research/04-cell-linking.md

#### Calculations Pattern
```typescript
// ✅ CORRECT: Derived state with useMemo, NOT useEffect
function CategoryRow({ category }: { category: CategoryNode }) {
  // Calculate actual from linked cells (don't store in state)
  const actual = useMemo(() => {
    if (category.manualActual !== null) {
      return category.manualActual;  // User override
    }
    return category.linkedCells.reduce((sum, ref) => sum + ref.amount, 0);
  }, [category.linkedCells, category.manualActual]);

  const difference = category.projected - actual;

  return (
    <tr>
      <td>{category.name}</td>
      <td>{category.projected}</td>
      <td>{actual}</td>
      <td>{difference}</td>
    </tr>
  );
}

// ❌ WRONG: Don't use useEffect for calculations
// useEffect(() => setActual(...), [linkedCells]);  // NO!
```
**Reference**: research/05-calculations.md

### When Writing Backend Code

#### API Routing with Gin
```go
// ✅ CORRECT: RESTful routes with Gin
func setupRoutes(r *gin.Engine, deps Dependencies) {
    api := r.Group("/api/v1")
    
    // Auth routes (public)
    auth := api.Group("/auth")
    {
        auth.POST("/register", deps.AuthHandler.Register)
        auth.POST("/login", deps.AuthHandler.Login)
        auth.POST("/refresh", deps.AuthHandler.RefreshToken)
        auth.POST("/logout", deps.AuthHandler.Logout)
    }
    
    // Protected routes
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware(deps.JWTConfig))
    {
        protected.GET("/auth/me", deps.AuthHandler.Me)
        
        // Budget routes
        budgets := protected.Group("/budgets")
        {
            budgets.GET("/:month", deps.BudgetHandler.Get)
            budgets.POST("", deps.BudgetHandler.Create)
            budgets.PUT("/:month", deps.BudgetHandler.Update)
        }
        
        // Category routes with hierarchy
        categories := protected.Group("/budgets/:month/categories")
        {
            categories.GET("", deps.CategoryHandler.List)  // Returns tree structure
            categories.POST("", deps.CategoryHandler.Create)
            categories.PUT("/:id", deps.CategoryHandler.Update)
            categories.DELETE("/:id", deps.CategoryHandler.Delete)
        }
        
        // CSV upload
        protected.POST("/budgets/:month/files", deps.FileHandler.Upload)
    }
}
```
**Reference**: research/08-golang-framework.md, research/12-api-design.md

#### Database Access with sqlx + Repository Pattern
```go
// ✅ CORRECT: Repository interface with explicit SQL
type UserRepository interface {
    Create(ctx context.Context, email, passwordHash string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByID(ctx context.Context, id int) (*User, error)
}

type userRepository struct {
    db *sqlx.DB
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
    err := r.db.GetContext(ctx, &user, query, email)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil  // Not found (not an error)
    }
    return &user, err
}

// ❌ WRONG: Don't use GORM
// db.Where("email = ?", email).First(&user)  // NO! Use explicit SQL
```
**Reference**: research/13-database-integration.md

#### Connection Pooling Configuration
```go
// ✅ CORRECT: Configure connection pool for low-traffic app
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
    
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // Budget app is low-traffic, use conservative settings
    db.SetMaxOpenConns(25)                     // Max concurrent connections
    db.SetMaxIdleConns(5)                      // Idle connections ready
    db.SetConnMaxLifetime(30 * time.Minute)    // Recycle connections
    db.SetConnMaxIdleTime(5 * time.Minute)     // Close unused idle
    
    return db, nil
}
```
**Reference**: research/13-database-integration.md

#### JWT Authentication Middleware
```go
// ✅ CORRECT: JWT middleware for protected routes
func AuthMiddleware(jwtConfig auth.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": gin.H{"code": "MISSING_TOKEN"}})
            c.Abort()
            return
        }
        
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": gin.H{"code": "INVALID_TOKEN_FORMAT"}})
            c.Abort()
            return
        }
        
        claims, err := auth.ValidateAccessToken(jwtConfig, parts[1])
        if err != nil {
            c.JSON(401, gin.H{"error": gin.H{"code": "INVALID_TOKEN"}})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Next()
    }
}
```
**Reference**: research/14-auth.md

### When Writing Database Code

#### Migration Files
```sql
-- ✅ CORRECT: Sequential migration with explicit up/down
-- db/migrations/000001_create_initial_schema.up.sql
BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE monthly_budgets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    month VARCHAR(7) NOT NULL,  -- Format: "YYYY-MM"
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, month),
    CHECK (month ~ '^\d{4}-\d{2}$')
);

CREATE INDEX idx_monthly_budgets_user_month ON monthly_budgets(user_id, month);

COMMIT;

-- db/migrations/000001_create_initial_schema.down.sql
BEGIN;
DROP TABLE IF EXISTS monthly_budgets CASCADE;
DROP TABLE IF EXISTS users CASCADE;
COMMIT;
```
**Reference**: research/11-postgresql-migrations.md

#### Closure Table for Hierarchy
```sql
-- ✅ CORRECT: Closure table for efficient hierarchy queries
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    projected NUMERIC(12,2) NOT NULL DEFAULT 0,
    manual_actual NUMERIC(12,2) DEFAULT NULL,
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE
);

CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Query all descendants of a category (e.g., get House → Mortgage + Insurance)
SELECT c.*, cp.depth
FROM categories c
JOIN category_paths cp ON c.id = cp.descendant_id
WHERE cp.ancestor_id = ? AND cp.depth > 0
ORDER BY cp.depth, c.name;
```
**Reference**: research/09-postgresql-schema-hierarchical.md

#### Financial Data Types
```sql
-- ✅ CORRECT: Use NUMERIC(12,2) for exact decimal precision
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    csv_file_id INTEGER NOT NULL,
    amount NUMERIC(12,2) NOT NULL,  -- Exact decimal, NO FLOAT!
    -- ... other columns
);

-- ❌ WRONG: Don't use FLOAT, DOUBLE, or MONEY type
-- amount FLOAT  -- NO! Causes rounding errors
-- amount MONEY  -- NO! PostgreSQL-specific, limited precision
```
**Reference**: research/10-postgresql-schema-budget-tables.md

---

## Common Patterns & Best Practices

### Frontend Patterns

#### ✅ DO: Lift State Up for Complex Coordination
```typescript
// Complex state (CSV + Categories + Cell Links) managed at app root
const [state, dispatch] = useReducer(appReducer, initialState);

// Pass dispatch down to components
<CategoryTree dispatch={dispatch} categories={state.categories} />
<CSVTable dispatch={dispatch} data={state.csvFiles[0].data} />
```

#### ✅ DO: Use Derived State for Calculations
```typescript
// Calculate during render, don't store in state
const total = useMemo(() => 
  items.reduce((sum, item) => sum + item.amount, 0),
  [items]
);
```

#### ❌ DON'T: Use useEffect for Calculations
```typescript
// ❌ WRONG: Causes extra re-renders
useEffect(() => {
  setTotal(items.reduce((sum, item) => sum + item.amount, 0));
}, [items]);
```

### Backend Patterns

#### ✅ DO: Use Repository Pattern
```typescript
// Separate database access from HTTP handlers
type BudgetRepository interface {
    GetByMonth(ctx, userId, month) (*Budget, error)
    Create(ctx, userId, month) (*Budget, error)
}

// Handler uses repository
func (h *BudgetHandler) Get(c *gin.Context) {
    userId := c.GetInt("user_id")
    month := c.Param("month")
    budget, err := h.repo.GetByMonth(c.Request.Context(), userId, month)
    // ...
}
```

#### ✅ DO: Use Context for Cancellation
```go
// Always pass context through the call stack
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, query, email)  // Uses context
    return &user, err
}
```

#### ❌ DON'T: Use ORM Query Builders
```go
// ❌ WRONG: Use explicit SQL, not GORM
// db.Where("user_id = ?", userId).Find(&budgets)  // NO!

// ✅ RIGHT: Explicit SQL with sqlx
query := `SELECT * FROM budgets WHERE user_id = $1`
err := r.db.SelectContext(ctx, &budgets, query, userId)
```

### Database Patterns

#### ✅ DO: Use Foreign Keys with CASCADE
```sql
-- Automatic cleanup when parent deleted
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
```

#### ✅ DO: Add Indexes on Foreign Keys
```sql
CREATE INDEX idx_categories_budget ON categories(monthly_budget_id);
CREATE INDEX idx_transactions_csv_file ON transactions(csv_file_id);
```

#### ✅ DO: Use CHECK Constraints for Validation
```sql
month VARCHAR(7) NOT NULL,
CHECK (month ~ '^\d{4}-\d{2}$')  -- Validates YYYY-MM format
```

---

## Critical Implementation Notes

### 1. Cell Reference ID Mapping
**Issue**: Frontend uses `csvIndex + rowIndex` but backend uses `transactionId`.

**Solution**: Include both in CellReference interface:
```typescript
interface CellReference {
  transactionId: number  // For backend API sync
  csvIndex: number       // For frontend UI reference
  rowIndex: number       // For frontend UI reference
  amount: number
  cellId: string         // "csv0_row5" for React keys
}
```

### 2. 10-File Limit Enforcement
Enforce at multiple levels:
- **Frontend**: `react-dropzone` with `maxFiles: 10`
- **Backend**: Query count before accepting upload
- **Application**: Business logic validation

```go
func (h *FileHandler) Upload(c *gin.Context) {
    budgetId := getBudgetId(c)
    count := h.repo.CountFiles(budgetId)
    if count >= 10 {
        c.JSON(400, gin.H{"error": "Maximum 10 files per budget"})
        return
    }
    // ... proceed with upload
}
```

### 3. Transaction Amount Precision
Always use exact decimals:
- **Frontend**: Store as `number`, display with `.toFixed(2)`
- **Backend**: Parse as `float64`, store in `NUMERIC(12,2)`
- **Database**: `NUMERIC(12,2)` type (never FLOAT or MONEY)

### 4. Timezone Handling
- **Transaction dates**: Use `DATE` type (calendar dates, no timezone)
- **System timestamps**: Use `TIMESTAMPTZ` (created_at, updated_at)
- **API format**: ISO 8601 strings ("2025-10-19" for dates, "2025-10-19T14:30:00Z" for timestamps)

---

## Quick Reference: Where to Find Answers

| Question | Research Document |
|----------|-------------------|
| How to parse CSV files? | 01-csv-parsing.md |
| How to build editable tables? | 02-editable-tables.md |
| How to create category tree? | 03-category-tree.md |
| How to manage complex state? | 04-cell-linking.md |
| How to calculate totals? | 05-calculations.md |
| How to switch between months? | 06-multi-month-ux.md |
| Which Go framework to use? | 08-golang-framework.md |
| How to model hierarchical data? | 09-postgresql-schema-hierarchical.md |
| What database tables are needed? | 10-postgresql-schema-budget-tables.md |
| How to manage schema changes? | 11-postgresql-migrations.md |
| What API endpoints to create? | 12-api-design.md |
| How to access the database? | 13-database-integration.md |
| How to implement authentication? | 14-auth.md |
| Are there any contradictions? | ARCHITECTURE_REVIEW.md |

---

## Architecture Validation

✅ **All research complete** - 14 comprehensive documents  
✅ **Architecture reviewed** - No major contradictions found  
✅ **Technology validated** - All libraries work together  
✅ **Requirements covered** - 100% of spike requirements addressed  

**Overall Score: 9.4/10** - Production-ready with minor adjustments

**Status**: Ready for implementation phase

---

## Project Structure

```
every-sunday/
├── frontend/                 # React + TypeScript
│   ├── src/
│   │   ├── components/
│   │   │   ├── CSVUpload.tsx
│   │   │   ├── TransactionTable.tsx
│   │   │   ├── CategoryTree.tsx
│   │   │   └── MonthSelector.tsx
│   │   ├── hooks/
│   │   │   └── useAppState.ts
│   │   ├── reducers/
│   │   │   └── appReducer.ts
│   │   └── App.tsx
│   └── package.json
├── backend/                  # Golang
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── database/
│   │   ├── repositories/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── auth/
│   ├── db/migrations/
│   └── go.mod
├── research/                 # ✅ COMPLETE - Reference these docs
│   ├── 01-csv-parsing.md
│   ├── 02-editable-tables.md
│   ├── ... (14 total)
│   └── ARCHITECTURE_REVIEW.md
└── README.md
```

---

## Next Steps for Implementation

1. **Backend Setup**
   - Initialize Go module: `go mod init github.com/user/every-sunday-backend`
   - Install dependencies: Gin, sqlx, pgx, golang-migrate, golang-jwt
   - Create initial migration: `migrate create -ext sql -dir db/migrations -seq create_initial_schema`
   - Set up database connection with pooling (doc 13)

2. **Database Setup**
   - Create PostgreSQL database
   - Run migrations: `migrate -path db/migrations -database "postgres://..." up`
   - Verify schema with psql

3. **Backend Implementation**
   - Implement repositories (User, Budget, Category, Transaction)
   - Implement handlers (Auth, Budget, Category, File)
   - Implement middleware (JWT auth, CORS, error handling)
   - Wire up routes in main.go

4. **Frontend Setup**
   - Create React app: `npm create vite@latest frontend -- --template react-ts`
   - Install dependencies: PapaParse, react-dropzone, TanStack Table, react-arborist
   - Set up app-level useReducer state management
   - Create component structure

5. **Frontend Implementation**
   - Build CSV upload component
   - Build transaction table with editing
   - Build category tree with depth limiting
   - Implement cell linking interaction
   - Add calculations with useMemo

6. **Integration & Testing**
   - Connect frontend to backend API
   - Test full data flow: Upload → Parse → Store → Retrieve
   - Test cell linking and deletion
   - Test multi-month switching

---

**Last Updated**: October 19, 2025  
**Research Phase**: ✅ Complete  
**Implementation Phase**: Ready to begin