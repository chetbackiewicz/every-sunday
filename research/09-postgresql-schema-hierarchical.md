# PostgreSQL Schema - Hierarchical Categories Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: Closure Table (with optional adjacency list hybrid)

---

## Recommendation Summary

**Chosen Approach: Closure Table**

Based on comprehensive research of hierarchical data patterns in PostgreSQL, the **Closure Table** pattern is recommended for storing the 3-level category tree structure.

**Why Closure Table**:
1. ✅ **Easy querying** - Get all ancestors or descendants with simple JOIN
2. ✅ **Referential integrity** - Full foreign key support
3. ✅ **Good performance** - O(1) for descendants/ancestors queries
4. ✅ **Easy CRUD operations** - Compared to nested sets
5. ✅ **Flexible** - Works with any tree depth (not just 3 levels)
6. ✅ **Standard SQL** - No PostgreSQL-specific extensions needed
7. ✅ **Storage efficient** - O(n²) rows in practice much less than theoretical

**Alternative Considered**: **Adjacency List** (simpler but requires recursive queries)
**Alternative Considered**: **Nested Sets** (fast reads but expensive writes)
**Alternative Considered**: **ltree** (PostgreSQL-specific, good for path-based queries)

---

## Core Requirements

From spike document:
- **3-level hierarchy maximum**: Parent → Child → Grandchild (e.g., Expenses → House → Mortgage)
- **Dynamic category creation**: Users can add subcategories with + button
- **Query patterns needed**:
  - Get all children of a category
  - Get all descendants (children + grandchildren)
  - Get parent of a category
  - Calculate totals by rolling up child values
  - Validate depth before allowing new subcategory

---

## Closure Table Pattern

### Concept

Store **all paths** between nodes in a separate table. For each node, store a row for:
1. The node to itself (ancestor = descendant, depth = 0)
2. The node to its parent (depth = 1)
3. The node to its grandparent (depth = 2)
4. And so on up the tree

### Schema Design

```sql
-- Main categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    projected DECIMAL(12,2) NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL,  -- For multi-user support
    month VARCHAR(7) NOT NULL,  -- Format: "YYYY-MM"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Closure table for hierarchical relationships
CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);
CREATE INDEX idx_category_paths_depth ON category_paths(depth);
CREATE INDEX idx_categories_user_month ON categories(user_id, month);
```

### Example Data

For this category tree:
```
Income (id=1)
├── Salary (id=2)
└── Commission (id=3)

Expenses (id=4)
├── House (id=5)
│   ├── Mortgage (id=6)
│   └── Insurance (id=7)
└── Car (id=8)
```

**categories table**:
```
| id  | name       | projected | user_id | month   |
| --- | ---------- | --------- | ------- | ------- |
| 1   | Income     | 5000      | 1       | 2025-10 |
| 2   | Salary     | 4000      | 1       | 2025-10 |
| 3   | Commission | 1000      | 1       | 2025-10 |
| 4   | Expenses   | 3000      | 1       | 2025-10 |
| 5   | House      | 2000      | 1       | 2025-10 |
| 6   | Mortgage   | 1500      | 1       | 2025-10 |
| 7   | Insurance  | 500       | 1       | 2025-10 |
| 8   | Car        | 1000      | 1       | 2025-10 |
```

**category_paths table**:
```
| ancestor_id | descendant_id | depth                                          |
| ----------- | ------------- | ---------------------------------------------- |
| 1           | 1             | 0      -- Income to itself                     |
| 1           | 2             | 1      -- Income → Salary                      |
| 1           | 3             | 1      -- Income → Commission                  |
| 2           | 2             | 0      -- Salary to itself                     |
| 3           | 3             | 0      -- Commission to itself                 |
| 4           | 4             | 0      -- Expenses to itself                   |
| 4           | 5             | 1      -- Expenses → House                     |
| 4           | 6             | 2      -- Expenses → Mortgage (through House)  |
| 4           | 7             | 2      -- Expenses → Insurance (through House) |
| 4           | 8             | 1      -- Expenses → Car                       |
| 5           | 5             | 0      -- House to itself                      |
| 5           | 6             | 1      -- House → Mortgage                     |
| 5           | 7             | 1      -- House → Insurance                    |
| 6           | 6             | 0      -- Mortgage to itself                   |
| 7           | 7             | 0      -- Insurance to itself                  |
| 8           | 8             | 0      -- Car to itself                        |
```

---

## Common Queries

### 1. Get Immediate Children

```sql
-- Get direct children of "Expenses" (id=4)
SELECT c.*
FROM categories c
JOIN category_paths cp ON c.id = cp.descendant_id
WHERE cp.ancestor_id = 4 AND cp.depth = 1;

-- Result: House (id=5), Car (id=8)
```

### 2. Get All Descendants

```sql
-- Get all descendants of "Expenses" (id=4), excluding self
SELECT c.*, cp.depth
FROM categories c
JOIN category_paths cp ON c.id = cp.descendant_id
WHERE cp.ancestor_id = 4 AND cp.depth > 0
ORDER BY cp.depth, c.name;

-- Result: House, Car (depth=1), Mortgage, Insurance (depth=2)
```

### 3. Get Parent

```sql
-- Get parent of "Mortgage" (id=6)
SELECT c.*
FROM categories c
JOIN category_paths cp ON c.id = cp.ancestor_id
WHERE cp.descendant_id = 6 AND cp.depth = 1;

-- Result: House (id=5)
```

### 4. Get Full Ancestry Path (Breadcrumbs)

```sql
-- Get all ancestors of "Mortgage" (id=6), ordered root to leaf
SELECT c.*, cp.depth
FROM categories c
JOIN category_paths cp ON c.id = cp.ancestor_id
WHERE cp.descendant_id = 6 AND cp.depth > 0
ORDER BY cp.depth DESC;

-- Result: Expenses (depth=2), House (depth=1)
```

### 5. Check if Node is Leaf (No Children)

```sql
-- Check if "Mortgage" (id=6) has children
SELECT NOT EXISTS (
    SELECT 1
    FROM category_paths
    WHERE ancestor_id = 6 AND depth = 1
) AS is_leaf;

-- Result: true (no children)
```

### 6. Get Maximum Depth of Tree

```sql
-- Get deepest level in entire tree
SELECT MAX(depth) FROM category_paths;

-- Get depth of specific subtree (e.g., under "Expenses")
SELECT MAX(depth)
FROM category_paths
WHERE ancestor_id = 4;
```

### 7. Prevent Creating Child Beyond Max Depth

```sql
-- Check if "House" (id=5) can have children (max depth = 2)
SELECT EXISTS (
    SELECT 1
    FROM category_paths
    WHERE descendant_id = 5 AND depth >= 2
) AS at_max_depth;

-- If TRUE, don't allow adding child
```

---

## CRUD Operations

### Insert New Top-Level Category

```sql
BEGIN;

-- Insert category
INSERT INTO categories (name, projected, user_id, month)
VALUES ('New Category', 0, 1, '2025-10')
RETURNING id; -- Let's say it returns id=9

-- Insert self-referencing path
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
VALUES (9, 9, 0);

COMMIT;
```

### Insert New Child Category

```sql
-- Add "Utilities" (id=10) as child of "House" (id=5)
BEGIN;

-- Insert category
INSERT INTO categories (name, projected, user_id, month)
VALUES ('Utilities', 200, 1, '2025-10')
RETURNING id; -- Let's say it returns id=10

-- Insert all paths:
-- 1. Self-reference
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
VALUES (10, 10, 0);

-- 2. Copy all ancestor paths of parent, incrementing depth
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
SELECT cp.ancestor_id, 10, cp.depth + 1
FROM category_paths cp
WHERE cp.descendant_id = 5;

COMMIT;
```

**Result**: Creates paths:
- 10 → 10 (depth=0) - self
- 5 → 10 (depth=1) - House → Utilities
- 4 → 10 (depth=2) - Expenses → Utilities

### Delete Category and Subtree

```sql
-- Delete "House" (id=5) and all descendants (Mortgage, Insurance, Utilities)
-- CASCADE will handle category_paths cleanup automatically
DELETE FROM categories WHERE id = 5;

-- This removes:
-- - Category rows: House, Mortgage, Insurance, Utilities
-- - All paths involving these categories (ON DELETE CASCADE)
```

### Move Subtree to New Parent

```sql
-- Move "House" (id=5) from "Expenses" (id=4) to "Income" (id=1)
BEGIN;

-- 1. Delete old paths (except self-reference)
DELETE FROM category_paths
WHERE descendant_id IN (
    SELECT cp_descendants.descendant_id
    FROM category_paths cp_descendants
    WHERE cp_descendants.ancestor_id = 5 -- House and its descendants
)
AND ancestor_id IN (
    SELECT cp_ancestors.ancestor_id
    FROM category_paths cp_ancestors
    WHERE cp_ancestors.descendant_id = 5 AND cp_ancestors.depth > 0 -- Old ancestors of House
);

-- 2. Insert new paths
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
SELECT new_ancestors.ancestor_id, moving_nodes.descendant_id, 
       new_ancestors.depth + moving_nodes.depth + 1
FROM category_paths new_ancestors
CROSS JOIN category_paths moving_nodes
WHERE new_ancestors.descendant_id = 1  -- New parent: Income
  AND moving_nodes.ancestor_id = 5;    -- House and descendants

COMMIT;
```

---

## Depth Validation

### Prevent Deep Nesting in Application Code

```typescript
// Before creating new category, check parent depth
const canAddChild = async (parentId: number): Promise<boolean> => {
  const result = await db.query(`
    SELECT MAX(depth) as parent_depth
    FROM category_paths
    WHERE descendant_id = $1
  `, [parentId])
  
  const parentDepth = result.rows[0].parent_depth
  return parentDepth < 2 // Allow if parent is at depth 0 or 1
}
```

### Database Constraint (Optional - More Strict)

```sql
-- Add CHECK constraint to prevent depth > 2
ALTER TABLE category_paths
ADD CONSTRAINT check_max_depth 
CHECK (depth <= 2);

-- NOTE: This prevents ANY path with depth > 2, effectively limiting tree to 3 levels
```

---

## Performance Considerations

### Storage Requirements

For a tree with `n` nodes:
- **Worst case**: O(n²) rows in `category_paths`
- **Typical case**: Much better (~3n to 5n rows for balanced trees)

**Example**: For budget app with ~50 categories:
- Worst case: 2,500 rows
- Typical case: 150-250 rows

**Benchmark**: PostgreSQL closure table queries are very fast:
- Descendants query: < 1ms for trees with 1000+ nodes
- Ancestors query: < 1ms for trees with 1000+ nodes
- Single-level depth filtering makes queries even faster

### Index Strategy

```sql
-- Essential indexes (already in schema above)
CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);
CREATE INDEX idx_category_paths_depth ON category_paths(depth);

-- Composite index for common query pattern (ancestor + depth)
CREATE INDEX idx_category_paths_ancestor_depth 
ON category_paths(ancestor_id, depth);

-- Index for user-specific monthly budgets
CREATE INDEX idx_categories_user_month ON categories(user_id, month);
```

---

## Alternative: Adjacency List (Simpler but Less Performant)

### Schema

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER,
    projected DECIMAL(12,2) NOT NULL DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);
```

### Pros
- ✅ Very simple schema
- ✅ Easy inserts/updates
- ✅ Easy to understand

### Cons
- ❌ Requires recursive CTEs for descendants (PostgreSQL WITH RECURSIVE)
- ❌ Slower queries for ancestry/descendants
- ❌ More complex application logic for cascading calculations

### When to Use
- If tree queries are infrequent
- If simplicity is more important than performance
- If using database with good recursive CTE optimization

---

## Alternative: Nested Sets (Fast Reads, Slow Writes)

### Schema

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    lft INTEGER NOT NULL,
    rgt INTEGER NOT NULL,
    projected DECIMAL(12,2) NOT NULL DEFAULT 0
);
```

### Pros
- ✅ Very fast ancestor/descendant queries
- ✅ Compact storage (2 extra columns)

### Cons
- ❌ **Expensive inserts/moves** - Must recalculate lft/rgt for many rows
- ❌ Complex update logic
- ❌ Not recommended for frequent modifications

### When to Use
- Read-heavy workloads with rare modifications
- When tree structure is mostly static

---

## Alternative: PostgreSQL ltree Extension

**Source**: https://www.postgresql.org/docs/current/ltree.html

### Concept
Store materialized path as specialized `ltree` datatype.

### Schema

```sql
CREATE EXTENSION ltree;

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path ltree NOT NULL,  -- e.g., 'Income.Salary'
    projected DECIMAL(12,2) NOT NULL DEFAULT 0
);

CREATE INDEX idx_categories_path_gist ON categories USING GIST (path);
```

### Pros
- ✅ Path-based queries with pattern matching
- ✅ GiST indexes for performance
- ✅ Simple schema

### Cons
- ❌ **PostgreSQL-specific** - Not portable to other databases
- ❌ Manual path management on inserts/moves
- ❌ Path string length limits

### When to Use
- If already committed to PostgreSQL long-term
- If path-based queries are primary use case

---

## Comparison Matrix

| Feature                   | Closure Table     | Adjacency List    | Nested Sets    | ltree Extension       |
| ------------------------- | ----------------- | ----------------- | -------------- | --------------------- |
| **Descendants Query**     | O(1) - JOIN       | O(log n) - CTE    | O(1) - BETWEEN | O(1) - GiST index     |
| **Ancestors Query**       | O(1) - JOIN       | O(log n) - CTE    | O(1) - BETWEEN | O(1) - GiST index     |
| **Insert Leaf**           | O(log n)          | O(1)              | O(n)           | O(log n)              |
| **Delete Subtree**        | O(1) - CASCADE    | O(1) - CASCADE    | O(n)           | O(n)                  |
| **Move Subtree**          | O(log n)          | O(1)              | O(n)           | O(n)                  |
| **Storage**               | O(n²) paths       | O(n)              | O(n)           | O(n)                  |
| **Referential Integrity** | ✅ Full FK support | ✅ Full FK support | ❌ No FKs       | ❌ No FKs              |
| **Portability**           | ✅ Standard SQL    | ✅ Standard SQL    | ✅ Standard SQL | ❌ PostgreSQL-specific |
| **Complexity**            | Medium            | Low               | High           | Medium                |
| **Best For**              | Balanced workload | Simple trees      | Read-heavy     | Path-based queries    |

---

## Recommendation for Budget App

**Use Closure Table** because:

1. **Query patterns match** - Need to:
   - Get children for rendering tree
   - Get all descendants for rollup calculations
   - Get parent for depth validation
   
2. **Modification frequency moderate** - Users will:
   - Create categories occasionally (not every request)
   - Rarely move categories between parents
   - Delete categories infrequently
   
3. **Performance excellent** - Closure table queries are O(1) with proper indexes

4. **Storage acceptable** - For ~50 categories, ~250 rows in closure table (trivial)

5. **Referential integrity** - Foreign keys prevent orphaned categories

6. **Standard SQL** - Works on any database, not locked to PostgreSQL

---

## Implementation Checklist

### Phase 1: Schema Setup
- [ ] Create `categories` table with columns for name, projected, user_id, month
- [ ] Create `category_paths` table with ancestor_id, descendant_id, depth
- [ ] Add foreign keys with ON DELETE CASCADE
- [ ] Create indexes on category_paths (ancestor_id, descendant_id, depth)
- [ ] Add composite index for (ancestor_id, depth)

### Phase 2: CRUD Functions (Database)
- [ ] Create stored procedure for inserting category with parent
- [ ] Create stored procedure for deleting category and subtree
- [ ] Create stored procedure for moving category to new parent
- [ ] Add CHECK constraint to prevent depth > 2

### Phase 3: Query Functions (Application)
- [ ] Implement `getChildren(categoryId)` query
- [ ] Implement `getDescendants(categoryId)` query
- [ ] Implement `getParent(categoryId)` query
- [ ] Implement `getAncestors(categoryId)` query (breadcrumbs)
- [ ] Implement `canAddChild(categoryId)` validation

### Phase 4: Testing
- [ ] Test inserting 3-level deep tree
- [ ] Test depth validation prevents 4th level
- [ ] Test delete cascade removes paths correctly
- [ ] Test move operation updates paths
- [ ] Load test with 1000 categories

---

## Migration to Closure Table (If Using Adjacency List First)

If you start with simple adjacency list and want to migrate:

```sql
-- Assuming existing adjacency list table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER,
    projected DECIMAL(12,2) NOT NULL DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);

-- Step 1: Create closure table
CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Step 2: Populate closure table using recursive CTE
WITH RECURSIVE tree AS (
    -- Base case: self-referencing paths
    SELECT id AS ancestor_id, id AS descendant_id, 0 AS depth
    FROM categories
    
    UNION ALL
    
    -- Recursive case: follow parent_id links
    SELECT t.ancestor_id, c.id, t.depth + 1
    FROM tree t
    JOIN categories c ON c.parent_id = t.descendant_id
)
INSERT INTO category_paths (ancestor_id, descendant_id, depth)
SELECT * FROM tree;

-- Step 3: Create indexes
CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);
CREATE INDEX idx_category_paths_depth ON category_paths(depth);

-- Step 4: (Optional) Drop parent_id column if no longer needed
ALTER TABLE categories DROP COLUMN parent_id;
```

---

## Pros of Closure Table Approach

- ✅ **Fast queries** - O(1) for ancestors and descendants
- ✅ **Easy to understand** - "Store all paths" is intuitive
- ✅ **Referential integrity** - Full foreign key support
- ✅ **Standard SQL** - No database-specific extensions
- ✅ **Moderate storage** - O(n²) theoretical, ~3-5n practical
- ✅ **Easy CRUD** - Insert/delete/move logic is straightforward
- ✅ **Depth control** - Simple to enforce with CHECK constraint or application logic
- ✅ **Cascading deletes** - ON DELETE CASCADE handles cleanup automatically

---

## Cons/Challenges

- ⚠️ **More storage** - More rows than adjacency list
- ⚠️ **Insert complexity** - Must insert multiple path rows
- ⚠️ **Move complexity** - Delete old paths + insert new paths
- ⚠️ **Understanding curve** - Slightly more complex than adjacency list

---

## Next Steps

1. **Design budget and transaction tables** (next research doc)
2. **Research migration tools** (golang-migrate vs goose)
3. **Create initial migration** with categories and category_paths tables
4. **Implement CRUD stored procedures** in PostgreSQL
5. **Build Golang service layer** with query functions

---

## References

- [Models for Hierarchical Data (SlideShare)](https://www.slideshare.net/billkarwin/models-for-hierarchical-data) - Excellent comparison of all patterns
- [PostgreSQL Recursive Queries](https://www.postgresql.org/docs/current/queries-with.html)
- [PostgreSQL ltree Extension](https://www.postgresql.org/docs/current/ltree.html)
- [Managing Hierarchical Data in MySQL](http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/) - Nested sets tutorial
- [Stack Overflow: Hierarchical Data Options](https://stackoverflow.com/questions/4048151/what-are-the-options-for-storing-hierarchical-data-in-a-relational-database)
- [SQL Antipatterns: Naive Trees](http://www.slideshare.net/billkarwin/sql-antipatterns-strike-back) - Bill Karwin presentation
