package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type CellReferenceHandler struct {
	repos *repositories.Repositories
}

func NewCellReferenceHandler(repos *repositories.Repositories) *CellReferenceHandler {
	return &CellReferenceHandler{repos: repos}
}

func (h *CellReferenceHandler) Create(c *gin.Context) {
	// TODO: Create cell reference linking category to transaction
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *CellReferenceHandler) Delete(c *gin.Context) {
	// TODO: Delete cell reference (unlink transaction from category)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
