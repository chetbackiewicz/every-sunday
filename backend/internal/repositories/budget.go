package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type budgetRepository struct {
	db *sqlx.DB
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *sqlx.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(ctx context.Context, userID int, month string, name *string) (*models.MonthlyBudget, error) {
	var budget models.MonthlyBudget
	query := `
		INSERT INTO monthly_budgets (user_id, month, name, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, user_id, month, name, created_at, updated_at
	`
	err := r.db.GetContext(ctx, &budget, query, userID, month, name)
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) GetByMonth(ctx context.Context, userID int, month string) (*models.MonthlyBudget, error) {
	var budget models.MonthlyBudget
	query := `SELECT id, user_id, month, name, created_at, updated_at FROM monthly_budgets WHERE user_id = $1 AND month = $2`
	err := r.db.GetContext(ctx, &budget, query, userID, month)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &budget, err
}

func (r *budgetRepository) List(ctx context.Context, userID int) ([]*models.MonthlyBudget, error) {
	var budgets []*models.MonthlyBudget
	query := `SELECT id, user_id, month, name, created_at, updated_at FROM monthly_budgets WHERE user_id = $1 ORDER BY month DESC`
	err := r.db.SelectContext(ctx, &budgets, query, userID)
	return budgets, err
}

func (r *budgetRepository) Update(ctx context.Context, id int, name *string) error {
	query := `UPDATE monthly_budgets SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, name, id)
	return err
}

func (r *budgetRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM monthly_budgets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
