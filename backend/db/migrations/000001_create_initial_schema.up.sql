-- Initial schema migration for Every Sunday budget app
-- Reference: research/09-postgresql-schema-hierarchical.md and research/10-postgresql-schema-budget-tables.md
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
    month VARCHAR(7) NOT NULL,  -- Format: "YYYY-MM"
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, month),
    CHECK (month ~ '^\d{4}-\d{2}$')
);

CREATE INDEX idx_monthly_budgets_user_month ON monthly_budgets(user_id, month);

-- Categories table (3-level hierarchy support)
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

-- Category paths (closure table for hierarchical queries)
CREATE TABLE category_paths (
    ancestor_id INTEGER NOT NULL,
    descendant_id INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (ancestor_id, descendant_id),
    FOREIGN KEY (ancestor_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES categories(id) ON DELETE CASCADE,
    CHECK (depth >= 0 AND depth <= 2)  -- Max 3 levels (0=self, 1=child, 2=grandchild)
);

CREATE INDEX idx_category_paths_ancestor ON category_paths(ancestor_id);
CREATE INDEX idx_category_paths_descendant ON category_paths(descendant_id);

-- CSV files table
CREATE TABLE csv_files (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_size INTEGER NOT NULL,
    row_count INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (monthly_budget_id) REFERENCES monthly_budgets(id) ON DELETE CASCADE
);

CREATE INDEX idx_csv_files_budget ON csv_files(monthly_budget_id);

-- Transactions table (CSV row data)
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    csv_file_id INTEGER NOT NULL,
    row_number INTEGER NOT NULL,
    transaction_date DATE NOT NULL,
    post_date DATE NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(255),
    type VARCHAR(50),
    amount NUMERIC(12,2) NOT NULL,  -- Exact decimal precision for financial data
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (csv_file_id) REFERENCES csv_files(id) ON DELETE CASCADE
);

CREATE INDEX idx_transactions_csv_file ON transactions(csv_file_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

-- Cell references table (links category cells to CSV transactions)
CREATE TABLE cell_references (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    transaction_id INTEGER NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    UNIQUE(category_id, transaction_id)  -- Prevent duplicate links
);

CREATE INDEX idx_cell_references_category ON cell_references(category_id);
CREATE INDEX idx_cell_references_transaction ON cell_references(transaction_id);

COMMIT;
