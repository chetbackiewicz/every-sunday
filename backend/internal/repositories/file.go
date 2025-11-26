package repositories

import (
	"context"
	"fmt"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type fileRepository struct {
	db *sqlx.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *sqlx.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create inserts a new CSV file record
func (r *fileRepository) Create(ctx context.Context, file *models.CSVFile) (*models.CSVFile, error) {
	query := `
		INSERT INTO csv_files (monthly_budget_id, filename, file_size, row_count)
		VALUES ($1, $2, $3, $4)
		RETURNING id, upload_date
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		file.MonthlyBudgetID,
		file.Filename,
		file.FileSize,
		file.RowCount,
	).Scan(&file.ID, &file.UploadDate)

	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

// GetByID retrieves a file by ID
func (r *fileRepository) GetByID(ctx context.Context, id int) (*models.CSVFile, error) {
	file := &models.CSVFile{}
	query := `SELECT * FROM csv_files WHERE id = $1`
	
	err := r.db.GetContext(ctx, file, query, id)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// List retrieves all files for a monthly budget
func (r *fileRepository) List(ctx context.Context, monthlyBudgetID int) ([]*models.CSVFile, error) {
	files := []*models.CSVFile{}
	query := `SELECT * FROM csv_files WHERE monthly_budget_id = $1 ORDER BY upload_date DESC`
	
	err := r.db.SelectContext(ctx, &files, query, monthlyBudgetID)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Delete removes a file record
func (r *fileRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM csv_files WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("file not found")
	}

	return nil
}

// CountByBudget returns the number of files in a budget
func (r *fileRepository) CountByBudget(ctx context.Context, monthlyBudgetID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM csv_files WHERE monthly_budget_id = $1`
	
	err := r.db.GetContext(ctx, &count, query, monthlyBudgetID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

