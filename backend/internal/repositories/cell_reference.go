package repositories

import (
	"context"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type cellReferenceRepository struct {
	db *sqlx.DB
}

func NewCellReferenceRepository(db *sqlx.DB) CellReferenceRepository {
	return &cellReferenceRepository{db: db}
}

func (r *cellReferenceRepository) Create(ctx context.Context, categoryID, transactionID int, amount float64) (*models.CellReference, error) {
	// TODO: Implement
	return nil, nil
}

func (r *cellReferenceRepository) GetByCategory(ctx context.Context, categoryID int) ([]*models.CellReference, error) {
	// TODO: Implement
	return nil, nil
}

func (r *cellReferenceRepository) GetByTransaction(ctx context.Context, transactionID int) ([]*models.CellReference, error) {
	// TODO: Implement
	return nil, nil
}

func (r *cellReferenceRepository) Delete(ctx context.Context, id int) error {
	// TODO: Implement
	return nil
}

func (r *cellReferenceRepository) DeleteByTransaction(ctx context.Context, transactionID int) error {
	// TODO: Implement (for referential integrity when transaction deleted)
	return nil
}
