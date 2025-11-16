package models

import "time"

// User represents a user account
type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// MonthlyBudget represents a budget for a specific month
type MonthlyBudget struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Month     string    `db:"month" json:"month"` // Format: YYYY-MM
	Name      *string   `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Category represents a budget category with hierarchical structure
type Category struct {
	ID              int      `db:"id" json:"id"`
	MonthlyBudgetID int      `db:"monthly_budget_id" json:"monthly_budget_id"`
	Name            string   `db:"name" json:"name"`
	Projected       float64  `db:"projected" json:"projected"`
	ManualActual    *float64 `db:"manual_actual" json:"manual_actual"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	Children        []*Category `db:"-" json:"children,omitempty"` // For tree structure
}

// CategoryPath represents closure table for category hierarchy
type CategoryPath struct {
	AncestorID   int `db:"ancestor_id" json:"ancestor_id"`
	DescendantID int `db:"descendant_id" json:"descendant_id"`
	Depth        int `db:"depth" json:"depth"`
}

// CSVFile represents an uploaded CSV bank statement
type CSVFile struct {
	ID              int       `db:"id" json:"id"`
	MonthlyBudgetID int       `db:"monthly_budget_id" json:"monthly_budget_id"`
	Filename        string    `db:"filename" json:"filename"`
	UploadDate      time.Time `db:"upload_date" json:"upload_date"`
	FileSize        int       `db:"file_size" json:"file_size"`
	RowCount        int       `db:"row_count" json:"row_count"`
}

// Transaction represents a single row from CSV bank statement
type Transaction struct {
	ID              int       `db:"id" json:"id"`
	CSVFileID       int       `db:"csv_file_id" json:"csv_file_id"`
	RowNumber       int       `db:"row_number" json:"row_number"`
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"`
	PostDate        time.Time `db:"post_date" json:"post_date"`
	Description     string    `db:"description" json:"description"`
	Category        string    `db:"category" json:"category"`
	Type            string    `db:"type" json:"type"`
	Amount          float64   `db:"amount" json:"amount"` // Stored as NUMERIC(12,2) in DB
	Note            *string   `db:"note" json:"note"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// CellReference represents a link between a category cell and CSV transaction
type CellReference struct {
	ID            int       `db:"id" json:"id"`
	CategoryID    int       `db:"category_id" json:"category_id"`
	TransactionID int       `db:"transaction_id" json:"transaction_id"`
	Amount        float64   `db:"amount" json:"amount"` // Stored as NUMERIC(12,2) in DB
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
