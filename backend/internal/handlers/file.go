package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	repos *repositories.Repositories
}

func NewFileHandler(repos *repositories.Repositories) *FileHandler {
	return &FileHandler{repos: repos}
}

func (h *FileHandler) Upload(c *gin.Context) {
	// TODO: Handle CSV file upload with validation:
	// 1. Check file count (max 10 per budget)
	// 2. Parse CSV with proper column mapping
	// 3. Store file metadata in csv_files table
	// 4. Insert transactions into transactions table
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
