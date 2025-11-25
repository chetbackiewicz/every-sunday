# Async Agent Guidelines

This document serves as the central reference for asynchronous agents working on the **Every Sunday** budget application. It distills architectural decisions, best practices, and implementation patterns from the research phase.

**Status**: Living Document  
**Goal**: Ensure consistency across autonomous development tasks.

---

## 1. Frontend Agent Guidelines

### Stack & Philosophy
- **Stack**: React 18+, TypeScript, Vite.
- **Philosophy**: **"Headless Everything"**. Use headless libraries for complex logic (tables, trees, drag-and-drop) and implement custom UI/CSS. Avoid opinionated UI frameworks (like Material UI) that fight against custom designs.
- **State Pattern**: **"Lift State Up"**. Complex state (like cell linking) lives at the app root using `useReducer` + `Context`. Calculations are **derived** during render, not stored in state.

### Core Libraries
- **Tables**: `TanStack Table v8` (Headless, highly customizable).
- **Category Tree**: `react-arborist` (Headless, built-in drag-and-drop, virtual scrolling).
- **CSV Parsing**: `PapaParse v5` (Streaming support, worker threads).
- **File Upload**: `react-dropzone` (Headless, simple hooks).
- **Styling**: Custom CSS/SCSS Modules or Tailwind (prefer consistency).

### Implementation Details

#### State Management (Cell Linking)
- **Single Source of Truth**: Use a root-level `useReducer`.
- **Action Flow**: `CATEGORY_ACTUAL_CLICKED` → `CSV_CELL_CLICKED` → Update `linkedCells` array in category node.
- **Data Structure**:
  ```typescript
  interface CellReference {
    transactionId: number // Database ID (critical for sync)
    csvIndex: number      // UI reference
    rowIndex: number      // UI reference
    amount: number        // Cached snapshot
  }
  ```
- **Deletion**: Deleting a CSV row MUST automatically remove related `CellReference`s from categories.

#### Calculations (Real-time)
- **Pattern**: **Derived State**. Never `useEffect` to sync state.
- **Optimization**: Use `useMemo` for expensive aggregations (budget summary, deep category trees).
- **Manual Overrides**:
  - Priority: `manualActual` > `SUM(linkedCells)` > `0`.
  - UI: Show indicator (e.g., "pencil" icon) when override is active.

#### Month Switching
- **UI**: Controlled `<select>` dropdown + Prev/Next buttons.
- **Persistence**:
  - **Load**: Lazy initializer in `useState` (`() => localStorage.getItem(...)`).
  - **Save**: `useEffect` to write to `localStorage` on state change.
  - **Switching**: Save current month state → Switch month ID → Load new month state.

---

## 2. Backend Agent Guidelines

### Stack & Philosophy
- **Stack**: Go 1.21+, Gin Web Framework.
- **Philosophy**: **"Simple, Fast, Explicit"**. Minimal abstraction. No heavy ORMs (use `sqlx`).
- **Error Handling**: Standardized JSON error responses.

### API Design (`/api/v1`)
- **Style**: RESTful, JSON-first.
- **Resource Naming**: Plural nouns (`/budgets`, `/categories`), kebab-case URLs.
- **Validation**: Struct tags (`binding:"required,gte=0"`) in Gin handlers.
- **Response Wrapper**: Standard envelope for success/error not strictly required for success, but mandatory for errors.
- **Pagination**: Support `limit` and `offset` for transaction lists.

### Authentication
- **Mechanism**: JWT (JSON Web Tokens).
- **Token Strategy**:
  - **Access Token**: Short-lived (15 min), sent in `Authorization: Bearer` header.
  - **Refresh Token**: Long-lived (7 days), stored in **HttpOnly Cookie**.
- **Security**: `bcrypt` (cost 12) for password hashing. HTTPS enforcement in production.

### Database Access
- **Library**: `sqlx` (Struct scanning, named parameters, explicit SQL).
- **Pattern**: **Repository Pattern**. Define interfaces (`UserRepository`, `BudgetRepository`) for testability.
- **Transactions**: Use manual `BEGIN`/`COMMIT` for multi-step operations (e.g., create budget + seed categories).
- **Connection Pooling**: Configure `pgx` pool (Max 25 open, 5 idle).

### File Handling
- **Upload**: Multipart/form-data (`c.FormFile`).
- **Constraints**: Max 10MB per file, Max 10 files per budget.
- **Storage**: Store metadata in DB (`csv_files` table), file content in object storage or local filesystem (interface abstracted).

---

## 3. Database Agent Guidelines

### Stack & Schema
- **Database**: PostgreSQL 14+.
- **Money Type**: `NUMERIC(12,2)` (Strictly avoid `FLOAT` or PostgreSQL `MONEY` type).
- **Dates**: `TIMESTAMPTZ` for audit fields, `DATE` for transaction/post dates.
- **Keys**: `SERIAL` or `BIGSERIAL` primary keys. Foreign keys always use `ON DELETE CASCADE` for parent-child relationships.

### Hierarchical Data (Categories)
- **Pattern**: **Closure Table** (`category_paths`).
- **Structure**: Stores all ancestor-descendant paths + depth.
- **Constraints**: Enforce max depth of 3 levels (Parent → Child → Grandchild) via application logic or CHECK constraint (`depth <= 2`).

### Tables Overview
- `users`: Auth and profile.
- `monthly_budgets`: Root container for a month.
- `categories`: Hierarchical nodes.
- `category_paths`: Closure table for hierarchy.
- `csv_files`: Uploaded file metadata.
- `transactions`: Individual rows from CSVs.
- `cell_references`: Many-to-many link between Categories and Transactions (with `amount_snapshot`).

### Migrations
- **Tool**: `golang-migrate`.
- **Versioning**: Sequential integers (`000001_name.up.sql`).
- **Safety**: Always wrap in `BEGIN; ... COMMIT;` (except concurrent index creation).
- **Reversibility**: Always implement `.down.sql`.

---

## 4. QA/Testing Agent Guidelines

### Strategy
- **Frontend**:
  - **Unit**: Jest/Vitest for reducers, calculation utilities (pure functions).
  - **Component**: React Testing Library for interactive components (tables, trees).
- **Backend**:
  - **Unit**: Mock repositories to test service logic/handlers.
  - **Integration**: Test API endpoints against a **real ephemeral test database** (Docker container), not mocks. Use `testify` assertions.
- **E2E**: Playwright/Cypress for critical flows (Login → Upload CSV → Link Cell).

### Critical Test Scenarios
1. **Cell Linking**: Verify sums update correctly when links are added/removed.
2. **File Limits**: Verify uploading >10 files returns 4xx error.
3. **Cascading Deletes**: Delete a budget -> Verify categories, files, and transactions are gone.
4. **Isolation**: Verify User A cannot access User B's budget.

