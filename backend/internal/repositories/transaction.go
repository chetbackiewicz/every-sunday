package repositories

import (
	"context"
	"fmt"
	"strings"

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
	query := `
		INSERT INTO transactions (csv_file_id, row_number, transaction_date, post_date, description, category, type, amount, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		tx.CSVFileID,
		tx.RowNumber,
		tx.TransactionDate,
		tx.PostDate,
		tx.Description,
		tx.Category,
		tx.Type,
		tx.Amount,
		tx.Note,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return tx, nil
}

func (r *transactionRepository) CreateBatch(ctx context.Context, txs []*models.Transaction) error {
	if len(txs) == 0 {
		return nil
	}

	// Using sqlx.NamedExec or building a batch insert query
	// For better performance/portability with large datasets, let's construct a query
	// However, sqlx NamedExec is cleaner.
	
	query := `
		INSERT INTO transactions (csv_file_id, row_number, transaction_date, post_date, description, category, type, amount, note)
		VALUES (:csv_file_id, :row_number, :transaction_date, :post_date, :description, :category, :type, :amount, :note)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, txs)
	if err != nil {
		return fmt.Errorf("failed to batch insert transactions: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id int) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `SELECT * FROM transactions WHERE id = $1`
	err := r.db.GetContext(ctx, tx, query, id)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *transactionRepository) GetByCSVFile(ctx context.Context, csvFileID int) ([]*models.Transaction, error) {
	txs := []*models.Transaction{}
	query := `SELECT * FROM transactions WHERE csv_file_id = $1 ORDER BY row_number ASC`
	err := r.db.SelectContext(ctx, &txs, query, csvFileID)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepository) List(ctx context.Context, monthlyBudgetID int) ([]*models.Transaction, error) {
	txs := []*models.Transaction{}
	query := `
		SELECT t.* 
		FROM transactions t
		JOIN csv_files f ON t.csv_file_id = f.id
		WHERE f.monthly_budget_id = $1
		ORDER BY t.transaction_date DESC, t.id DESC
	`
	err := r.db.SelectContext(ctx, &txs, query, monthlyBudgetID)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepository) Update(ctx context.Context, id int, description *string, category *string, note *string) error {
	sets := []string{}
	args := []interface{}{}
	argID := 1

	if description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argID))
		args = append(args, *description)
		argID++
	}
	if category != nil {
		sets = append(sets, fmt.Sprintf("category = $%d", argID))
		args = append(args, *category)
		argID++
	}
	if note != nil {
		sets = append(sets, fmt.Sprintf("note = $%d", argID))
		args = append(args, *note)
		argID++
	}
	
	sets = append(sets, fmt.Sprintf("updated_at = NOW()"))

	if len(sets) == 0 { // Only updating timestamp if nothing else? Or no-op.
		return nil
	}

	query := fmt.Sprintf("UPDATE transactions SET %s WHERE id = $%d", strings.Join(sets, ", "), argID)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (r *transactionRepository) CountByCSVFile(ctx context.Context, csvFileID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM transactions WHERE csv_file_id = $1`
	err := r.db.GetContext(ctx, &count, query, csvFileID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
