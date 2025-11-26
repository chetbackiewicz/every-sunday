package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Create Category
	cat := &models.Category{
		MonthlyBudgetID: monthlyBudgetID,
		Name:            name,
		Projected:       projected,
	}
	
	query := `
		INSERT INTO categories (monthly_budget_id, name, projected)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, query, monthlyBudgetID, name, projected).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert category: %w", err)
	}

	// 2. Insert Self Path (depth 0)
	_, err = tx.ExecContext(ctx, `INSERT INTO category_paths (ancestor_id, descendant_id, depth) VALUES ($1, $1, 0)`, cat.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert self path: %w", err)
	}

	// 3. Insert Ancestor Paths if parent exists
	if parentID != nil {
		// Verify depth limit (max 2 levels deep means newly inserted can be at most depth 2, i.e. max existing parent depth is 1)
		// Check parent depth relative to root?
		// Actually, we just copy paths from parent.
		// Query: Insert (ancestor of parent, new node, depth+1)
		
		pathQuery := `
			INSERT INTO category_paths (ancestor_id, descendant_id, depth)
			SELECT p.ancestor_id, $1, p.depth + 1
			FROM category_paths p
			WHERE p.descendant_id = $2
		`
		_, err = tx.ExecContext(ctx, pathQuery, cat.ID, *parentID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert ancestor paths: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return cat, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	cat := &models.Category{}
	query := `SELECT * FROM categories WHERE id = $1`
	err := r.db.GetContext(ctx, cat, query, id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *categoryRepository) GetTree(ctx context.Context, monthlyBudgetID int) ([]*models.Category, error) {
	// Fetch all categories for budget
	// Also fetch their depth/parent info to reconstruct tree?
	// Or just fetch all and use closure table to find parents.
	// Actually, for simple tree construction, we need to know the parent of each node.
	// Parent is the ancestor with depth=1.
	
	query := `
		SELECT c.*, 
		       (SELECT ancestor_id FROM category_paths WHERE descendant_id = c.id AND depth = 1) as parent_id
		FROM categories c
		WHERE c.monthly_budget_id = $1
		ORDER BY c.name ASC
	`
	
	// models.Category doesn't have ParentID field mapped in struct tag usually, but let's check.
	// The struct definition in models.go doesn't have ParentID.
	// We can extend it or use a temporary struct.
	
	type CategoryWithParent struct {
		models.Category
		ParentID *int `db:"parent_id"`
	}
	
	var rows []CategoryWithParent
	err := r.db.SelectContext(ctx, &rows, query, monthlyBudgetID)
	if err != nil {
		return nil, err
	}

	// Reconstruct tree
	catMap := make(map[int]*models.Category)
	var rootCats []*models.Category

	// First pass: create nodes
	for _, row := range rows {
		cat := row.Category
		cat.Children = []*models.Category{}
		catMap[cat.ID] = &cat
	}

	// Second pass: link children
	for _, row := range rows {
		cat := catMap[row.ID]
		if row.ParentID != nil {
			if parent, ok := catMap[*row.ParentID]; ok {
				parent.Children = append(parent.Children, cat)
			}
		} else {
			// Root node (or at least no parent in this set/logic)
			rootCats = append(rootCats, cat)
		}
	}

	return rootCats, nil
}

func (r *categoryRepository) GetChildren(ctx context.Context, categoryID int) ([]*models.Category, error) {
	// Immediate children have depth=1 from this ancestor
	query := `
		SELECT c.*
		FROM categories c
		JOIN category_paths p ON c.id = p.descendant_id
		WHERE p.ancestor_id = $1 AND p.depth = 1
	`
	var cats []*models.Category
	err := r.db.SelectContext(ctx, &cats, query, categoryID)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *categoryRepository) GetDepth(ctx context.Context, categoryID int) (int, error) {
	// Depth is max path length from a root?
	// Or depth from root of tree.
	// Root has no ancestors with depth > 0.
	// Actually, we can count ancestors.
	var depth int
	query := `SELECT COUNT(*) - 1 FROM category_paths WHERE descendant_id = $1`
	err := r.db.GetContext(ctx, &depth, query, categoryID)
	return depth, err
}

func (r *categoryRepository) Update(ctx context.Context, id int, name *string, projected *float64, manualActual *float64) error {
	sets := []string{}
	args := []interface{}{}
	argID := 1

	if name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argID))
		args = append(args, *name)
		argID++
	}
	if projected != nil {
		sets = append(sets, fmt.Sprintf("projected = $%d", argID))
		args = append(args, *projected)
		argID++
	}
	if manualActual != nil {
		sets = append(sets, fmt.Sprintf("manual_actual = $%d", argID))
		args = append(args, *manualActual)
		argID++
	}
	
	sets = append(sets, "updated_at = NOW()")
	
	if len(sets) == 1 { // Only updated_at
		return nil
	}

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = $%d", strings.Join(sets, ", "), argID)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	// Cascade deletes paths and children in DB definition, so simple delete works
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
