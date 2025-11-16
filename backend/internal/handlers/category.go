package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repos *repositories.Repositories
}

func NewCategoryHandler(repos *repositories.Repositories) *CategoryHandler {
	return &CategoryHandler{repos: repos}
}

func (h *CategoryHandler) List(c *gin.Context) {
	// TODO: Get category tree for budget
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	// TODO: Create category with depth check (max 2 levels deep)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	// TODO: Update category name, projected, or manual_actual
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	// TODO: Delete category (cascade handles children and paths)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
