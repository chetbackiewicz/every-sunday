package handlers

import (
	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Auth          *AuthHandler
	Budget        *BudgetHandler
	Category      *CategoryHandler
	File          *FileHandler
	Transaction   *TransactionHandler
	CellReference *CellReferenceHandler
}

// NewHandlers creates new handler instances
func NewHandlers(repos *repositories.Repositories) *Handlers {
	return &Handlers{
		Auth:          NewAuthHandler(repos),
		Budget:        NewBudgetHandler(repos),
		Category:      NewCategoryHandler(repos),
		File:          NewFileHandler(repos),
		Transaction:   NewTransactionHandler(repos),
		CellReference: NewCellReferenceHandler(repos),
	}
}
