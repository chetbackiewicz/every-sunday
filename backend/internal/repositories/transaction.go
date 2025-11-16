package repositories

import (
	"context"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	// TODO: Implement
	return nil, nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id int) (*models.Transaction, error) {
	// TODO: Implement
	return nil, nil
}

func (r *transactionRepository) GetByCSVFile(ctx context.Context, csvFileID int) ([]*models.Transaction, error) {
	// TODO: Implement
	return nil, nil
}

func (r *transactionRepository) List(ctx context.Context, monthlyBudgetID int) ([]*models.Transaction, error) {
	// TODO: Implement with JOIN to csv_files
	return nil, nil
}

func (r *transactionRepository) Update(ctx context.Context, id int, description *string, category *string, note *string) error {
	// TODO: Implement
	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id int) error {
	// TODO: Implement (should also delete cell_references via CASCADE)
	return nil
}

func (r *transactionRepository) CountByCSVFile(ctx context.Context, csvFileID int) (int, error) {
	// TODO: Implement
	return 0, nil
}
