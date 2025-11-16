package repositories

import (
	"context"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, monthlyBudgetID int, name string, projected float64, parentID *int) (*models.Category, error) {
	// TODO: Implement with closure table path insertion
	return nil, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	// TODO: Implement
	return nil, nil
}

func (r *categoryRepository) GetTree(ctx context.Context, monthlyBudgetID int) ([]*models.Category, error) {
	// TODO: Implement with recursive query using closure table
	return nil, nil
}

func (r *categoryRepository) GetChildren(ctx context.Context, categoryID int) ([]*models.Category, error) {
	// TODO: Implement
	return nil, nil
}

func (r *categoryRepository) GetDepth(ctx context.Context, categoryID int) (int, error) {
	// TODO: Implement using category_paths depth column
	return 0, nil
}

func (r *categoryRepository) Update(ctx context.Context, id int, name *string, projected *float64, manualActual *float64) error {
	// TODO: Implement
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	// TODO: Implement (cascade will handle paths)
	return nil
}
