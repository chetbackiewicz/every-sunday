package handlers

import (
	"net/http"

	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repos *repositories.Repositories
}

func NewAuthHandler(repos *repositories.Repositories) *AuthHandler {
	return &AuthHandler{repos: repos}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// TODO: Implement user registration with bcrypt password hashing
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	// TODO: Implement login with JWT token generation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// TODO: Implement refresh token logic
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement logout (token invalidation if using blacklist)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	// TODO: Return current user info from JWT claims
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
