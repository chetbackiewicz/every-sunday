package handlers

import (
	"net/http"
	"strconv"

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
	month := c.Param("month")
	// Assume userID from middleware or mock
	userIDVal, _ := c.Get("userID")
	if userIDVal == nil { userIDVal = 1 }
	userID := userIDVal.(int)

	budget, err := h.repos.Budget.GetByMonth(c, userID, month)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	tree, err := h.repos.Category.GetTree(c, budget.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category tree"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": tree})
}

type CreateCategoryRequest struct {
	Name      string   `json:"name" binding:"required"`
	Projected float64  `json:"projected" binding:"gte=0"`
	ParentID  *int     `json:"parentId"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	month := c.Param("month")
	// Assume userID from middleware or mock
	userIDVal, _ := c.Get("userID")
	if userIDVal == nil { userIDVal = 1 }
	userID := userIDVal.(int)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget, err := h.repos.Budget.GetByMonth(c, userID, month)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	// Check depth if parent provided
	if req.ParentID != nil {
		depth, err := h.repos.Category.GetDepth(c, *req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent category"})
			return
		}
		if depth >= 2 { // Max depth 2 (0 -> 1 -> 2) means max parent depth is 1
			c.JSON(http.StatusBadRequest, gin.H{"error": "Max category depth exceeded"})
			return
		}
	}

	cat, err := h.repos.Category.Create(c, budget.ID, req.Name, req.Projected, req.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

type UpdateCategoryRequest struct {
	Name         *string  `json:"name"`
	Projected    *float64 `json:"projected"`
	ManualActual *float64 `json:"manualActual"`
}

func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repos.Category.Update(c, id, req.Name, req.Projected, req.ManualActual)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repos.Category.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
