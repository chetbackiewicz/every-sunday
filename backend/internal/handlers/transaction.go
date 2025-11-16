package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	repos *repositories.Repositories
}

func NewTransactionHandler(repos *repositories.Repositories) *TransactionHandler {
	return &TransactionHandler{repos: repos}
}

func (h *TransactionHandler) List(c *gin.Context) {
	// TODO: List all transactions for a budget month
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *TransactionHandler) Update(c *gin.Context) {
	// TODO: Update transaction fields (description, category, note)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	// TODO: Delete transaction (cascade removes cell_references)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
