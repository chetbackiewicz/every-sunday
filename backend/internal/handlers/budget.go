package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	repos *repositories.Repositories
}

func NewBudgetHandler(repos *repositories.Repositories) *BudgetHandler {
	return &BudgetHandler{repos: repos}
}

func (h *BudgetHandler) List(c *gin.Context) {
	// TODO: List all budgets for authenticated user
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *BudgetHandler) Get(c *gin.Context) {
	// TODO: Get budget by month
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *BudgetHandler) Create(c *gin.Context) {
	// TODO: Create new monthly budget
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *BudgetHandler) Update(c *gin.Context) {
	// TODO: Update budget name
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	// TODO: Delete budget (cascade will handle related data)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
