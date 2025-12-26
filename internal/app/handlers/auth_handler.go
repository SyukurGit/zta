package handlers

import (
	"backend-zerotrust-skripsi/internal/app/models"
	"backend-zerotrust-skripsi/internal/database"
	"backend-zerotrust-skripsi/internal/security/token"
	"backend-zerotrust-skripsi/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler adalah struct yang menampung dependency JWT
type AuthHandler struct {
	JWTService *token.JWTService
}

// NewAuthHandler adalah Constructor untuk membuat AuthHandler
func NewAuthHandler(jwtService *token.JWTService) *AuthHandler {
	return &AuthHandler{JWTService: jwtService}
}

// Login sekarang adalah Method milik AuthHandler (bukan fungsi biasa)
func (h *AuthHandler) Login(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	// 1. Cari User
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, "Username atau Password salah", nil)
		return
	}

	// 2. Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, "Username atau Password salah", nil)
		return
	}

	// 3. Generate Token (Panggil method dari h.JWTService yang di-inject)
	tokenString, err := h.JWTService.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal membuat token", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", gin.H{
		"token": tokenString,
		"user": gin.H{
			"username": user.Username,
			"role":     user.Role,
		},
	})
}