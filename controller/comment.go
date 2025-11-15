package controller

import (
	"api-undangan/database"
	"api-undangan/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
	comment := models.Comment{
		Name: req.Name,
		Message: req.Message,
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