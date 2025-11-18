package controller

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type CommentRequest struct {
	Name		string `json:"name" binding:"required"`
	Message		string `json:"message" binding:"required"`
}

func GetComments(c *gin.Context){
	var comments []models.Comment
	if err := database.DB.Order("created_at desc").Find(&comments).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch comments",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"message": "Get comments successfully!",
	})
}

func PostComment(c *gin.Context){
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name dan message wajib di isi",
		})
		return
	}
	if config.ContainsBadword(req.Message) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Jangan berkata kasar ya kak",
		})
		return
	}
	lastMinute := time.Now().Add(-1 * time.Minute)
	var recent models.Comment
	err := database.DB.Where("ip_address = ? AND created_at >= ?", c.ClientIP(), lastMinute).First(&recent).Error
	if err == nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Tunggu 1 menit sebelum ucapan lagi",
		})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cek spam gagal",
		})
		return
	}
	comment := models.Comment{
		Name: req.Name,
		Message: req.Message,
		IPAddress: c.ClientIP(),
	}
	if err := database.DB.Create(&comment).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to post comment",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": comment,
		"message": "Post comment successfully!",
	})
}