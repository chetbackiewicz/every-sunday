package repositories

import (
	"context"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// Repositories holds all repository interfaces
type Repositories struct {
	User          UserRepository
	Budget        BudgetRepository
	Category      CategoryRepository
	Transaction   TransactionRepository
	CellReference CellReferenceRepository
	File          FileRepository
}

// NewRepositories creates new repository instances
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:          NewUserRepository(db),
		Budget:        NewBudgetRepository(db),
		Category:      NewCategoryRepository(db),
		Transaction:   NewTransactionRepository(db),
		CellReference: NewCellReferenceRepository(db),
		File:          NewFileRepository(db),
	}
}

// UserRepository defines user data access methods
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
}

// BudgetRepository defines budget data access methods
type BudgetRepository interface {
	Create(ctx context.Context, userID int, month string, name *string) (*models.MonthlyBudget, error)
	GetByMonth(ctx context.Context, userID int, month string) (*models.MonthlyBudget, error)
	List(ctx context.Context, userID int) ([]*models.MonthlyBudget, error)
	Update(ctx context.Context, id int, name *string) error
	Delete(ctx context.Context, id int) error
}

// CategoryRepository defines category data access methods
type CategoryRepository interface {
	Create(ctx context.Context, monthlyBudgetID int, name string, projected float64, parentID *int) (*models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	GetTree(ctx context.Context, monthlyBudgetID int) ([]*models.Category, error)
	GetChildren(ctx context.Context, categoryID int) ([]*models.Category, error)
	GetDepth(ctx context.Context, categoryID int) (int, error)
	Update(ctx context.Context, id int, name *string, projected *float64, manualActual *float64) error
	Delete(ctx context.Context, id int) error
}

// TransactionRepository defines transaction data access methods
type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	CreateBatch(ctx context.Context, txs []*models.Transaction) error
	GetByID(ctx context.Context, id int) (*models.Transaction, error)
	GetByCSVFile(ctx context.Context, csvFileID int) ([]*models.Transaction, error)
	List(ctx context.Context, monthlyBudgetID int) ([]*models.Transaction, error)
	Update(ctx context.Context, id int, description *string, category *string, note *string) error
	Delete(ctx context.Context, id int) error
	CountByCSVFile(ctx context.Context, csvFileID int) (int, error)
}

// CellReferenceRepository defines cell reference data access methods
type CellReferenceRepository interface {
	Create(ctx context.Context, categoryID, transactionID int, amount float64) (*models.CellReference, error)
	GetByCategory(ctx context.Context, categoryID int) ([]*models.CellReference, error)
	GetByTransaction(ctx context.Context, transactionID int) ([]*models.CellReference, error)
	Delete(ctx context.Context, id int) error
	DeleteByTransaction(ctx context.Context, transactionID int) error
}

// FileRepository defines file data access methods
type FileRepository interface {
	Create(ctx context.Context, file *models.CSVFile) (*models.CSVFile, error)
	GetByID(ctx context.Context, id int) (*models.CSVFile, error)
	List(ctx context.Context, monthlyBudgetID int) ([]*models.CSVFile, error)
	Delete(ctx context.Context, id int) error
	CountByBudget(ctx context.Context, monthlyBudgetID int) (int, error)
}
