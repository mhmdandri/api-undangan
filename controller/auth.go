package controller

import (
	"api-undangan/auth"
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Login(c *gin.Context){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
    return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
    return
	}
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal membuat token",
		})
		return
	}
	ttl := int(config.Cfg.JWTExpiresIn.Seconds())
	domain := "localhost:3000"
	c.SetCookie("token", token, ttl, "/", domain, true, true)
	
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"name": user.Name,
			"email": user.Email,
			"role": user.Role,
		},
	})
}