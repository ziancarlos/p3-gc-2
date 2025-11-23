package controllers

import (
	"net/http"
	"p3-graded-challenge-2-ziancarlos/middleware"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type LoginRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login and get JWT token
// @Description Login with user ID to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}
